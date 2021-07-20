package types

// PropertyMap is a II.22.35 PropertyMap representation.
type PropertyMap struct {
	Parent       Index `table:"TypeDef"`
	PropertyList List  `table:"Property"`
}
