package internal

type Filter struct {
	Kind       string
	Value      string
	Start      int
	End        int
	Confidence float64
	EntityID   any
}
