// Command sdd-tdd-loop-demo reads temperature lines from stdin and prints them.
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/laodao-ai/sdd-tdd-loop-demo/sensor"
)

func main() {
	run(os.Stdin, os.Stdout, os.Stderr)
}

func run(in io.Reader, out, errOut io.Writer) {
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		celsius, err := sensor.ParseCelsius(scanner.Text())
		if err != nil {
			fmt.Fprintln(errOut, err)
			continue
		}
		fmt.Fprintf(out, "%.1f°C\n", celsius)
		if sensor.IsOverTemp(celsius) {
			fmt.Fprintf(errOut, "ALERT: %.1f°C > %.1f°C\n", celsius, sensor.MaxCelsius)
		}
	}
}
