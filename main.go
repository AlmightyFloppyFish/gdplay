package gdplay

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/rylio/ytdl"
	"io"
)

// AudioSession is used to manage your audio session after starting it
type AudioSession struct {
	manage     chan audioAction
	Vc         *discordgo.VoiceConnection
	dcaSession *dca.EncodeSession
	IsPlaying  bool
	IsPaused   bool
}

// AudioFromRaw allows you to play audio from any source in memory
// however verifying that the data is compatible becomes your job
func AudioFromRaw(src io.Reader, voice *discordgo.VoiceConnection) (*AudioSession, error) {
	session, err := dca.EncodeMem(src, defaultStreamSettings())
	if err != nil {
		return nil, err
	}
	s := AudioSession{
		manage:     make(chan audioAction),
		Vc:         voice,
		dcaSession: session,
		IsPlaying:  true,
		IsPaused:   false,
	}
	go s.playFromDCA()
	return &s, nil
}

// AudioFromFile loads from a filepath
func AudioFromFile(src string, voice *discordgo.VoiceConnection) (*AudioSession, error) {
	session, err := dca.EncodeFile(src, defaultStreamSettings())
	if err != nil {
		return nil, err
	}
	s := AudioSession{
		manage:     make(chan audioAction),
		Vc:         voice,
		dcaSession: session,
		IsPlaying:  true,
		IsPaused:   false,
	}
	go s.playFromDCA()
	return &s, nil
}

// AudioFromYoutubeLink loads from youtube link
func AudioFromYoutubeLink(src string, voice *discordgo.VoiceConnection) (*AudioSession, error) {
	videoInfo, err := ytdl.GetVideoInfo(src)
	if err != nil {
		return nil, fmt.Errorf("Unable to get youtube info from '%s': %s ", src, err.Error())
	}

	reader, err := ytCompatibleStreamFrom(videoInfo)
	if err != nil {
		return nil, err
	}
	// return AudioFromFile(link, voice, manage)
	return AudioFromRaw(reader, voice)
}

// AudioFromYoutubeSearch query's youtube for link and starts stream of the video from the top search
// WARNING: This method of searching can be quite slow. Reason being that I'm avoiding using the youtube API so
// you don't have to add an API key, if you want better speed then get a developer key and use something like
// godoc.org/google.golang.org/api/youtube/v3, then give the youtube link to AudioFromYoutubeLink
func AudioFromYoutubeSearch(search string, voice *discordgo.VoiceConnection) (*AudioSession, error) {
	vids := GetVideosFromSearch(search)
	return vids[0].Play(voice)
}
