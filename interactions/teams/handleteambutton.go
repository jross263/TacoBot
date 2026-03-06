package teams

import (
	"fmt"
	"strings"
	"tacobot/interactions"

	"github.com/bwmarrin/discordgo"
)

type handleTeamButton struct{}

func (t handleTeamButton) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	users, _ := Get(i.Member.User.ID)
	teams, _ := RandomizeTeams(users.NumberOfTeams, users.Users)

	var sb strings.Builder
	for i, team := range teams {
		sb.WriteString(fmt.Sprintf("Team %d: %s\n", i+1, strings.Join(team, ", ")))
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
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

func (t handleTeamButton) ID() string {
	return InteractionhandleTeamButton
}

func init() {
	interactions.Add(handleTeamButton{})
}
