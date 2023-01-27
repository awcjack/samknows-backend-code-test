package writer

import (
	"os"

	"github.com/awcjack/samknows-backend-code-test/types"
)

var (
	basePath = "output"
)

type ioWriter struct{}

func NewIOWriter() ioWriter {
	return ioWriter{}
}

// write multiple file to filesystem
func (w ioWriter) WriteMultipleOutput(outputs []types.OutputFormat) error {
	for _, output := range outputs {
		err := w.WriteOutput(output.Name, output.Content)
		if err != nil {
			return err
		}
	}

	return nil
}

// write one file to filesystem
func (w ioWriter) WriteOutput(name string, content []byte) error {
	err := os.WriteFile(basePath+"/"+name, content, 0644)
	if err != nil {
		return err
	}

	return nil
}
