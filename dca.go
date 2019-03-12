package hldiscordopus

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

func defaultStreamSettings() *dca.EncodeOptions {
	return &dca.EncodeOptions{
		Volume:        120,
		Channels:      2,
		FrameRate:     48000,
		FrameDuration: 20,
		Bitrate:       96,
		// Should be LowDelay?
		Application:      dca.AudioApplicationAudio,
		CompressionLevel: 5,
		PacketLoss:       1,
		BufferedFrames:   200,
		VBR:              true,
		Threads:          2,
	}
}

type dcaSession struct {
	dca    *dca.EncodeSession
	manage chan AudioAction
	voicec *discordgo.VoiceConnection
}

func (s *dcaSession) playFromDCA() {
	// options are already set in above scope
	defer s.dca.Cleanup()

	streamNaturalEnd := make(chan error)
	dcaStream := dca.NewStream(s.dca, s.voicec, streamNaturalEnd)

	for {
		select {
		case <-streamNaturalEnd:
			s.dca.Stop()
			s.dca.Cleanup()

		case action := <-s.manage:
			switch action {
			case AudioPause:
				dcaStream.SetPaused(true)
			case AudioResume:
				dcaStream.SetPaused(false)
			case AudioStop:
				s.dca.Stop()
				s.dca.Cleanup()
			}
		}
	}
}
