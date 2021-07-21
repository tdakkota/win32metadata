package types

// GenericParamConstraint is a II.22.21 GenericParamConstraint representation.
type GenericParamConstraint struct {
	Owner      Index `table:"GenericParam"`
	Constraint TypeDefOrRef
}
