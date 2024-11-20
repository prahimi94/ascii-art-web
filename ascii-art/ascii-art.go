package ascii

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
)

// initialize flags prefix
const (
	// OUTPUT_FLAG = "--output="
	// COLOR_FLAG = "--color="
	// ALIGN_FLAG = "--align="
	OUTPUT_DIR = "./outputs"
)


func HandleAsciiArt(str string, subStr string, banner string, flags map[string]string) string {

	// // Read the banner file
	baseFormat, err := readFile(banner + ".txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	// Process and print ASCII art for the input string
	oo := printAsciiArt(str, subStr, baseFormat, flags)
	fmt.Println(oo)
	return oo
}

func checkInvalidFlags() {
	// Define allowed flags
	allowedFlags := map[string]bool{
		"color":  true,
		"align":  true,
		"output": true,
	}

	// Iterate over the command-line arguments (excluding the program name)
	for _, arg := range os.Args[1:] {
		// Extract flag name (everything before the '=' sign or flag itself)
		if strings.HasPrefix(arg, "--") {
			// Get the flag name (before '=' if present)
			flagName := strings.TrimPrefix(arg, "--")
			flagName = strings.Split(flagName, "=")[0] // Get the part before '='

			// If the flag name is not allowed, print an error and exit
			if !allowedFlags[flagName] {
				fmt.Printf("Usage: go run . [OPTION] [STRING]\nEX: go run . --color=<color> <substring to be colored> something\n")
				os.Exit(1)
			}
		}
	}
}

// validateInput function: checks if the command-line input is valid
func validateInput(args []string) ([]string, error) {
	//go run . [OPTION] [STRING] [BANNER]
	switch len(args) {
	case 3:
		return []string{args[0], args[1], args[2]}, nil
	case 2:
		return []string{args[0], args[1]}, nil
	case 1:
		return []string{args[0]}, nil
	default:
		return []string{}, fmt.Errorf("\nUsage: go run . [OPTION] [STRING] [BANNER]\n\nEX: go run . flag something standard\n")
	}
}

// readFile function: reads the content of a file and returns it as a string
func readFile(filename string) (string, error) {
	//set directory of banners
	directory := os.DirFS("./banners")
	data, err := fs.ReadFile(directory, filename)
	if err != nil {
		return "", err
	}
	cleanedData := strings.ReplaceAll(string(data), "\r", "")
	return string(cleanedData), nil
}

// printAsciiArt function: converts the input string to ASCII art and prints it
func printAsciiArt(inputString string, subStr string, baseFormat string, flags map[string]string) string {
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
