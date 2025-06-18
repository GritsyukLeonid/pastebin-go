package handlers

import (
	"github.com/GritsyukLeonid/pastebin-go/internal/service"
)

var (
	Paste    *PasteHandler
	User     *UserHandler
	Stats    *StatsHandler
	ShortURL *ShortURLHandler
)

func InitHandlers(
	ps service.PasteService,
	us service.UserService,
	ss service.StatsService,
	sh service.ShortURLService,
) {
	Paste = NewPasteHandler(ps)
	User = NewUserHandler(us)
	Stats = NewStatsHandler(ss)
	ShortURL = NewShortURLHandler(sh, ps, ss)
}
