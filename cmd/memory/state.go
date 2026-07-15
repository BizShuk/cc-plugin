package memory

import (
	"github.com/bizshuk/cc-plugin/model"
)

// StateStore wraps model.StateStore for backward compatibility.
type StateStore = model.StateStore

func NewStateStore() (*StateStore, error) {
	return model.NewStateStore()
}

func Fingerprint(text string, entities []string) string {
	return model.Fingerprint(text, entities)
}
