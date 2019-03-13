package hldiscordopus

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/rylio/ytdl"
	"io"
)

// Signals for manage channel
const (
	AudioStop = iota
	AudioPause
	AudioResume
)

// AudioAction is a Signal for a manage channel
type AudioAction int

// AudioFromRaw allows you to play audio from any source in memory
// however verifying that the data is compatible becomes your job
func AudioFromRaw(src io.Reader, voice *discordgo.VoiceConnection, manage chan AudioAction) error {
	session, err := dca.EncodeMem(src, defaultStreamSettings())
	if err != nil {
		return err
	}
	s := dcaSession{session, manage, voice}
	go s.playFromDCA()
	return nil
}

// AudioFromFile loads from a filepath
func AudioFromFile(src string, voice *discordgo.VoiceConnection, manage chan AudioAction) error {
	session, err := dca.EncodeFile(src, defaultStreamSettings())
	if err != nil {
		return err
	}
	s := dcaSession{session, manage, voice}
	go s.playFromDCA()
	return nil
}

// AudioFromYoutubeLink loads from youtube link
func AudioFromYoutubeLink(src string, voice *discordgo.VoiceConnection, manage chan AudioAction) error {
	videoInfo, err := ytdl.GetVideoInfo(src)
	if err != nil {
		return fmt.Errorf("Unable to get youtube info from '%s': %s ", src, err.Error())
	}

	reader, err := ytCompatibleStreamFrom(videoInfo)
	if err != nil {
		return err
	}
	// return AudioFromFile(link, voice, manage)
	return AudioFromRaw(reader, voice, manage)
}

// AudioFromYoutubeSearch query's youtube for link
// and starts stream of the video from the top search
func AudioFromYoutubeSearch(search string, voice *discordgo.VoiceConnection, manage chan AudioAction) {
	vids := GetVideosFromSearch(search)
	vids[0].Play(voice, manage)
}
