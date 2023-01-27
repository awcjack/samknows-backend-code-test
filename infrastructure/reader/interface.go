package reader

import (
	"github.com/awcjack/samknows-backend-code-test/types"
)

// interface that expect to be provided in reader implementation
type Reader interface {
	GetInputs() ([]types.InputFormat, error)
	GetInput(name string) (types.InputFormat, error)
}
