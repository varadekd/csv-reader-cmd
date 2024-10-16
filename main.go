package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Record represents a single record in the CSV file.
type Record map[string]string

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <input.csv> [--json] [--print]")
		return
	}

	inputFile := os.Args[1]
	outputFormat := "csv"
	printOutput := false

	// Parse command line arguments
	for _, arg := range os.Args[2:] {
		switch {
		case strings.HasPrefix(arg, "--print"):
			printOutput = true
		case strings.HasPrefix(arg, "--json"):
			outputFormat = "json"
		}
	}

	// Read the CSV file
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read all records from the CSV
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Error reading CSV: %v\n", err)
		return
	}

	if len(records) == 0 {
		fmt.Println("No records found in the CSV file.")
		return
	}

	// Display headers
	headers := records[0]
	fmt.Println("Headers:")
	for i, header := range headers {
		fmt.Printf("%d: %s\n", i, header)
	}

	fmt.Printf("\n\n")
	// Wait for user input to specify column=value or press Enter for first 10 entries
	fmt.Println("Enter column=value to filter or press Enter to see the first 10 entries:")
	var columnFilter string
	fmt.Scanln(&columnFilter) // Wait for user input
	fmt.Printf("\n\n")

	// Prepare data for filtering
	var filteredRecords []Record
	for _, record := range records[1:] {
		rec := make(Record)
		for i, value := range record {
			rec[headers[i]] = value
		}
		filteredRecords = append(filteredRecords, rec)
	}

	// Apply filter if provided
	if columnFilter != "" {
		parts := strings.SplitN(columnFilter, "=", 2)
		if len(parts) != 2 {
			fmt.Println("\nInvalid filter format. Use <column>=<value>.")
			return
		}

		columnName := parts[0]
		filterValue := parts[1]

		// Find the index of the column
		columnIndex := -1
		for i, header := range headers {
			if header == columnName {
				columnIndex = i
				break
			}
		}

		if columnIndex == -1 {
			fmt.Printf("\n Column name '%s' not found.\n", columnName)
			return
		}

		// Filter records
		var filtered []Record
		for _, record := range filteredRecords {
			if record[columnName] == filterValue {
				filtered = append(filtered, record)
			}
		}
		filteredRecords = filtered
		if len(filtered) == 0 {
			fmt.Println("No matching records found.")
			return
		} else {
			filteredRecords = filtered
		}
	} else {
		// Print first 10 records if no filter is applied
		if len(filteredRecords) > 10 {
			filteredRecords = filteredRecords[:10]
		}
	}

	// Create output directory if it doesn't exist
	outputDir := "outputs"
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		fmt.Printf("\n Error creating output directory: %v\n", err)
		return
	}

	// Create output file with timestamp
	timestamp := time.Now().Format("20060102150405")
	outputFile := filepath.Join(outputDir, fmt.Sprintf("output_%s.%s", timestamp, outputFormat))

	var outputData []byte
	var err2 error
	var csvOutput strings.Builder

	if outputFormat == "json" {
		outputData, err2 = json.MarshalIndent(filteredRecords, "", "  ")
	} else {
		// Default to CSV
		writer := csv.NewWriter(&csvOutput)
		writer.Write(headers) // write headers

		for _, record := range filteredRecords {
			recordSlice := make([]string, len(headers))
			for i, header := range headers {
				recordSlice[i] = record[header]
			}
			writer.Write(recordSlice)
		}
		writer.Flush()
		outputData = []byte(csvOutput.String())
	}

	if err2 != nil {
		fmt.Printf("Error creating output data: %v\n", err2)
		return
	}

	// Write output to file
	err = ioutil.WriteFile(outputFile, outputData, 0644)
	if err != nil {
		fmt.Printf("Error writing to output file: %v\n", err)
		return
	}

	if printOutput {
		// Print output to console
		if outputFormat == "json" {
			fmt.Println(string(outputData))
		} else {
			fmt.Println(csvOutput.String())
		}
	} else {
		fmt.Printf("Output written to %s\n", outputFile)
	}
}
