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

func (h *TeamsHandlers) handleUsersSelect(ctx interactions.InteractionContext, params map[string]string) error {
	h.store.Set(ctx.Event.Member.User.ID, ctx.Event.MessageComponentData().Values)

	return ctx.Session.InteractionRespond(ctx.Event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
}

func (h *TeamsHandlers) handleUsersButton(ctx interactions.InteractionContext, params map[string]string) error {
	users, ok := h.store.Get(ctx.Event.Member.User.ID)
	if !ok {
		return respondWithError(ctx, ErrSessionNotFound)
	}
	h.store.Delete(ctx.Event.Member.User.ID)

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
		return respondWithError(ctx, err)
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
		return respondWithError(ctx, ErrSessionNotFound)
	}

	return h.respondWithTeams(ctx, discordgo.InteractionResponseChannelMessageWithSource, sessionID, ctx.Event.MessageComponentData().Values[0])
}

func (h *TeamsHandlers) handleTeamButton(ctx interactions.InteractionContext, params map[string]string) error {
	sessionID, ok := params["sessionID"]
	if !ok {
		return respondWithError(ctx, ErrSessionNotFound)
	}

	numberOfTeams, ok := params["numTeams"]
	if !ok {
		return respondWithError(ctx, ErrSessionNotFound)
	}

	return h.respondWithTeams(ctx, discordgo.InteractionResponseUpdateMessage, sessionID, numberOfTeams)
}

func (h *TeamsHandlers) respondWithTeams(ctx interactions.InteractionContext, t discordgo.InteractionResponseType, sessionID string, numberOfTeams string) error {
	users, ok := h.store.Get(sessionID)
	if !ok {
		return respondWithError(ctx, ErrSessionNotFound)
	}

	n, err := strconv.Atoi(numberOfTeams)
	if err != nil {
		return respondWithError(ctx, err)
	}

	teams, err := RandomizeTeams(n, users)
	if err != nil {
		return respondWithError(ctx, err)
	}

	var sb strings.Builder
	for i, team := range teams {
		sb.WriteString(fmt.Sprintf("Team %d: %s\n", i+1, strings.Join(team, ", ")))
	}

	customID, err := interactions.Encode(HandleTeamButton, interactions.Param{Key: "sessionID", Value: sessionID}, interactions.Param{Key: "numTeams", Value: numberOfTeams})
	if err != nil {
		return respondWithError(ctx, err)
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

func respondWithError(ctx interactions.InteractionContext, originalErr error) error {
	respondErr := ctx.Session.InteractionRespond(ctx.Event.Interaction, &discordgo.InteractionResponse{
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
