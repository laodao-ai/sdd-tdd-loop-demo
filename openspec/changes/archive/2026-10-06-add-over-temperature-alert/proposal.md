# Proposal

## Why

The temperature reader prints every reading but never flags a dangerous one. An over-temperature reading is easy to miss in a stream of normal values, so the reader should call it out.

## What Changes

- Readings strictly above 80.0°C produce an alert line on stderr, in addition to the normal stdout reading line.
- Stdout is unchanged: every valid reading is still printed as `%.1f°C`.
- The threshold is a hardcoded 80.0°C. A reading of exactly 80.0°C does not alert.
- Invalid lines keep today's behavior (error on stderr, no alert, processing continues).

## Capabilities

### New Capabilities
- `over-temperature-alert`: Flags readings above the 80.0°C limit with an alert line on stderr while leaving the stdout reading stream unchanged.

### Modified Capabilities

## Impact

- `sensor` package: a new pure check for whether a reading is over the limit.
- `main.go`: calls the check after parsing and writes the alert to stderr.
- No new dependencies, and no change to existing stdout output.
- Out of scope: a configurable threshold (flag or env var), exit codes, and non-Celsius units.
