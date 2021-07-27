package types

// Param is a II.22.33 Param representation.
type Param struct {
	Flags    ParamAttributes
	Sequence uint16
	Name     string
}
