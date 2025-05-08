package model

import (
	"time"
)

type Paste struct {
	id        int64
	hash      string
	content   string
	createdAt time.Time
	expiresAt time.Time
	views     int
	metrics   *Stats
}

func NewPaste(content string, ttl time.Duration) *Paste {
	now := time.Now()
	return &Paste{
		content:   content,
		createdAt: now,
		expiresAt: now.Add(ttl),
		views:     0,
	}
}

func (p *Paste) ID() int64 {
	return p.id
}

func (p *Paste) SetID(id int64) {
	p.id = id
}

func (p *Paste) Hash() string {
	return p.hash
}

func (p *Paste) SetHash(h string) {
	p.hash = h
}

func (p *Paste) Content() string {
	return p.content
}

func (p *Paste) SetContent(content string) {
	p.content = content
}

func (p *Paste) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Paste) ExpiresAt() time.Time {
	return p.expiresAt
}

func (p *Paste) SetExpiresAt(t time.Time) {
	p.expiresAt = t
}

func (p *Paste) IncrementViews() {
	p.views++
}

func (p *Paste) Views() int {
	return p.views
}

func (p *Paste) AttachStats(s *Stats) {
	p.metrics = s
}

func (p *Paste) Stats() *Stats {
	return p.metrics
}

func (p *Paste) GetTypeName() string {
	return "Paste"
}
