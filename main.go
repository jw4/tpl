package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"text/template"
)

func main() {
	var dataFile string
	flag.StringVar(&dataFile, "d", "", "JSON data file (optional, reads from stdin if not provided)")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Error: exactly one template file is required\n")
		flag.Usage()
		os.Exit(1)
	}

	templateFile := flag.Arg(0)

	if err := renderTemplate(dataFile, templateFile, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func renderTemplate(dataFile, templateFile string, output io.Writer) error {
	data, err := loadData(dataFile)
	if err != nil {
		return fmt.Errorf("loading data: %w", err)
	}

	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	if err := tmpl.Execute(output, data); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	return nil
}

func loadData(filename string) (map[string]any, error) {
	if filename == "" {
		// Check if stdin has data available
		stat, err := os.Stdin.Stat()
		if err != nil {
			return make(map[string]any), nil
		}

		// If stdin is a terminal (no piped data), return empty map
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			return make(map[string]any), nil
		}

		// Read from stdin
		var data map[string]any
		if err := json.NewDecoder(os.Stdin).Decode(&data); err != nil {
			return nil, fmt.Errorf("decoding JSON from stdin: %w", err)
		}
		return data, nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	var data map[string]any
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, fmt.Errorf("decoding JSON from file: %w", err)
	}

	return data, nil
}
