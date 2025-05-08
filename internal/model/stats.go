package model

type Stats struct {
	pasteHash string
	views     int
}

func NewStats(hash string) *Stats {
	return &Stats{
		pasteHash: hash,
		views:     0,
	}
}

func (s *Stats) IncrementViews() {
	s.views++
}
func (s *Stats) Views() int {
	return s.views
}
func (s *Stats) Hash() string {
	return s.pasteHash
}
