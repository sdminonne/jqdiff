package jqdiff

type DifferenceKind string

const (
	DifferentValue DifferenceKind = "DifferentValue"
	DifferentType  DifferenceKind = "DifferentType"
)

type Diff interface {
}

type typedDiff[R any, A any] struct {
	// selector contains the jq element path
	Selector string
	// Reference is the expected value
	Reference R
	// Actual is the current value
	Actual A
	// Difference is the kind of difference
	Kind DifferenceKind
}
