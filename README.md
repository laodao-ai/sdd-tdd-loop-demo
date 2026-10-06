# sdd-tdd-loop-demo

A tiny temperature reader used to demo a minimal AI coding loop:
OpenSpec for the spec, Superpowers for test-first implementation.

The tag `start` marks the starting point. The demo adds an over-temperature alert from there.

## Run

```bash
go test ./...
printf '23.5\n81.2\n' | go run .
```

The program reads temperatures from stdin (one per line) and prints them to stdout with the unit. Readings strictly above 80.0°C also emit an alert to stderr: `ALERT: <reading>°C > 80.0°C`.

Requires Go 1.27 or later. Standard library only.
