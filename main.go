package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"tacobot/interactions"
	_ "tacobot/interactions/teams"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

type config struct {
	appID   string
	guildID string
	token   string
}

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	_ = godotenv.Load()

	config, err := loadConfig()
	if err != nil {
		return err
	}

	session, err := discordgo.New("Bot " + config.token)
	if err != nil {
		return fmt.Errorf("creating session: %w", err)
	}

	slog.Info("Attempting to start bot session...")

	registerInteractionHandlers(session)
	registeredCommands := registerCommands(session, config)

	err = session.Open()
	if err != nil {
		return fmt.Errorf("opening session: %w", err)
	}

	slog.Info("TacoBot running!")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	removeCommands(session, registeredCommands, config)

	slog.Info("Attempting to stop bot session...")

	err = session.Close()
	if err != nil {
		return fmt.Errorf("closing session: %w", err)
	}

	slog.Info("TacoBot stopped!")

	return nil
}

func loadConfig() (config, error) {
	appID := os.Getenv("DISCORD_APP_ID")
	guildID := os.Getenv("DISCORD_GUILD_ID")
	token := os.Getenv("DISCORD_TOKEN")

	var missing []string
	if strings.TrimSpace(appID) == "" {
		missing = append(missing, "DISCORD_APP_ID")
	}
	if strings.TrimSpace(guildID) == "" {
		missing = append(missing, "DISCORD_GUILD_ID")
	}
	if strings.TrimSpace(token) == "" {
		missing = append(missing, "DISCORD_TOKEN")
	}

	if len(missing) > 0 {
		return config{}, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}
	return config{appID: appID, guildID: guildID, token: token}, nil
}

func registerInteractionHandlers(s *discordgo.Session) {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		slog.Info("Received interaction", "Type", i.Type)
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			routeInteractionType(i.ApplicationCommandData().Name, s, i)
		case discordgo.InteractionMessageComponent:
			routeInteractionType(i.MessageComponentData().CustomID, s, i)
		}
	})
}

func routeInteractionType(id string, s *discordgo.Session, i *discordgo.InteractionCreate) {
	id, params, err := interactions.Decode(id)
	if err != nil {
		slog.Error("Error executing handler", "ID", id, "err", err)
		return
	}

	if h, ok := interactions.Get(id); ok {
		slog.Info("Executing handler", "ID", id)
		ctx := interactions.InteractionContext{
			Session: interactions.NewDiscordSession(s),
			Event:   i,
		}
		err := h(ctx, params)
		if err != nil {
			slog.Error("Error executing handler", "ID", id, "err", err)
		}
	}
}

func registerCommands(s *discordgo.Session, config config) []*discordgo.ApplicationCommand {
	commands := interactions.AllCommands()

	registeredCommands := make([]*discordgo.ApplicationCommand, 0, len(commands))

	for _, cmd := range commands {
		slog.Info("Attempting to add command", "Name", cmd.Name)
		createdCmd, err := s.ApplicationCommandCreate(config.appID, config.guildID, cmd)
		if err != nil {
			slog.Error("Error adding command", "Name", cmd.Name, "err", err)
			continue
		}
		slog.Info("Added command", "Name", cmd.Name)
		registeredCommands = append(registeredCommands, createdCmd)
	}

	return registeredCommands
}

func removeCommands(s *discordgo.Session, registeredCommands []*discordgo.ApplicationCommand, config config) {
	cmds, err := s.ApplicationCommands(config.appID, config.guildID)

	if err != nil {
		slog.Warn("Error retrieving commands, defaulting to registered list", "err", err)
		cmds = registeredCommands
	}

	for _, cmd := range cmds {
		slog.Info("Attempting to remove command", "Name", cmd.Name, "ID", cmd.ID)
		err := s.ApplicationCommandDelete(config.appID, config.guildID, cmd.ID)
		if err != nil {
			slog.Error("Error trying to delete command", "Name", cmd.Name, "ID", cmd.ID, "err", err)
			continue
		}
		slog.Info("Removed command", "Name", cmd.Name)
	}
}
