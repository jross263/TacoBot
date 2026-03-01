package commands

import "github.com/bwmarrin/discordgo"

type Teams struct{}

func (t Teams) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Yo",
		},
	})
}

func (t Teams) Definition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "testcommand",
		Description: "Test Command",
	}
}

func init() {
	Add(Teams{})
}
