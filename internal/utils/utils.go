// internal/utils/utils.go
package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ClearScreen clears the terminal screen (LIMPIA PANTALLA).
func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

// WaitForEnter waits for the user to press Enter.
func WaitForEnter() {
	fmt.Print("Presione Enter para continuar...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// ReadLine prompts the user and reads input from stdin.
func ReadLine(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

// ToInt converts a string to int.
func ToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// ToFloat converts a string to float64.
func ToFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// Normalize trims spaces and standardizes lowercase.
func Normalize(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

// IsEmpty checks if a string is empty after trimming.
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}
