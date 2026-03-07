package teams

import (
	"log/slog"
	"tacobot/interactions"
	"tacobot/util"

	"github.com/bwmarrin/discordgo"
)

const userSelectMinimum int = 3
const userSelectMaximum int = 25

func handleTeamCommand(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	defaultMembers := getDefaultMembers(s, i)
	Add(i.Member.User.ID, TeamSession{Users: defaultMembers})
	defaultSelectValues := getDefaultSelectValues(defaultMembers)

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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

func getDefaultMembers(s *discordgo.Session, i *discordgo.InteractionCreate) []string {
	var members []string

	vs, err := s.State.VoiceState(i.GuildID, i.Member.User.ID)
	if err != nil {
		slog.Warn("Error getting voice state, defaulting to no users", "err", err)
		return members
	}

	guild, err := s.State.Guild(i.GuildID)
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
	interactions.Add(interactions.Command{
		Definition: &discordgo.ApplicationCommand{
			Name:        TeamsCommand,
			Description: "Test Command",
		},
		Handle: handleTeamCommand,
		Interactions: []interactions.Handler{
			{ID: HandleUsersSelect, Handle: handleUsersSelect},
			{ID: HandleUsersButton, Handle: handleUsersButton},
			{ID: HandleTeamSelect, Handle: handleTeamSelect},
			{ID: HandleTeamButton, Handle: handleTeamButton},
		},
	})
}
