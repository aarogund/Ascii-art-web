package main

import (
	"fmt"
	"os"
	"strings"
)

// func main() {
// 	if len(os.Args) < 2 || len(os.Args) > 3 {
// 		fmt.Println("Usage: go run . [STRING] [BANNER]")
// 		fmt.Println("EX: go run . something standard")
// 		return
// 	}

// 	filename := "standard.txt"

// 	if len(os.Args) == 3 {
// 		filename = os.Args[2] + ".txt"
// 	}

// 	bannerMap := loadBanner(filename)
// 	if bannerMap == nil {
// 		fmt.Println("Usage: go run . [STRING] [BANNER]")
// 		fmt.Println("EX: go run . something standard")
// 		return
// 	}
// 	result, err := printArt(os.Args[1], bannerMap)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Print(result)
// }

func loadBanner(bannerName string) map[rune][]string {

	data, err := os.ReadFile(bannerName + ".txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	dataString := string(data)
	// Normalize line endings
	dataString = strings.ReplaceAll(dataString, "\r\n", "\n")

	dataLines := strings.Split(dataString, "\n")

	result := make(map[rune][]string)

	for char := rune(32); char <= 126; char++ {
		// i.e ( ' ' , '!' , '"' , '#' ... 'A' , 'B' ... 'z' , '~' )

		index := (char - 32) * 9
		result[char] = dataLines[index+1 : index+9]
	}

	return result
}

func printArt(input string, bannerMap map[rune][]string) (string, error) {
	result := ""
	lines := strings.Split(input, "\\n")
	for _, line := range lines {
		if line == "" {
			result += "\n"
			continue
		}
		for _, char := range line {
			if _, ok := bannerMap[char]; !ok {
				return "", fmt.Errorf("character %q is not supported", char)
			}
		}
		for i := 0; i < 8; i++ {
			for _, char := range line {
				result += (bannerMap[char][i])
			}
			result += "\n"
		}
	}
	return result, nil
}
func printPlain(wordArts []string) string {
	result := ""
	for i := 0; i < 8; i++ {
		for j, w := range wordArts {
			lines := strings.Split(w, "\n")
			result += (lines[i])
			if j < len(wordArts)-1 {
				result += (strings.Repeat(" ", 6))
			}
		}
		result += ("\n")
		continue
	}
	
	
	return result
}
func findPosition(s string, w string) []int {
	positions := []int{}
	start := 0
	for {
		pos := strings.Index(s[start:], w)
		if pos == -1 {
			break
		}
		// offset pos by start to get index in original string
		positions = append(positions, start+pos)
		start = start + pos + len(w)
	}
	return positions
}

// isColored reports whether position i falls within
// any colored range defined by positions and subLen.
func isColored(i int, positions []int, subLen int) bool {

	for j := 0; j < len(positions); j++ {
		if i >= positions[j] && i < positions[j]+subLen {
			return true
		}
	}
	return false
}

// getColor returns the ANSI escape code for the given color.
// Supports named colors, rgb() format, and raw ANSI codes.
func getColor(color string) string {
	color = strings.ToLower(color)
	colorsMap := map[string]string{
		"red":    "\033[0;31m",
		"green":  "\033[0;32m",
		"yellow": "\033[0;33m",
		"blue":   "\033[0;34m",
		"purple": "\033[0;35m",
		"white":  "\033[0;37m",
		"cyan":   "\033[0;36m",
	}
	if strings.HasPrefix(color, "rgb(") {

		inner := strings.TrimPrefix(color, "rgb(")

		inner = strings.TrimSuffix(inner, ")")

		parts := strings.Split(inner, ",")

		return fmt.Sprintf("\033[38;2;%s;%s;%sm", parts[0], parts[1], parts[2])

	}
	if val, ok := colorsMap[color]; ok {
		return val
	}
	return color // raw ANSI fallback
}

// printArt prints the ASCII art for input, coloring all
// occurrences of substring using colorCode.
func printArtColor(colorCode string, input string, substring string, bannerMap map[rune][]string) string {
	result := ""
	positions := findPosition(input, substring)
	lines := strings.Split(input, "\\n")
	color := getColor(colorCode)
	for _, line := range lines {
		if line == "" {
			result += "\n"
			continue
		}
		for row := 0; row < 8; row++ {
			if substring == "" {
				substring = input
			}
			for i, char := range line {
				colored := isColored(i, positions, len(substring))
				if colored {
					result += (color)
				}
				result += (bannerMap[char][row])
				
				if colored {
					result += ("\033[0m")
				}
			}
			result += "\n"
		}
	}
	return result
}

func alignment(groupArts []string, plainGroup []string, gaps, termWidth int, flag string) string {
	result := ""
	wordGap := 6
	extraSpaces := 0
	if flag == "justify" && gaps > 0 {
		lineLen := 0
		for _, w := range plainGroup {
			lineLen += getArtWidth(w)
		}
		spacing := termWidth - lineLen - 2
		if spacing < 0 {
			spacing = 0
			wordGap = 6
		} else {
			wordGap = spacing / gaps
			extraSpaces = spacing % gaps
		}
	}
	originalExtra := extraSpaces

	for i := 0; i < 8; i++ {
		lineLen := 0
		for _, w := range plainGroup {
			lines := strings.Split(w, "\n")
			lineLen += len(lines[i])
		}
		lineLen += 6 * gaps
		spacing := termWidth - lineLen - 2
		if spacing < 0 {
			spacing = 0
		}
		leftPad := 0
		rightPad := 0

		switch flag {
		case "left":
			leftPad = 0
			rightPad = spacing
		case "right":
			leftPad = spacing
			rightPad = 0
		case "center":
			leftPad = spacing / 2
			rightPad = spacing - (spacing / 2)
		case "justify":
			leftPad = 0
			rightPad = 0
			if gaps == 0 {
				flag = "left"
			}
		default:
			return "Error: wrong/Incorrect flag"
		}

		result += strings.Repeat(" ", leftPad)
		extraSpaces = originalExtra
		for j, w := range groupArts {
			lines := strings.Split(w, "\n")
			result += lines[i]
			if j < len(groupArts)-1 {
				result += strings.Repeat(" ", wordGap)
				if extraSpaces > 0 {
					result += " "
					extraSpaces--
				}
			}
		}

		result += strings.Repeat(" ", rightPad)
		result += "\n"
	}
	return result
}

func getArtWidth(art string) int {
	lines := strings.Split(art, "\n")
	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}
	return maxWidth
	// first 50 chars of art
}
func printArtColorWeb(colorCode string, input string, substring string, bannerMap map[rune][]string) string {
	if substring == "" {
		substring = input
	}
	result := ""
	positions := findPosition(input, substring)
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		if line == "" {
			result += "\n"
			continue
		}
		for row := 0; row < 8; row++ {
			for i, char := range line {
				colored := isColored(i, positions, len(substring))
				if colored {
					result += fmt.Sprintf("<span style='color:%s'>", colorCode)
				}
				result += (bannerMap[char][row])
				if colored {
					result += "</span>"
				}
			}
			result += "\n"
		}
	}
	return result
}
