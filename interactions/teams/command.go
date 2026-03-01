package teams

import (
	"tacobot/interactions"

	"github.com/bwmarrin/discordgo"
)

type Teams struct{}

func (t Teams) ID() string {
	return CommandTeams
}

func (t Teams) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Yo",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Yes",
							Style:    discordgo.SuccessButton,
							CustomID: InteractionPickButton,
						},
					},
				},
			},
		},
	})
}

func (t Teams) Definition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        CommandTeams,
		Description: "Test Command",
	}
}

func init() {
	interactions.Add(Teams{})
}
