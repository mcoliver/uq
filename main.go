package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/fatih/color"
)

const version = "1.0.0"

func main() {
	helpFlag := flag.Bool("help", false, "Show help message")
	versionFlag := flag.Bool("version", false, "Show version information")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: uq [options] <url>\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if *versionFlag {
		fmt.Printf("uq version %s\n", version)
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Error: A single URL must be provided.")
		flag.Usage()
		os.Exit(1)
	}

	inputURL := args[0]
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Invalid URL: %s\n", err)
		os.Exit(1)
	}

	output := parseURL(parsedURL)

	encodedOutput, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to encode output: %s\n", err)
		os.Exit(1)
	}

	colorizeOutput(string(encodedOutput))
}

func parseURL(parsedURL *url.URL) map[string]interface{} {
	baseURL := parsedURL.Scheme + "://" + parsedURL.Host + parsedURL.Path
	queryParams := parsedURL.Query()

	output := map[string]interface{}{
		"base_url": baseURL,
		"query":    make(map[string]interface{}),
	}

	for key, values := range queryParams {
		if len(values) > 0 {
			value := values[0]
			if subURL, err := url.Parse(value); err == nil && subURL.Scheme != "" {
				output["query"].(map[string]interface{})[key] = parseURL(subURL)
			} else {
				output["query"].(map[string]interface{})[key] = value
			}
		}
	}

	return output
}

func colorizeOutput(jsonString string) {
	keyColor := color.New(color.FgCyan).SprintFunc()
	stringColor := color.New(color.FgGreen).SprintFunc()

	lines := strings.Split(jsonString, "\n")
	for _, line := range lines {
		line = strings.ReplaceAll(line, "\"", "") // Strip extra quotes for formatting.
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			fmt.Printf("  %s: %s\n", keyColor(strings.TrimSpace(parts[0])), stringColor(strings.TrimSpace(parts[1])))
		} else {
			fmt.Println(line)
		}
	}
}
