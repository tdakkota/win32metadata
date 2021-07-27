package types

// File is a II.22.19 File representation.
type File struct {
	Flags     FileAttributes
	Name      string
	HashValue Blob
}
