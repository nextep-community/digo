package commands

import (
	"github.com/nextep-community/gocord/discord"
	"github.com/nextep-community/gocord/handler"
)

var playCommand = discord.SlashCommandCreate{
	Name:        "play",
	Description: "teste",
}

func PlayCommandHandler(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return e.CreateMessage(discord.MessageCreate{
		Content: "Play",
	})
}
