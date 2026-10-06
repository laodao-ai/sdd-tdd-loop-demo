# Over-Temperature Alert Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Write an `ALERT: <reading>°C > 80.0°C` line to stderr for every valid reading strictly above 80.0°C, leaving stdout unchanged.

**Architecture:** A pure `IsOverTemp` check plus exported `MaxCelsius` live in `sensor`. The read loop in `main.go` moves into `run(in, out, errOut)` so the stdout/stderr split is testable with buffers; `run` formats and writes the alert.

**Tech Stack:** Go 1.27.1, standard library only.

**Spec:** `openspec/changes/add-over-temperature-alert/` (`proposal.md`, `design.md`, `specs/over-temperature-alert/spec.md`, `tasks.md`)

## Global Constraints

- Limit is hardcoded `80.0`°C; alert only when reading is strictly greater (80.0 does not alert).
- Alert format: `ALERT: <reading>°C > 80.0°C`, reading with one decimal place, written to stderr only.
- Stdout stays `%.1f°C\n` per valid reading, alerting or not; never contains `ALERT`.
- Per valid line: print the reading to stdout first, then the alert to stderr.
- Invalid lines: parse error on stderr (existing behavior), no alert, processing continues.
- Standard library only; no new dependencies; no configurable threshold, exit codes, or other units.
- The alert prints the limit from `sensor.MaxCelsius`, not a second literal.

## Review Focus

- Reading `80.0000001` (prints as `80.0°C`): it is strictly above the limit, so it alerts as `ALERT: 80.0°C > 80.0°C`. Pinned by a `run` test in Task 2 so the odd-looking line is a deliberate, tested outcome.
- Blank line / whitespace-only line between readings: parse error on stderr, no alert, later lines still processed. Pinned in Task 2.
- Input with no trailing newline (`81.2` with no `\n`): still alerts. Pinned in Task 2.
- `+Inf` / `NaN`: parse accepts them (design.md risk, out of scope); behavior is left unspecified and deliberately not tested.

## File Structure

- `sensor/sensor.go`: add `MaxCelsius`, `IsOverTemp`.
- `sensor/sensor_test.go`: add `TestIsOverTemp`.
- `main.go`: extract `run`, add alert.
- `main_test.go` (create): tests for `run`.
- `README.md`: Run section.

---

### Task 1: Threshold check in `sensor`

**Files:**
- Modify: `sensor/sensor.go`
- Test: `sensor/sensor_test.go`

**Interfaces:**
- Consumes: nothing.
- Produces: `const MaxCelsius = 80.0` (untyped float const) and `func IsOverTemp(celsius float64) bool` in package `sensor`.

- [ ] **Step 1: Write the failing test** `TestIsOverTemp` in `sensor/sensor_test.go`, table style like `TestParseCelsius`: `81.2`→true, `80.1`→true, `80.0`→false (exactly at limit), `23.5`→false, `-4.0`→false.

- [ ] **Step 2: Run to verify it fails**

Run: `go test ./sensor`
Expected: FAIL to compile, `undefined: IsOverTemp`.

- [ ] **Step 3: Implement `const MaxCelsius = 80.0` and `func IsOverTemp(celsius float64) bool` in `sensor/sensor.go`**

Return `celsius > MaxCelsius`; add doc comments in the file's existing style.

- [ ] **Step 4: Run to verify it passes**

