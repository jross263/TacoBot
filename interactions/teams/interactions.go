package teams

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"tacobot/interactions"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

var ErrSessionNotFound = errors.New("session not found, please run /teams again")

func handleUsersSelect(s *discordgo.Session, i *discordgo.InteractionCreate, params map[string]string) error {
	cache.Set(i.Member.User.ID, i.MessageComponentData().Values)

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
}

func handleUsersButton(s *discordgo.Session, i *discordgo.InteractionCreate, params map[string]string) error {
	users, ok := cache.Get(i.Member.User.ID)
	if !ok {
		return respondWithError(s, i, ErrSessionNotFound)
	}
	cache.Delete(i.Member.User.ID)

	session := uuid.New().String()
	cache.Set(session, users)

	numberOfTeams := GetNumberOfTeams(len(users))

	var selectOptions []discordgo.SelectMenuOption
	for _, i := range numberOfTeams {
		selectOptions = append(selectOptions, discordgo.SelectMenuOption{
			Label: strconv.Itoa(i),
			Value: strconv.Itoa(i),
		})
	}

	customID, err := interactions.Encode(HandleTeamSelect, interactions.Param{Key: "sessionId", Value: session})
	if err != nil {
		return respondWithError(s, i, err)
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
							CustomID: customID,
						},
					},
				},
			},
		},
	})
}

func handleTeamSelect(s *discordgo.Session, i *discordgo.InteractionCreate, params map[string]string) error {
	sessionId, ok := params["sessionId"]
	if !ok {
		return respondWithError(s, i, ErrSessionNotFound)
	}

	return respondWithTeams(s, i, discordgo.InteractionResponseChannelMessageWithSource, sessionId, i.MessageComponentData().Values[0])
}

func handleTeamButton(s *discordgo.Session, i *discordgo.InteractionCreate, params map[string]string) error {
	sessionId, ok := params["sessionId"]
	if !ok {
		return respondWithError(s, i, ErrSessionNotFound)
	}

	numberOfTeams, ok := params["numTeams"]
	if !ok {
		return respondWithError(s, i, ErrSessionNotFound)
	}

	return respondWithTeams(s, i, discordgo.InteractionResponseUpdateMessage, sessionId, numberOfTeams)
}

func respondWithTeams(s *discordgo.Session, i *discordgo.InteractionCreate, t discordgo.InteractionResponseType, sessionId string, numberOfTeams string) error {
	users, ok := cache.Get(sessionId)
	if !ok {
		return respondWithError(s, i, ErrSessionNotFound)
	}

	n, err := strconv.Atoi(numberOfTeams)
	if err != nil {
		return respondWithError(s, i, err)
	}

	teams, err := RandomizeTeams(n, users)
	if err != nil {
		return respondWithError(s, i, err)
	}

	var sb strings.Builder
	for i, team := range teams {
		sb.WriteString(fmt.Sprintf("Team %d: %s\n", i+1, strings.Join(team, ", ")))
	}

	customID, err := interactions.Encode(HandleTeamButton, interactions.Param{Key: "sessionId", Value: sessionId}, interactions.Param{Key: "numTeams", Value: numberOfTeams})
	if err != nil {
		return respondWithError(s, i, err)
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
							CustomID: customID,
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
