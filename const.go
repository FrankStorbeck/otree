package otree

import "errors"

// a list
const (
	AtStart = -1                 // index number that can be used to indicate the start of a node list
	AtEnd   = int(^uint(0) >> 1) // index number that can be used to indicate the end of a node list

)

// Error codes
var (
	ErrCannotRemoveRootNode    = errors.New("otree: cannot remove root node")
	ErrDuplicateNodeFound      = errors.New("otree: duplicate node found")
	ErrNodeMustNotHaveSiblings = errors.New("otree: node must not have siblings")
	ErrNodeNotFound            = errors.New("otree: node not found")
	ErrNodesNotInSameTree      = errors.New("otree: nodes not in same tree")
	ErrParentMissing           = errors.New("otree: parent missing")
)
