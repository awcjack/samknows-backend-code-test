package writer

import (
	"github.com/awcjack/samknows-backend-code-test/types"
)

// interface that expect to be provided in writer implementation
type Writer interface {
	WriteMultipleOutput([]types.OutputFormat) error
	WriteOutput(name string, content []byte) error
}
