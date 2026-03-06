package teams

type TeamSession struct {
	Users         []string
	NumberOfTeams int
}

var cache = make(map[string]TeamSession)

func Get(id string) (TeamSession, bool) {
	teamSession, ok := cache[id]
	return teamSession, ok
}

func Add(id string, users TeamSession) {
	cache[id] = users
}

func Update(id string, numberOfTeams int) {
	session := cache[id]
	session.NumberOfTeams = numberOfTeams
	cache[id] = session
}
