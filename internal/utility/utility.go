package utility

import "fmt"

func NewError(message string) error {
	return fmt.Errorf("processinfo: %s", message)
}

func FormatError(message string, err error) error {
	return fmt.Errorf("processinfo: %s: %w", message, err)
}
