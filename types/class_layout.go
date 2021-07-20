package types

// ClassLayout is a II.22.8 ClassLayout representation.
type ClassLayout struct {
	PackingSize uint16
	ClassSize   uint32
	Parent      Index `table:"TypeDef"`
}
