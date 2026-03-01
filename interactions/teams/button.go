package teams

import (
	"tacobot/interactions"

	"github.com/bwmarrin/discordgo"
)

type Button struct{}

func (t Button) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Great! If you wanna know more or just have questions, feel free to visit Discord Devs and Discord Gophers server. " +
				"But now, when you know how buttons work, let's move onto select menus (execute `/selects single`)",
			Flags: discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Emoji: &discordgo.ComponentEmoji{
								Name: "🔧",
							},
							Label: "Discord developers",
							Style: discordgo.LinkButton,
							URL:   "https://discord.gg/discord-developers",
						},
						discordgo.Button{
							Emoji: &discordgo.ComponentEmoji{
								Name: "🦫",
							},
							Label: "Discord Gophers",
							Style: discordgo.LinkButton,
							URL:   "https://discord.gg/7RuRrVHyXF",
						},
					},
				},
			},
		},
	})
}

func (t Button) ID() string {
	return InteractionPickButton
}

func init() {
	interactions.Add(Button{})
}
