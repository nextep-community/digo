package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
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
