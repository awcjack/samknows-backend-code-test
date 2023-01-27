package reader

import (
	"encoding/json"
	"os"

	"github.com/awcjack/samknows-backend-code-test/types"
)

var (
	basePath = "./input"
)

type ioReader struct{}

func NewIOReader() ioReader {
	return ioReader{}
}

// Get all inputs files under directory
func (r ioReader) GetInputs() ([]types.InputFormat, error) {
	entries, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	result := make([]types.InputFormat, 0)

	for _, entry := range entries {
		// ignore directory
		if !entry.IsDir() {
			input, err := r.GetInput(entry.Name())
			if err != nil {
				return nil, err
			}

			result = append(result, input)
		}
	}

	return result, nil
}

// Get input file based on name
func (r ioReader) GetInput(name string) (types.InputFormat, error) {
	content, err := os.ReadFile(basePath + "/" + name)
	if err != nil {
		return types.InputFormat{}, err
	}

	var mesurement []types.Mesurement
	err = json.Unmarshal(content, &mesurement)
	if err != nil {
		return types.InputFormat{}, err
	}

	return types.InputFormat{
		Name:    name,
		Content: mesurement,
	}, nil
}
