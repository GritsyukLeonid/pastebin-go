package model

type ShortURL struct {
	Original string `json:"original"`
	hash     string `json:"hash"`
}

func NewShortURL(original string) *ShortURL {
	return &ShortURL{Original: original}
}

func (s *ShortURL) SetHash(h string) {
	s.hash = h
}

func (s *ShortURL) GetHash() string {
	return s.hash
}

func (s *ShortURL) GetTypeName() string {
	return "ShortURL"
}
