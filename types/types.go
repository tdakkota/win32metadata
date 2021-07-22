package types

type (
	// Index represents simple index field type.
	Index = uint32
	// List represents list index field type.
	List [2]Index
	// Blob represents #Blob heap index type.
	Blob = []byte
	// Signature represents Signature blob.
	Signature Blob
	// GUID represents #GUID heap index type.
	GUID = Index
)
