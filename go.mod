module discord-bot

go 1.21

require (
    github.com/bwmarrin/discordgo v0.27.1
    github.com/spf13/viper v1.17.0
    github.com/spf13/cobra v1.7.0
)

// Indirect dependencies will be added automatically when running go mod tidy