package teams

import (
	"fmt"
	"strconv"
	"strings"
	"tacobot/interactions"

	"github.com/bwmarrin/discordgo"
)

type handleTeamSelect struct{}

func (t handleTeamSelect) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	var numberOfTeams = i.MessageComponentData().Values[0]
	n, _ := strconv.Atoi(numberOfTeams)

	Update(i.Member.User.ID, n)
	users, _ := Get(i.Member.User.ID)

	teams, _ := RandomizeTeams(n, users.Users)

	var sb strings.Builder
	for i, team := range teams {
		sb.WriteString(fmt.Sprintf("Team %d: %s\n", i+1, strings.Join(team, ", ")))
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: sb.String(),
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Shuffle!",
							CustomID: InteractionhandleTeamButton,
						},
					},
				},
			},
		},
	})
}

func (t handleTeamSelect) ID() string {
	return InteractionHandleTeamSelect
}

func init() {
	interactions.Add(handleTeamSelect{})
}
