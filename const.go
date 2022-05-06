package otree

import "errors"

// Error codes
var (
	ErrCannotRemoveRootNode    = errors.New("cannot remove root node")
	ErrDuplicateNodeFound      = errors.New("duplicate node found")
	ErrNodeMustNotHaveSiblings = errors.New("node must not have siblings")
	ErrNodeNotFound            = errors.New("node not found")
	ErrNodesNotInSameTree      = errors.New("nodes not in same tree")
	ErrParentMissing           = errors.New("parent missing")
)
