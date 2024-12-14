package commands

import (
	"github.com/nextep-community/gocord/discord"
	"github.com/nextep-community/gocord/handler"
)

var pingCommand = discord.SlashCommandCreate{
	Name:        "ping",
	Description: "Replies with pong",
}

func HandlePing(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return e.CreateMessage(discord.MessageCreate{
		Content: "pong",
		Components: []discord.ContainerComponent{
			discord.ActionRowComponent{
				discord.NewPrimaryButton("Test", "/test"),
			},
		},
	})
}
