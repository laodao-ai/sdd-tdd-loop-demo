# Tasks

## 1. Threshold check in `sensor`

- [ ] 1.1 Add a failing table test `TestIsOverTemp` in `sensor/sensor_test.go` covering 81.2, 80.1, 80.0 (exactly at the limit), 23.5 and -4.0; verify it fails with `go test ./sensor` because `IsOverTemp` does not exist yet
- [ ] 1.2 Add exported `MaxCelsius = 80.0` and `IsOverTemp(celsius float64) bool` (strictly greater than) in `sensor/sensor.go`; verify `go test ./sensor` passes

## 2. Testable loop and alert output in `main`

- [ ] 2.1 Extract the read loop into `run(in io.Reader, out, errOut io.Writer)` with `main()` calling `run(os.Stdin, os.Stdout, os.Stderr)`, adding a `main_test.go` test first that feeds `23.5` and `81.2` and asserts stdout is `23.5°C\n81.2°C\n` with nothing on stderr (the alert is not added yet); verify `go test ./...` passes before and after the extraction
- [ ] 2.2 Add failing `run` tests for the spec scenarios: `81.2` and `80.1` alert on stderr as `ALERT: <reading>°C > 80.0°C`; `80.0`, `23.5` and `-4.0` do not alert; stdout never contains `ALERT`; input `81.2`, `hot`, `82.0` gives two alerts plus one parse error on stderr; verify the new tests fail for the right reason
- [ ] 2.3 Write the alert to `errOut` after printing the reading, using `sensor.IsOverTemp` and `sensor.MaxCelsius`; verify `go test ./...` passes

## 3. Docs and end-to-end check

- [ ] 3.1 Update the README Run section to describe the alert on stderr and the 80.0°C limit; verify the documented `printf '23.5\n81.2\n' | go run .` command shows `81.2°C` on stdout and the `ALERT` line on stderr (for example with `2>&1 >/dev/null` to see only stderr)
