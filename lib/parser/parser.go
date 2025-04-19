package parser

import (
	"os"

	"gopkg.in/yaml.v3"
)

// ParseYamlFile function:
// out argument must be passed as reference
func ParseYamlFile(path string, out any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, out)
	if err != nil {
		return err
	}

	return nil
}

