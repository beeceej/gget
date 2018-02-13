package gget

import "log"

// ToStdOut provides a strategy for printing file contents to std out
type ToStdOut struct {
	// Size specifies the max length that you want to print,
	// If not specified will display all the bytes as a String
	Size int
}

// Handle is a strategy for printing file contents out to std out
func (t *ToStdOut) Handle(r result) error {
	if t.Size != 0 {
		log.Printf(string(r.b)[:t.Size])
	} else {
		log.Printf(string(r.b))
	}

	return nil
}
