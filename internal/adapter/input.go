package adapter

import (
	"fmt"
	"os"
)

func OpenInput(args []string) (*os.File, func(), error) {
	if len(args) >= 1 {
		f, err := os.Open(args[0])
		if err != nil {
			return nil, nil, fmt.Errorf("cannot open events file: %w", err)
		}
		return f, func() { _ = f.Close() }, nil
	}
	return os.Stdin, func() {}, nil
}
