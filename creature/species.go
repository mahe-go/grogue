package creature

type Species struct {
	Movement int
	Rune     rune
}

func NewSpecies(movement int, r rune) *Species {
	return &Species{movement, r}
}
