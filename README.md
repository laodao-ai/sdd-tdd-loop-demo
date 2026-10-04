# sdd-tdd-loop-demo

A tiny temperature reader used to demo a minimal AI coding loop:
OpenSpec for the spec, Superpowers for test-first implementation.

The tag `start` marks the starting point. The demo adds an over-temperature alert from there.

## Run

```bash
go test ./...
printf '23.5\n81.2\n' | go run .
```

Requires Go 1.27 or later. Standard library only.
