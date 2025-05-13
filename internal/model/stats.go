package model

type Stats struct {
	ID    string `json:"id"`
	Views int    `json:"views"`
}

func NewStats(hash string) *Stats {
	return &Stats{
		ID:    hash,
		Views: 0,
	}
}

func (s *Stats) IncrementViews() {
	s.Views++
}

func (s *Stats) GetTypeName() string {
	return "Stats"
}
