package model

type ShortURL struct {
	original string
	hash     string
}

func NewShortURL(original string) *ShortURL {
	return &ShortURL{original: original}
}

func (s *ShortURL) Original() string {
	return s.original
}

func (s *ShortURL) SetHash(h string) {
	s.hash = h
}

func (s *ShortURL) Hash() string {
	return s.hash
}

func (s *ShortURL) GetTypeName() string {
	return "ShortURL"
}
