package types

// Start returns first Index in a List.
func (l List) Start() Index {
	return l[0]
}

// End returns last Index in a List.
func (l List) End() Index {
	return l[1]
}

// Size returns List size.
func (l List) Size() int {
	return int(l.End() - l.Start())
}

// Empty denotes List is empty.
func (l List) Empty() bool {
	return l.Size() < 1
}
