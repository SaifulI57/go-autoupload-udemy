package messager

import (
	"fmt"
	"os"
	"os/signal"
	"strings"

	log "github.com/SaifulI57/uploader-udemy/logger"
	"github.com/bwmarrin/discordgo"
)

type store struct {
	ChannelID string
}

var v []*store

func Start() {
	bot, err := discordgo.New("Bot " + os.Getenv("Token"))
	if err != nil {
		log.Logger.Info(fmt.Sprintf("Error creating session bot: %s", err))
		return
	}
	bot.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Println("Bot is ready")
	})
	bot.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {

		if m.Author.ID == s.State.User.ID {
			return
		}
		fmt.Println(m.Content)
		if strings.Contains(m.Content, "ping") {
			s.ChannelMessageSend(m.ChannelID, "/thread name:start message:work!!")
		}
		splited := strings.Split(m.Content, " ")
		if strings.Contains(splited[0], "!set") && strings.Contains(splited[1], "subscribe") {
			newS := &store{
				ChannelID: m.ChannelID,
			}
			v = append(v, newS)
			s.ChannelMessageSend(m.ChannelID, v[0].ChannelID)
		}
		if ch, err := s.State.Channel(m.ChannelID); err != nil || !ch.IsThread() {
			if strings.Contains(m.Content, "pong") {
				thread, err := s.MessageThreadStartComplex(m.ChannelID, m.ID, &discordgo.ThreadStart{
					Name:                "Pong game with ",
					AutoArchiveDuration: 60,
					Invitable:           false,
				})
				if err != nil {
					panic(err)
				}
				_, _ = s.ChannelMessageSend(thread.ID, "pong")
			}
		} else {
			_, _ = s.ChannelMessageSendReply(m.ChannelID, "pong", m.Reference())
		}

	})
	bot.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
	err = bot.Open()
	if err != nil {
		log.Logger.Error(fmt.Sprintf("Error opening session bot: %s", err))
	}
	defer bot.Close()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc
}
