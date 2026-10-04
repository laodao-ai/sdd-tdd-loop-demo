package sensor

import "testing"

func TestParseCelsius(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    float64
		wantErr bool
	}{
		{name: "plain value", raw: "23.5", want: 23.5},
		{name: "surrounding spaces", raw: "  -4.0\n", want: -4.0},
		{name: "not a number", raw: "hot", wantErr: true},
		{name: "empty line", raw: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCelsius(tt.raw)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ParseCelsius(%q) = %v, want an error", tt.raw, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseCelsius(%q) returned error: %v", tt.raw, err)
			}
			if got != tt.want {
				t.Fatalf("ParseCelsius(%q) = %v, want %v", tt.raw, got, tt.want)
			}
		})
	}
}
