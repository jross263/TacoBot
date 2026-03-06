package teams

import (
	"tacobot/interactions"

	"github.com/bwmarrin/discordgo"
)

type handleUserSelect struct{}

func (t handleUserSelect) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	Add(i.Member.User.ID, TeamSession{Users: i.MessageComponentData().Values})

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
}

func (t handleUserSelect) ID() string {
	return InteractionHandleUserSelect
}

func init() {
	interactions.Add(handleUserSelect{})
}
