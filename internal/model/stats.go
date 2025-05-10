package model

type Stats struct {
	PasteHash string `json:"paste_hash"`
	Views     int    `json:"views"`
}

func NewStats(hash string) *Stats {
	return &Stats{
		PasteHash: hash,
		Views:     0,
	}
}

func (s *Stats) IncrementViews() {
	s.Views++
}

func (s *Stats) GetTypeName() string {
	return "Stats"
}
