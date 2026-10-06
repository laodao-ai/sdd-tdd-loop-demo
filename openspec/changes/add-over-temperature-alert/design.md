# Design

## Context

`main.go` reads stdin in `main()`, calls `sensor.ParseCelsius`, and writes straight to `os.Stdout` / `os.Stderr`. The stdout/stderr split that the specs depend on has no test seam today; only `ParseCelsius` is tested. See proposal.md for motivation and the spec for required behavior.

## Goals / Non-Goals

**Goals:**
- Put the threshold decision in a pure, table-testable function.
- Make the stdout/stderr split testable without running a subprocess.

**Non-Goals:**
- A configurable threshold, alert levels, or any structured alert type.

## Decisions

**1. Pure check in `sensor`: `IsOverTemp(celsius float64) bool`, backed by an exported `MaxCelsius = 80.0` constant.**
The comparison is strictly `celsius > maxCelsius`, which gives the exact-80.0 boundary from the spec. It sits beside `ParseCelsius` and is tested in the same table style. Alternative: take the limit as a parameter. Rejected for now, since the limit is hardcoded and a parameter adds an unused knob. It is an easy change if a flag is added later.

**2. Move the loop into `run(in io.Reader, out, errOut io.Writer)` in `main.go`; `main()` just calls `run(os.Stdin, os.Stdout, os.Stderr)`.**
Tests can then feed lines and assert on both buffers, covering the stdout-unchanged and invalid-line scenarios. Alternative: test through `go run` / `exec.Command`. Rejected as slower and noisier. Alternative: leave `main.go` untested and trust the pure function. Rejected, because the stderr-vs-stdout requirement would go unchecked.

**3. The alert text is formatted in `run`, not in `sensor`.**
Presentation stays with the output code, and `sensor` stays a parsing/threshold package. The alert prints the limit from `sensor.MaxCelsius`, so the message and the check can't drift apart.

**4. Order per valid line: print the reading to stdout first, then the alert to stderr.**
Matches the spec examples and keeps the existing line first when both streams share a terminal.

## Risks / Trade-offs

- [`strconv.ParseFloat` accepts `Inf` and `NaN`] → `+Inf` would alert and `NaN` would not. This is existing parse behavior and outside this change's scope, so it is noted and left unspecified. Rejecting non-finite values would be a separate change to the parse requirement.
- [Refactoring `main` into `run` touches existing code] → The existing stdout behavior is covered by a test written first, before the alert logic is added.
- [stdout and stderr interleave nondeterministically when both go to a terminal] → Accepted; this is inherent to the chosen separate-stream design.
