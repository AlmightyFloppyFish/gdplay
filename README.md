# High-Level discord OPUS wrapper

## Example

```go
manage := make(chan audio.AudioAction)
vc, err := audio.JoinUserVoiceChannel(m.Author.ID, s)
if err != nil {
    fmt.Println(err)
}
if err := audio.AudioFromYoutubeLink("https://www.youtube.com/watch?v=GX8Hg6kWQYI", vc, manage); err != nil {
    fmt.Println(err)
}

time.Sleep(4 * time.Second)
manage <- audio.AudioPause

time.Sleep(2 * time.Second)
manage <- audio.AudioResume

time.Sleep(2 * time.Second)
manage <- audio.AudioStop
```
