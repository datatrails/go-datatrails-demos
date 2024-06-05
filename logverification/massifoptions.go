package logverification

type MassifOptions struct {

	// nonLeafNode is an optional suppression
	//
	//	of errors that occur due to attempting to get
	//  a massif based on a non leaf node mmrIndex.
	nonLeafNode bool
}

type MassifOption func(*MassifOptions)

// WithNonLeafNode is an optional suppression
//
//	of errors that occur due to attempting to get
//	a massif based on a non leaf node mmrIndex.
func WithNonLeafNode(nonLeafNode bool) MassifOption {
	return func(mo *MassifOptions) { mo.nonLeafNode = nonLeafNode }
}

// ParseMassifOptions parses the given options into a MassifOptions struct
func ParseMassifOptions(options ...MassifOption) MassifOptions {
	massifOptions := MassifOptions{
		nonLeafNode: false, // default to erroring on non leaf nodes
	}

	for _, option := range options {
		option(&massifOptions)
	}

	return massifOptions
}
