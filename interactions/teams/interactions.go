package teams

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func handleUsersSelect(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	Add(i.Member.User.ID, TeamSession{Users: i.MessageComponentData().Values})

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
}

func handleUsersButton(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	users, ok := Get(i.Member.User.ID)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content: "Error",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
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
	n, _ := strconv.Atoi(numberOfTeams)

	Update(i.Member.User.ID, n)
	users, _ := Get(i.Member.User.ID)

	teams, _ := RandomizeTeams(n, users.Users)

	var sb strings.Builder
	for i, team := range teams {
		sb.WriteString(fmt.Sprintf("Team %d: %s\n", i+1, strings.Join(team, ", ")))
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
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

func handleTeamButton(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	users, _ := Get(i.Member.User.ID)
	teams, _ := RandomizeTeams(users.NumberOfTeams, users.Users)

	var sb strings.Builder
	for i, team := range teams {
		sb.WriteString(fmt.Sprintf("Team %d: %s\n", i+1, strings.Join(team, ", ")))
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
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
