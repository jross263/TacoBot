package interactions

import (
	"errors"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func GetUserID(ctx InteractionContext) (string, error) {
	if ctx.Event.GuildID == "" {
		if ctx.Event.User == nil || ctx.Event.User.ID == "" {
			return "", errors.New("unable to find user id")
		}

		return ctx.Event.User.ID, nil
	}

	if ctx.Event.Member == nil || ctx.Event.Member.User == nil || ctx.Event.Member.User.ID == "" {
		return "", errors.New("unable to find user id")
	}

	return ctx.Event.Member.User.ID, nil
}

func RespondWithError(ctx InteractionContext, originalErr error, message string) error {
	respondErr := ctx.Session.InteractionRespond(ctx.Event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})

	if respondErr != nil {
		slog.Error("failed to send error response", "err", respondErr)
		return errors.Join(originalErr, respondErr)
	}

	return originalErr
}
