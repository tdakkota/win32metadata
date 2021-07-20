package types

// Module is a II.22.30 Module representation.
type Module struct {
	Generation uint16
	Name       string
	Mvid       GUID
	EncId      GUID
	EncBaseId  GUID
}
