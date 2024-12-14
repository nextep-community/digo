package commands

import (
	"github.com/nextep-community/gocord/discord"
)

var Commands = []discord.ApplicationCommandCreate{
	pingCommand,
	playCommand,
}
