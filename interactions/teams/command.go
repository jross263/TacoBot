package teams

import (
	"log/slog"
	"tacobot/interactions"
	"tacobot/util"

	"github.com/bwmarrin/discordgo"
)

type SessionStore interface {
	Get(key string) ([]string, bool)
	Set(key string, value []string)
	Delete(key string)
}

type TeamsHandlers struct {
	store SessionStore
}

const userSelectMinimum int = 3
const userSelectMaximum int = 25

func (h *TeamsHandlers) handleTeamCommand(ctx interactions.InteractionContext, params map[string]string) error {
	defaultMembers := getDefaultMembers(ctx)
	h.store.Set(ctx.Event.Member.User.ID, defaultMembers)
	defaultSelectValues := getDefaultSelectValues(defaultMembers)

	return ctx.Session.InteractionRespond(ctx.Event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pick your players",
			Flags:   discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							MenuType:      discordgo.UserSelectMenu,
							MinValues:     util.Ptr(userSelectMinimum),
							MaxValues:     userSelectMaximum,
							CustomID:      HandleUsersSelect,
							DefaultValues: defaultSelectValues,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Go!",
							CustomID: HandleUsersButton,
							Style:    discordgo.SuccessButton,
						},
					},
				},
			},
		},
	})
}

func getDefaultMembers(ctx interactions.InteractionContext) []string {
	var members []string

	vs, err := ctx.Session.StateVoiceState(ctx.Event.GuildID, ctx.Event.Member.User.ID)
	if err != nil {
		slog.Warn("Error getting voice state, defaulting to no users", "err", err)
		return members
	}

	guild, err := ctx.Session.StateGuild(ctx.Event.GuildID)
	if err != nil {
		slog.Warn("Error getting voice state, defaulting to no users", "err", err)
		return members
	}
	for _, v := range guild.VoiceStates {
		if v.ChannelID == vs.ChannelID {
			members = append(members, v.UserID)
			if len(members) == 25 {
				break
			}
		}
	}

	return members
}

func getDefaultSelectValues(users []string) []discordgo.SelectMenuDefaultValue {
	var members []discordgo.SelectMenuDefaultValue

	for _, user := range users {
		members = append(members, discordgo.SelectMenuDefaultValue{
			ID:   user,
			Type: discordgo.SelectMenuDefaultValueUser,
		})
	}

	return members
}

func init() {
	handlers := &TeamsHandlers{store: util.NewCache[string, []string]("TeamCache")}

	interactions.Add(interactions.Command{
		Definition: &discordgo.ApplicationCommand{
			Name:        TeamsCommand,
			Description: "Test Command",
		},
		Handle: handlers.handleTeamCommand,
		Interactions: []interactions.Handler{
			{ID: HandleUsersSelect, Handle: handlers.handleUsersSelect},
			{ID: HandleUsersButton, Handle: handlers.handleUsersButton},
			{ID: HandleTeamSelect, Handle: handlers.handleTeamSelect},
			{ID: HandleTeamButton, Handle: handlers.handleTeamButton},
		},
	})
}
