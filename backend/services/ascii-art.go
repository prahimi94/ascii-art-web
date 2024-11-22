package ascii

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
)

// initialize flags prefix
const (
	OUTPUT_DIR = "./outputs"
)

func HandleAsciiArt(str string, banner string, flags map[string]string) string {
	// // Read the banner file
	baseFormat, err := readFile(banner + ".txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	// Process and print ASCII art for the input string
	result := printAsciiArt(str, baseFormat, flags)
	return result
}

// readFile function: reads the content of a file and returns it as a string
func readFile(filename string) (string, error) {
	//set directory of banners
	directory := os.DirFS("backend/banners")

	data, err := fs.ReadFile(directory, filename)
	if err != nil {
		return "", err
	}

	cleanedData := strings.ReplaceAll(string(data), "\r", "")
	return string(cleanedData), nil
}

// printAsciiArt function: converts the input string to ASCII art and prints it
func printAsciiArt(inputString string, baseFormat string, flags map[string]string) string {
	const ASCII_HEIGHT = 8
	const ASCII_OFFSET = 32

	inputString = strings.ReplaceAll(inputString, "\r\n", "\\n")
	inputLines := strings.Split(inputString, "\\n")
	asciiLines := strings.Split(baseFormat, "\n")

	var outputData string
	var outputText string
	// Process ASCII art for each row
	for i, inputString := range inputLines {
		inputLength := len(inputString)
		if inputString == "" {
			outputData += "\n"
		}
		for row := 1; row <= ASCII_HEIGHT; row++ {
			var lineData strings.Builder

			for col := 0; col < inputLength; col++ {
				char := inputString[col]
				asciiIndex := int(char) - ASCII_OFFSET
				lineNumber := (asciiIndex * (ASCII_HEIGHT + 1)) + row

				if lineNumber < len(asciiLines) {
					segment := asciiLines[lineNumber]

					lineData.WriteString(segment)

				}
			}

			outputText = lineData.String()

			if lineData.Len() > 0 {
				if flags["output"] == "" {
					outputData += outputText
					if i != len(inputLines)-1 || row != ASCII_HEIGHT {
						outputData += "\n"
					}
				} else {
					outputToFile(flags["output"], outputText)
				}
			}
		}
	}

	return outputData
}

func emptyOutputFile(filePath string) {
	// Check if file exists
	if _, err := os.Stat(filePath); err == nil {
		// File exists, now empty it
		// Open the file with the O_TRUNC flag to truncate it
		file, err := os.OpenFile(OUTPUT_DIR+"/"+filePath, os.O_RDWR|os.O_TRUNC, 0666)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()
		// The file is now empty
		fmt.Println("File exists and has been emptied.")
	}
}

func outputToFile(filePath string, lineData string) error {
	// Ensure the output directory exists
	if err := os.MkdirAll(OUTPUT_DIR, 0755); err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}
	// Open the file in append mode (os.O_APPEND), create it if it doesn't exist (os.O_CREATE)
	outputFile, err := os.OpenFile(OUTPUT_DIR+"/"+filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer outputFile.Close()

	// Write the line data to the file
	_, err = outputFile.WriteString(lineData + "\n")
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}

func checkEmpty(inputString string) {
	if len(inputString) == 0 {
		return
	}
	if inputString == "\\n" {
		fmt.Println()
		return
	}
}

// Check for errors
func checkError(err error) bool {
	if err != nil {
		fmt.Println("Error:", err)
		return true
	}
	return false
}
