package handlers

import (
	"github.com/GritsyukLeonid/pastebin-go/internal/service"
)

var (
	pasteService    service.PasteService
	userService     service.UserService
	statsService    service.StatsService
	shortURLService service.ShortURLService
)

func SetServices(
	ps service.PasteService,
	us service.UserService,
	ss service.StatsService,
	sh service.ShortURLService,
) {
	pasteService = ps
	userService = us
	statsService = ss
	shortURLService = sh
}
