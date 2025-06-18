package model

import (
	"time"
)

type Paste struct {
	ID        string    `json:"id"`
	Hash      string    `json:"hash"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
	Views     int       `json:"views"`
}

func NewPaste(content string, ttl time.Duration) *Paste {
	now := time.Now()
	return &Paste{
		Content:   content,
		CreatedAt: now,
		ExpiresAt: now.Add(ttl),
		Views:     0,
	}
}

func (p *Paste) IncrementViews() {
	p.Views++
}

func (p *Paste) GetTypeName() string {
	return "Paste"
}