Run: `go test ./sensor`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add sensor/sensor.go sensor/sensor_test.go
git commit -m "feat(sensor): add IsOverTemp threshold check"
```

---

### Task 2: Testable loop and alert output in `main`

**Files:**
- Modify: `main.go`
- Create: `main_test.go`

**Interfaces:**
- Consumes: `sensor.ParseCelsius(raw string) (float64, error)`, `sensor.IsOverTemp(float64) bool`, `sensor.MaxCelsius` (Task 1).
- Produces: `func run(in io.Reader, out, errOut io.Writer)` in package `main`; `main()` calls `run(os.Stdin, os.Stdout, os.Stderr)`.

- [ ] **Step 1: Write the characterization test** `TestRunPrintsReadings` in `main_test.go` (package `main`): feed `"23.5\n81.2\n"` via `strings.NewReader`, assert `out.String() == "23.5°C\n81.2°C\n"` and `errOut.Len() == 0`. It cannot compile until `run` exists, so the refactor in Step 2 is what makes it pass.

- [ ] **Step 2: Extract `run(in io.Reader, out, errOut io.Writer)` in `main.go`**

Move the existing scanner loop unchanged except writing to `out` / `errOut`; no alert yet.

- [ ] **Step 3: Run to verify it passes**

Run: `go test ./...`
Expected: PASS (existing behavior preserved).

- [ ] **Step 4: Write the failing alert tests in `main_test.go`**, using a helper `runString(t, input string) (stdout, stderr string)`:
  - `TestRunAlerts`, table: `"81.2\n"` → stderr `"ALERT: 81.2°C > 80.0°C\n"`; `"80.1\n"` → `"ALERT: 80.1°C > 80.0°C\n"`; `"81.2"` (no trailing newline) → `"ALERT: 81.2°C > 80.0°C\n"`; `"80.0000001\n"` → `"ALERT: 80.0°C > 80.0°C\n"`.
  - `TestRunNoAlert`, table: `80.0`, `23.5`, `-4.0` → stderr empty.
  - `TestRunAlertNeverOnStdout`: input `"81.2\n"` → stdout `"81.2°C\n"` and does not contain `ALERT`.
  - `TestRunInvalidLineBetweenReadings`: input `"81.2\nhot\n82.0\n"` → stdout `"81.2°C\n82.0°C\n"`; stderr has exactly three lines: alert for 81.2, a line containing `parse temperature "hot"` not starting with `ALERT`, alert for 82.0.
  - `TestRunBlankLineDoesNotAlert`: input `"81.2\n\n82.0\n"` → stderr has two `ALERT:` lines and one non-alert parse-error line.

- [ ] **Step 5: Run to verify the new tests fail for the right reason**

Run: `go test . -run 'TestRunAlerts|TestRunInvalid|TestRunBlank|TestRunAlertNever' -v`
Expected: FAIL on missing alert output (stderr mismatch), not on compile errors; `TestRunNoAlert` passes already.

- [ ] **Step 6: Write the alert in `run`**

After printing the reading to `out`, if `sensor.IsOverTemp(celsius)`, write `ALERT: %.1f°C > %.1f°C\n` with `celsius` and `sensor.MaxCelsius` to `errOut`.

- [ ] **Step 7: Run to verify everything passes**

Run: `go test ./...`
Expected: PASS.

- [ ] **Step 8: Commit**

```bash
git add main.go main_test.go
git commit -m "feat: alert on stderr for readings above 80.0°C"
```

---

### Task 3: Docs and end-to-end check

**Files:**
- Modify: `README.md` (Run section)

**Interfaces:**
- Consumes: the working program from Task 2.
- Produces: nothing.

- [ ] **Step 1: Update the README Run section** to say readings above 80.0°C (strictly) also emit an `ALERT: <reading>°C > 80.0°C` line on stderr, while stdout is unchanged. Keep the existing commands.

- [ ] **Step 2: Verify end to end**

Run: `printf '23.5\n81.2\n' | go run . 2>/dev/null`
Expected: stdout is exactly `23.5°C` and `81.2°C`.

Run: `printf '23.5\n81.2\n' | go run . 2>&1 >/dev/null`
Expected: exactly `ALERT: 81.2°C > 80.0°C`.

- [ ] **Step 3: Commit**

```bash
git add README.md
git commit -m "docs: describe over-temperature alert in README"
```

---

## Self-Review Notes

- Spec coverage: alert requirement and all four scenarios → Tasks 1–2; stdout-unchanged → Task 2 (`TestRunPrintsReadings`, `TestRunAlertNeverOnStdout`); invalid lines → Task 2; docs/e2e → Task 3. Every `tasks.md` item (1.1–3.1) maps to a step above.
- Names are consistent across tasks: `MaxCelsius`, `IsOverTemp`, `run`.
