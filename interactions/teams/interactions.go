package teams

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var ErrSessionNotFound = errors.New("session not found, please run /teams again")

func handleUsersSelect(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	cache.Set(i.Member.User.ID, TeamSession{Users: i.MessageComponentData().Values})

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
}

func handleUsersButton(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	users, ok := cache.Get(i.Member.User.ID)
	if !ok {
		return respondWithError(s, i, ErrSessionNotFound)
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
							CustomID: HandleTeamSelect,
						},
					},
				},
			},
		},
	})
}

func handleTeamSelect(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	var numberOfTeams = i.MessageComponentData().Values[0]
	n, err := strconv.Atoi(numberOfTeams)
	if err != nil {
		return respondWithError(s, i, err)
	}

	err = cache.Update(i.Member.User.ID, func(ts TeamSession) TeamSession {
		ts.NumberOfTeams = n
		return ts
	})
	if err != nil {
		return respondWithError(s, i, err)
	}

	return respondWithTeams(s, i, discordgo.InteractionResponseChannelMessageWithSource)
}

func handleTeamButton(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return respondWithTeams(s, i, discordgo.InteractionResponseUpdateMessage)
}

func respondWithTeams(s *discordgo.Session, i *discordgo.InteractionCreate, t discordgo.InteractionResponseType) error {
	users, ok := cache.Get(i.Member.User.ID)
	if !ok {
		return respondWithError(s, i, ErrSessionNotFound)
	}

	teams, err := RandomizeTeams(users.NumberOfTeams, users.Users)
	if err != nil {
		return respondWithError(s, i, err)
	}

	var sb strings.Builder
	for i, team := range teams {
		sb.WriteString(fmt.Sprintf("Team %d: %s\n", i+1, strings.Join(team, ", ")))
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: t,
		Data: &discordgo.InteractionResponseData{
			Content: sb.String(),
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Shuffle!",
							CustomID: HandleTeamButton,
						},
					},
				},
			},
		},
	})
}

func respondWithError(s *discordgo.Session, i *discordgo.InteractionCreate, originalErr error) error {
	respondErr := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: "Error encountered, please run /teams again!",
		},
	})

	if respondErr != nil {
		slog.Error("failed to send error response", "err", respondErr)
		return errors.Join(originalErr, respondErr)
	}

	return originalErr
}
