package model

import (
	"time"
)

type Paste struct {
	ID        int64     `json:"id"`
	Hash      string    `json:"hash"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
	Views     int       `json:"views"`
	Metrics   *Stats    `json:"metrics"`
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

func (p *Paste) SetID(id int64) {
	p.ID = id
}

func (p *Paste) SetHash(h string) {
	p.Hash = h
}

func (p *Paste) SetContent(content string) {
	p.Content = content
}

func (p *Paste) SetExpiresAt(t time.Time) {
	p.ExpiresAt = t
}

func (p *Paste) IncrementViews() {
	p.Views++
}

func (p *Paste) AttachStats(s *Stats) {
	p.Metrics = s
}

func (p *Paste) Stats() *Stats {
	return p.Metrics
}

func (p *Paste) GetTypeName() string {
	return "Paste"
}
