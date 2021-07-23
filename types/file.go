package types

type (
	// FileAttributes describes II.23.1.6 Flags for files [FileAttributes].
	FileAttributes uint32

	// File is a II.22.19 File representation.
	File struct {
		Flags     FileAttributes
		Name      string
		HashValue Blob
	}
)
