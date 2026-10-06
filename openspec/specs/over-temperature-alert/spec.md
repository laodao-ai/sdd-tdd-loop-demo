# over-temperature-alert Specification

## Purpose

Flags temperature readings above a fixed safety limit so an over-temperature condition is not lost in the stream of normal readings, without disturbing the normal output.

## Requirements

### Requirement: Alert on readings above the limit
The system SHALL write an alert line to stderr for every valid reading strictly greater than 80.0°C. The alert line SHALL state the reading and the limit in the form `ALERT: <reading>°C > 80.0°C`, with the reading shown to one decimal place.

#### Scenario: Reading above the limit
- **WHEN** the input line is `81.2`
- **THEN** stderr receives `ALERT: 81.2°C > 80.0°C`

#### Scenario: Reading just above the limit
- **WHEN** the input line is `80.1`
- **THEN** stderr receives `ALERT: 80.1°C > 80.0°C`

#### Scenario: Reading exactly at the limit
- **WHEN** the input line is `80.0`
- **THEN** no alert is written

#### Scenario: Reading below the limit
- **WHEN** the input line is `23.5` or `-4.0`
- **THEN** no alert is written

### Requirement: Readings stay on stdout unchanged
The system SHALL continue to print every valid reading to stdout as `<reading>°C` with one decimal place, whether or not it triggers an alert. Alerts SHALL NOT be written to stdout.

#### Scenario: Alerting reading is still printed
- **WHEN** the input line is `81.2`
- **THEN** stdout receives `81.2°C`
- **AND** stdout does not contain `ALERT`

### Requirement: Invalid lines never alert
The system SHALL NOT write an alert for a line that is not a valid temperature, and SHALL continue processing subsequent lines.

#### Scenario: Invalid line between readings
- **WHEN** the input lines are `81.2`, `hot`, `82.0`
- **THEN** stderr receives an alert for `81.2`, a parse error for `hot`, and an alert for `82.0`
- **AND** the parse error is not reported as an alert
