package ascii

import (
	"fmt"
	"io/fs"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// initialize flags prefix
const (
	// OUTPUT_FLAG = "--output="
	// COLOR_FLAG = "--color="
	// ALIGN_FLAG = "--align="
	OUTPUT_DIR = "./outputs"
)

var banners = []string{"apple", "shadow", "standard", "thinkertoy"}

func HandleAsciiArt(str string, subStr string, banner string, flags map[string]string) string {

	fmt.Println(strings.ReplaceAll(str, "\n", "\\n"))
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
	if flags["output"] != "" {
		emptyOutputFile(flags["output"])
	}

	subIndexes := findSubStr(inputString, subStr)
	const ASCII_HEIGHT = 8
	const ASCII_OFFSET = 32

	asciiLines := strings.Split(baseFormat, "\n")
	inputLength := len(inputString)
	var outputData strings.Builder
	var o string
	// Process ASCII art for each row
	for row := 1; row <= ASCII_HEIGHT; row++ {
		var lineData strings.Builder

		for col := 0; col < inputLength; col++ {
			char := inputString[col]
			asciiIndex := int(char) - ASCII_OFFSET
			lineNumber := (asciiIndex * (ASCII_HEIGHT + 1)) + row

			if lineNumber < len(asciiLines) {
				segment := asciiLines[lineNumber]

				// Check if the character is part of the substring to be colored
				isColored := false
				for _, startIdx := range subIndexes {
					if col >= startIdx && col < startIdx+len(subStr) {
						isColored = true
						break
					}
				}

				if isColored && flags["color"] != "" {
					coloredSegment := colorizeText(segment, []int{}, subStr, flags["color"])
					lineData.WriteString(coloredSegment)
				} else {
					lineData.WriteString(segment)
				}
			}
		}

		outputText := lineData.String()
		if flags["align"] != "" {
			outputText = applyAlign(outputText, flags["align"], getTerminalWidth())
		}
		if lineData.Len() > 0 {
			if flags["output"] == "" {
				outputData.WriteString(outputText + "\n")
				o += outputText + "\n"
				fmt.Println(outputText)
			} else {
				outputToFile(flags["output"], outputText)
			}
		}
	}

	return o
}

func getTerminalWidth() int {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin // Use os.Stdin instead of exec.Stdin
	output, err := cmd.Output()
	if err != nil {
		// Fallback to a default value if command fails
		return 80
	}

	parts := strings.Fields(string(output))
	if len(parts) < 2 {
		return 80 // Default width for compatibility
	}

	width, err := strconv.Atoi(parts[1])
	if err != nil {
		return 80
	}

	return width
}

func applyAlign(text string, align string, termWidth int) string {
	textLen := len(text)
	if align == "right" && textLen < termWidth {
		padding := termWidth - textLen
		return strings.Repeat(" ", padding) + text
	} else if align == "center" && textLen < termWidth {
		padding := (termWidth - textLen) / 2
		return strings.Repeat(" ", padding) + text
	}
	// Default to left-align if align type is unknown or fits in width
	return text
}

func findSubStr(str, substr string) []int {
	var indices []int
	offset := 0
	for {
		index := strings.Index(str[offset:], substr)
		if index == -1 {
			break
		}
		indices = append(indices, offset+index)
		offset += index + len(substr)
	}
	return indices
}

func colorizeText(inputText string, subIndexes []int, subStr string, color string) string {
	colorMap := map[string]string{
		"reset":   "\033[0m",
		"red":     "\033[31m",
		"green":   "\033[32m",
		"yellow":  "\033[33m",
		"blue":    "\033[34m",
		"magenta": "\033[35m",
		"cyan":    "\033[36m",
		"gray":    "\033[37m",
		"white":   "\033[97m",
	}
	colorCode := "\033[97m"

	color = strings.ToLower(color)

	if strings.HasPrefix(color, "rgb(") {
		colorCode = convertRgbColorToANSII(color)
	} else if strings.HasPrefix(color, "hsl(") {
		rgbColor := convertHslColorToRgb(color)
		colorCode = convertRgbColorToANSII(rgbColor)
	} else {
		colorCode = colorMap[color]
	}

	// Wrap the whole text if no specific indexes
	if len(subIndexes) == 0 {
		return colorCode + inputText + colorMap["reset"]
	}

	// Return the original text if color is not found
	return inputText
}

func convertRgbColorToANSII(color string) string {
	rgbs := strings.Split(strings.TrimSpace(color[4:len(color)-1]), ",")

	red, err := strconv.Atoi(rgbs[0])
	if checkError(err) {
		os.Exit(1)
	}
	green, err := strconv.Atoi(rgbs[1])
	if checkError(err) {
		os.Exit(1)
	}
	blue, err := strconv.Atoi(rgbs[2])
	if checkError(err) {
		os.Exit(1)
	}
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", red, green, blue)
}

// this part has error
func convertHslColorToRgb(color string) string {
	// Remove 'hsl(' and ')' and split by ','
	color = strings.TrimPrefix(color, "hsl(")
	color = strings.TrimSuffix(color, ")")
	hslParts := strings.Split(color, ",")

	// Parse H, S, L values from the string
	h, _ := strconv.ParseFloat(strings.TrimSpace(hslParts[0]), 64)
	s, _ := strconv.ParseFloat(strings.TrimSpace(hslParts[1]), 64)
	l, _ := strconv.ParseFloat(strings.TrimSpace(hslParts[2]), 64)

	// Convert percentages to fractions
	s /= 100
	l /= 100

	c := (1 - math.Abs(2*l-1)) * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := l - c/2

	var r, g, b float64

	switch {
	case h >= 0 && h < 60:
		r, g, b = c, x, 0
	case h >= 60 && h < 120:
		r, g, b = x, c, 0
	case h >= 120 && h < 180:
		r, g, b = 0, c, x
	case h >= 180 && h < 240:
		r, g, b = 0, x, c
	case h >= 240 && h < 300:
		r, g, b = x, 0, c
	case h >= 300 && h < 360:
		r, g, b = c, 0, x
	}

	// Convert to 0-255 range
	r = (r + m) * 255
	g = (g + m) * 255
	b = (b + m) * 255

	return fmt.Sprintf("rgb(%d, %d, %d)", int(r), int(g), int(b))
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
