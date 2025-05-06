package shorturl

type ShortURL struct {
	original string
	hash     string
}

func New(original string) *ShortURL {
	return &ShortURL{
		original: original,
	}
}

func (s *ShortURL) GetOriginal() string {
	return s.original
}

func (s *ShortURL) SetHash(h string) {
	s.hash = h
}

func (s *ShortURL) GetHash() string {
	return s.hash
}
