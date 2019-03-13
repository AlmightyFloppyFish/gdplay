package main

import (
	"fmt"
	"github.com/AlmightyFloppyFish/gdplay"
	"github.com/bwmarrin/discordgo"
	"os"
	"strings"
	"time"
)

var isPlaying = false

func main() {
	token := os.Getenv("TOKEN")
	if len(token) == 0 {
		fmt.Println("Please provide a discord token by setting the TOKEN environment variable")
		return
	}
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	if dg.Open() != nil {
		panic("Could not open discord session")
	}
	defer dg.Close()

	dg.AddHandler(onMessage)

	select {}
}

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	words := strings.Split(m.Content, " ")
	if []byte(m.Content)[0] != byte('!') {
		return
	}

	var (
		playError error
		session   *gdplay.AudioSession
	)

	if isPlaying {
		s.ChannelMessageSend(m.ChannelID, "Already playing!")
		return
	}

	switch words[0] {
	case "!link":
		if len(words) != 2 {
			return
		}
		vc, err := gdplay.JoinUserVoiceChannel(m.Author.ID, s)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Could not join your voice channel")
			return
		}
		session, playError = gdplay.AudioFromYoutubeLink(words[1], vc)
	case "!search":
		if len(words) < 2 {
			return
		}
		vc, err := gdplay.JoinUserVoiceChannel(m.Author.ID, s)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Could not join your voice channel")
			return
		}
		search := strings.Join(words[1:], " ")
		fmt.Println(search)
		vids := gdplay.GetVideosFromSearch(search)
		for _, v := range vids {
			fmt.Println(v.Title)
		}
		s.ChannelMessageSend(m.ChannelID, "Playing "+vids[0].Title)
		session, playError = vids[0].Play(vc)
	case "!oof":
		vc, err := gdplay.JoinUserVoiceChannel(m.Author.ID, s)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Could not join your voice channel")
			return
		}
		session, playError = gdplay.AudioFromFile("assets/oof.mp3", vc)

		session.Wait()
		session.Vc.Disconnect()
		session.Vc.Close()
		return
	}
	if playError != nil {
		s.ChannelMessageSend(m.ChannelID, "Failed to play audio: "+playError.Error())
		return
	}

	isPlaying = true
	defer func() { isPlaying = false }()

	// Wait for 10 seconds then stops audio and exit
	time.Sleep(10 * time.Second)
	if session.IsPlaying {
		session.Stop(true)
	}
}
