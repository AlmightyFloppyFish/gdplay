# High-Level discord OPUS wrapper

## Example

```go
import audio "github.com/AlmightyFloppyFish/highlevel-discordgo-opus"

-- snip --

manage := make(chan audio.AudioAction)
vc, err := audio.JoinUserVoiceChannel(m.Author.ID, s)
if err != nil {
    // Handle
}
if err := audio.AudioFromYoutubeLink("https://www.youtube.com/watch?v=GX8Hg6kWQYI", vc, manage); err != nil {
    // Handle
}

time.Sleep(4 * time.Second)
manage <- audio.AudioPause

time.Sleep(2 * time.Second)
manage <- audio.AudioResume

time.Sleep(2 * time.Second)
manage <- audio.AudioStop
```
