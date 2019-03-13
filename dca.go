package gdplay

import (
	"fmt"
	// "github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

const debug = false

func defaultStreamSettings() *dca.EncodeOptions {
	return &dca.EncodeOptions{
		Volume:        120,
		Channels:      2,
		FrameRate:     48000,
		FrameDuration: 20,
		Bitrate:       96,
		// Should be LowDelay?
		Application:      dca.AudioApplicationAudio,
		CompressionLevel: 1,
		PacketLoss:       1,
		BufferedFrames:   200,
		VBR:              true,
		Threads:          4,
	}
}

func (s *AudioSession) playFromDCA() {
	// options are already set in above scope
	defer s.dcaSession.Cleanup()

	streamNaturalEnd := make(chan error)
	dcaStream := dca.NewStream(s.dcaSession, s.Vc, streamNaturalEnd)

	for {
		select {
		case e := <-streamNaturalEnd:
			if debug {
				fmt.Println(e)
			}
			s.dcaSession.Stop()
			s.dcaSession.Cleanup()
			s.IsPlaying = false
			return

		case action := <-s.manage:
			switch action {
			case audioPause:
				dcaStream.SetPaused(true)
				s.IsPaused = true
			case audioResume:
				dcaStream.SetPaused(false)
				s.IsPaused = false
			case audioStop:
				s.dcaSession.Stop()
				s.dcaSession.Cleanup()
				s.IsPlaying = false
				return
			}
		}
	}
}
