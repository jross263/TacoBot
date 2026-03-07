package interactions

import "github.com/bwmarrin/discordgo"

type InteractionResponder interface {
	InteractionRespond(interaction *discordgo.Interaction, resp *discordgo.InteractionResponse, options ...discordgo.RequestOption) error
	StateVoiceState(guildID, userID string) (*discordgo.VoiceState, error)
	StateGuild(guildID string) (*discordgo.Guild, error)
}

type InteractionContext struct {
	Session InteractionResponder
	Event   *discordgo.InteractionCreate
}

type HandleFunc func(InteractionContext, map[string]string) error

type discordSession struct {
	s *discordgo.Session
}

func (d *discordSession) InteractionRespond(i *discordgo.Interaction, r *discordgo.InteractionResponse, opts ...discordgo.RequestOption) error {
	return d.s.InteractionRespond(i, r, opts...)
}

func (d *discordSession) StateVoiceState(guildID, userID string) (*discordgo.VoiceState, error) {
	return d.s.State.VoiceState(guildID, userID)
}

func (d *discordSession) StateGuild(guildID string) (*discordgo.Guild, error) {
	return d.s.State.Guild(guildID)
}

func NewDiscordSession(s *discordgo.Session) InteractionResponder {
	return &discordSession{s: s}
}
