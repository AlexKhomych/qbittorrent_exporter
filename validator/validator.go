package validator

import (
	"fmt"
	"net/url"
	"os"
)

func ValidatePath(path string, checkIfDir bool) error {
	if len(path) == 0 {
		return fmt.Errorf("validation failed, empty path value")
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("validation failed, encountered an error: %v", err)
	}
	if checkIfDir && !info.IsDir() {
		return fmt.Errorf("validation failed, not a directory")
	}
	return nil
}

func ValidateURL(input string) error {
	_, err := url.ParseRequestURI(input)
	return err
}
