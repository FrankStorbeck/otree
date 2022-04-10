package otree

import "errors"

// Error codes
var (
	ErrDuplicateNodeFound      = errors.New("duplicate node found")
	ErrNodeMustNotHaveSiblings = errors.New("node must not have siblings")
	ErrNodeNotFound            = errors.New("node not found")
	ErrNoNodeFound             = errors.New("no node found ")
	ErrParentMissing           = errors.New("parent missing")
)
