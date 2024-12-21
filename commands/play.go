package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
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
