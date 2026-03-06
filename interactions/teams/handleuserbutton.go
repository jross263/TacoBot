package teams

import (
	"strconv"
	"tacobot/interactions"

	"github.com/bwmarrin/discordgo"
)

type handleUserbutton struct{}

// Handle implements [interactions.Handler].
func (h handleUserbutton) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	users, ok := Get(i.Member.User.ID)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content: "Error",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	}

	numberOfTeams := GetNumberOfTeams(len(users.Users))

	var selectOptions []discordgo.SelectMenuOption
	for _, i := range numberOfTeams {
		selectOptions = append(selectOptions, discordgo.SelectMenuOption{
			Label: strconv.Itoa(i),
			Value: strconv.Itoa(i),
		})
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: "How many teams?",
			Flags:   discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							MenuType: discordgo.StringSelectMenu,
							Options:  selectOptions,
							CustomID: InteractionHandleTeamSelect,
						},
					},
				},
			},
		},
	})
}

// ID implements [interactions.Handler].
func (h handleUserbutton) ID() string {
	return InteractionHandleUserButton
}

func init() {
	interactions.Add(handleUserbutton{})
}
