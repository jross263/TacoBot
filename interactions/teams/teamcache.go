package teams

import "tacobot/util"

type TeamSession struct {
	Users         []string
	NumberOfTeams int
}

var cache = util.NewCache[string, TeamSession]("TeamCache")
