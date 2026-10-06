package main

import (
	"strings"
	"testing"
)

func TestRunPrintsReadings(t *testing.T) {
	var out, errOut strings.Builder
	run(strings.NewReader("23.5\n80.0\n"), &out, &errOut)
	if got, want := out.String(), "23.5°C\n80.0°C\n"; got != want {
		t.Errorf("stdout = %q, want %q", got, want)
	}
	if errOut.Len() != 0 {
		t.Errorf("stderr = %q, want empty", errOut.String())
	}
}

func runString(t *testing.T, input string) (stdout, stderr string) {
	t.Helper()
	var out, errOut strings.Builder
	run(strings.NewReader(input), &out, &errOut)
	return out.String(), errOut.String()
}

func TestRunAlerts(t *testing.T) {
	tests := []struct{ name, input, want string }{
		{"over", "81.2\n", "ALERT: 81.2°C > 80.0°C\n"},
		{"just over", "80.1\n", "ALERT: 80.1°C > 80.0°C\n"},
		{"no trailing newline", "81.2", "ALERT: 81.2°C > 80.0°C\n"},
		{"rounds to threshold", "80.0000001\n", "ALERT: 80.0°C > 80.0°C\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, stderr := runString(t, tt.input)
			if stderr != tt.want {
				t.Errorf("stderr = %q, want %q", stderr, tt.want)
			}
		})
	}
}

func TestRunNoAlert(t *testing.T) {
	for _, in := range []string{"80.0", "23.5", "-4.0"} {
		t.Run(in, func(t *testing.T) {
			_, stderr := runString(t, in+"\n")
			if stderr != "" {
				t.Errorf("stderr = %q, want empty", stderr)
			}
		})
	}
}

func TestRunAlertNeverOnStdout(t *testing.T) {
	stdout, _ := runString(t, "81.2\n")
	if stdout != "81.2°C\n" {
		t.Errorf("stdout = %q, want %q", stdout, "81.2°C\n")
	}
	if strings.Contains(stdout, "ALERT") {
		t.Errorf("stdout contains ALERT: %q", stdout)
	}
}

func TestRunInvalidLineBetweenReadings(t *testing.T) {
	stdout, stderr := runString(t, "81.2\nhot\n82.0\n")
	if want := "81.2°C\n82.0°C\n"; stdout != want {
		t.Errorf("stdout = %q, want %q", stdout, want)
	}
	lines := strings.Split(strings.TrimSuffix(stderr, "\n"), "\n")
	if len(lines) != 3 {
		t.Fatalf("stderr has %d lines, want 3: %q", len(lines), stderr)
	}
	if lines[0] != "ALERT: 81.2°C > 80.0°C" {
		t.Errorf("line 0 = %q", lines[0])
	}
	if !strings.Contains(lines[1], `parse temperature "hot"`) || strings.HasPrefix(lines[1], "ALERT") {
		t.Errorf("line 1 = %q, want parse error", lines[1])
	}
	if lines[2] != "ALERT: 82.0°C > 80.0°C" {
		t.Errorf("line 2 = %q", lines[2])
	}
}

func TestRunBlankLineDoesNotAlert(t *testing.T) {
	_, stderr := runString(t, "81.2\n\n82.0\n")
	lines := strings.Split(strings.TrimSuffix(stderr, "\n"), "\n")
	alerts := 0
	for _, l := range lines {
		if strings.HasPrefix(l, "ALERT:") {
			alerts++
		}
	}
	if len(lines) != 3 || alerts != 2 {
		t.Errorf("want 3 lines with 2 alerts, got %d lines, %d alerts: %q", len(lines), alerts, stderr)
	}
}
