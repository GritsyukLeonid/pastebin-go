package model

type ShortURL struct {
	Original string `json:"original"`
	ID       string `json:"id"`
}

func NewShortURL(original string, hash string) *ShortURL {
	return &ShortURL{Original: original, ID: hash}
}

func (s *ShortURL) GetHash() string {
	return s.ID
}

func (s *ShortURL) GetTypeName() string {
	return "ShortURL"
}
