package paste

import (
	"pastebin-go/internal/model/stats"
	"time"
)

type Paste struct {
	id        int64
	hash      string
	content   string
	createdAt time.Time
	expiresAt time.Time
	views     int
	metrics   *stats.Stats
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

func (p *Paste) GetID() int64 {
	return p.id
}

func (p *Paste) SetID(id int64) {
	p.id = id
}

func (p *Paste) GetHash() string {
	return p.hash
}

func (p *Paste) SetHash(h string) {
	p.hash = h
}

func (p *Paste) GetContent() string {
	return p.content
}

func (p *Paste) SetContent(content string) {
	p.content = content
}

func (p *Paste) GetCreatedAt() time.Time {
	return p.createdAt
}

func (p *Paste) GetExpiresAt() time.Time {
	return p.expiresAt
}

func (p *Paste) SetExpiresAt(t time.Time) {
	p.expiresAt = t
}

func (p *Paste) IncrementViews() {
	p.views++
}

func (p *Paste) GetViews() int {
	return p.views
}

func (p *Paste) AttachStats(s *stats.Stats) {
	p.metrics = s
}

func (p *Paste) GetStats() *stats.Stats {
	return p.metrics
}
