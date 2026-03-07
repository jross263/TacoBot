package teams

import (
	"errors"
	"math/rand"
	"slices"
	"tacobot/util"
)

func GetNumberOfTeams(numberOfUsers int) []int {
	if numberOfUsers < 3 || numberOfUsers > 25 {
		return []int{}
	}

	const minTeams = 2
	maxTeams := max(numberOfUsers-1, minTeams)

	teamSizes := make([]int, (maxTeams-minTeams)+1)
	for i := 0; i < len(teamSizes); i++ {
		teamSizes[i] = i + 2
	}

	return teamSizes
}

func RandomizeTeams(numberOfTeams int, users []string) ([][]string, error) {
	if numberOfTeams < 2 || numberOfTeams > 24 {
		return nil, errors.New("invalid number of teams")
	}

	if len(users) > 25 || len(users) <= numberOfTeams {
		return nil, errors.New("invalid number of users")
	}

	teams := make([][]string, numberOfTeams)

	usersCopy := slices.Clone(users)

	rand.Shuffle(len(usersCopy), func(i, j int) {
		usersCopy[i], usersCopy[j] = usersCopy[j], usersCopy[i]
	})

	teamSize := len(usersCopy) / numberOfTeams
	remainder := len(usersCopy) % numberOfTeams

	for i := 0; i < numberOfTeams; i++ {
		start := i*teamSize + min(remainder, i)
		end := start + teamSize + util.If(i < remainder, 1, 0)
		teams[i] = usersCopy[start:end]
	}

	return teams, nil
}
