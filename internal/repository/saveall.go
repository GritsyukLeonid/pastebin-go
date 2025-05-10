package repository

func SaveAll() {
	saveJSON("pastes.json", Pastes)
	saveJSON("users.json", Users)
	saveJSON("stats.json", StatsSet)
	saveJSON("urls.json", URLs)
}
