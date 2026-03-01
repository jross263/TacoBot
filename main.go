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
	appId   string
	guildId string
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
	appId := os.Getenv("DISCORD_APP_ID")
	guildId := os.Getenv("DISCORD_GUILD_ID")
	token := os.Getenv("DISCORD_TOKEN")

	var missing []string
	if strings.TrimSpace(appId) == "" {
		missing = append(missing, "DISCORD_APP_ID")
	}
	if strings.TrimSpace(guildId) == "" {
		missing = append(missing, "DISCORD_GUILD_ID")
	}
	if strings.TrimSpace(token) == "" {
		missing = append(missing, "DISCORD_TOKEN")
	}

	if len(missing) > 0 {
		return config{}, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}
	return config{appId: appId, guildId: guildId, token: token}, nil
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
	if h, ok := interactions.Get(id); ok {
		slog.Info("Executing handler", "ID", h.ID())
		err := h.Handle(s, i)
		if err != nil {
			slog.Error("Error executing handler", "ID", h.ID(), "err", err)
		}
	}
}

func registerCommands(s *discordgo.Session, config config) []*discordgo.ApplicationCommand {
	commands := interactions.AllCommands()

	registeredCommands := make([]*discordgo.ApplicationCommand, 0, len(commands))

	for _, cmd := range commands {
		slog.Info("Attempting to add command", "Name", cmd.Definition().Name)
		createdCmd, err := s.ApplicationCommandCreate(config.appId, config.guildId, cmd.Definition())
		if err != nil {
			slog.Error("Error adding command", "Name", cmd.Definition().Name, "err", err)
			continue
		}
		slog.Info("Added command", "Name", cmd.Definition().Name)
		registeredCommands = append(registeredCommands, createdCmd)
	}

	return registeredCommands
}

func removeCommands(s *discordgo.Session, registeredCommands []*discordgo.ApplicationCommand, config config) {
	cmds, err := s.ApplicationCommands(config.appId, config.guildId)

	if err != nil {
		slog.Warn("Error retrieving commands, defaulting to registered list", "err", err)
		cmds = registeredCommands
	}

	for _, cmd := range cmds {
		slog.Info("Attempting to remove command", "Name", cmd.Name, "ID", cmd.ID)
		err := s.ApplicationCommandDelete(config.appId, config.guildId, cmd.ID)
		if err != nil {
			slog.Error("Error trying to delete command", "Name", cmd.Name, "ID", cmd.ID, "err", err)
			continue
		}
		slog.Info("Removed command", "Name", cmd.Name)
	}
}
