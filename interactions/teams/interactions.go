package teams

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"tacobot/interactions"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

var ErrSessionNotFound = errors.New("session not found")
var ErrNumberOfTeamsNotFound = errors.New("numTeams not found")

func (h *TeamsHandlers) handleUsersSelect(ctx interactions.InteractionContext, params map[string]string) error {
	userId, err := interactions.GetUserID(ctx)
	if err != nil {
		return interactions.RespondWithError(ctx, err, genericError)
	}

	h.store.Set(userId, ctx.Event.MessageComponentData().Values)

	return ctx.Session.InteractionRespond(ctx.Event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
}

func (h *TeamsHandlers) handleUsersButton(ctx interactions.InteractionContext, params map[string]string) error {
	userId, err := interactions.GetUserID(ctx)
	if err != nil {
		return interactions.RespondWithError(ctx, err, genericError)
	}

	users, ok := h.store.Get(userId)
	if !ok {
		return interactions.RespondWithError(ctx, ErrSessionNotFound, genericError)
	}
	h.store.Delete(userId)

	session := uuid.New().String()
	h.store.Set(session, users)

	numberOfTeams := GetNumberOfTeams(len(users))

	var selectOptions []discordgo.SelectMenuOption
	for _, i := range numberOfTeams {
		selectOptions = append(selectOptions, discordgo.SelectMenuOption{
			Label: strconv.Itoa(i),
			Value: strconv.Itoa(i),
		})
	}

	customID, err := interactions.Encode(HandleTeamSelect, interactions.Param{Key: "sessionID", Value: session})
	if err != nil {
		return interactions.RespondWithError(ctx, err, genericError)
	}

	return ctx.Session.InteractionRespond(ctx.Event.Interaction, &discordgo.InteractionResponse{
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

func (h *TeamsHandlers) handleTeamSelect(ctx interactions.InteractionContext, params map[string]string) error {
	sessionID, ok := params["sessionID"]
	if !ok {
		return interactions.RespondWithError(ctx, ErrSessionNotFound, genericError)
	}

	values := ctx.Event.MessageComponentData().Values
	if len(values) != 1 {
		return errors.New("no values in selection")
	}

	return h.respondWithTeams(ctx, discordgo.InteractionResponseChannelMessageWithSource, sessionID, values[0])
}

func (h *TeamsHandlers) handleTeamButton(ctx interactions.InteractionContext, params map[string]string) error {
	sessionID, ok := params["sessionID"]
	if !ok {
		return interactions.RespondWithError(ctx, ErrSessionNotFound, genericError)
	}

	numberOfTeams, ok := params["numTeams"]
	if !ok {
		return interactions.RespondWithError(ctx, ErrNumberOfTeamsNotFound, genericError)
	}

	return h.respondWithTeams(ctx, discordgo.InteractionResponseUpdateMessage, sessionID, numberOfTeams)
}

func (h *TeamsHandlers) respondWithTeams(ctx interactions.InteractionContext, t discordgo.InteractionResponseType, sessionID string, numberOfTeams string) error {
	users, ok := h.store.Get(sessionID)
	if !ok {
		return interactions.RespondWithError(ctx, ErrSessionNotFound, genericError)
	}

	n, err := strconv.Atoi(numberOfTeams)
	if err != nil {
		return interactions.RespondWithError(ctx, err, genericError)
	}

	teams, err := RandomizeTeams(n, users)
	if err != nil {
		return interactions.RespondWithError(ctx, err, genericError)
	}

	var sb strings.Builder
	for i, team := range teams {
		sb.WriteString(fmt.Sprintf("Team %d: ", i+1))
		for j, name := range team {
			if j > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString("@")
			sb.WriteString(name)
		}
		sb.WriteByte('\n')
	}

	customID, err := interactions.Encode(HandleTeamButton, interactions.Param{Key: "sessionID", Value: sessionID}, interactions.Param{Key: "numTeams", Value: numberOfTeams})
	if err != nil {
		return interactions.RespondWithError(ctx, err, genericError)
	}

	return ctx.Session.InteractionRespond(ctx.Event.Interaction, &discordgo.InteractionResponse{
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
