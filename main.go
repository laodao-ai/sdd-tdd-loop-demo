// Command sdd-tdd-loop-demo reads temperature lines from stdin and prints them.
package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/laodao-ai/sdd-tdd-loop-demo/sensor"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		celsius, err := sensor.ParseCelsius(scanner.Text())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		fmt.Printf("%.1f°C\n", celsius)
	}
}
