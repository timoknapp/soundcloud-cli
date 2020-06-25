# 🔊☁️ soundcloud-cli
Simple CLI to fetch the SoundCloud streams

## Getting started

### Build
```
go build -v -o sc-cli cmd/main.go
```

### Info
```
./sc-cli -h
Usage of ./sc-cli:
  -downloadPath string
    	the path you want to download the tracks to (default "download")
  -downloadQuality string
    	the quality of the music files (mp3/ogg) (default "mp3")
  -trackID string
    	the SoundCloud track ID you want to download
  -version
    	prints current version
```

### Run
```
./sc-cli -trackID=${TRACK_ID}
```


## Contribution

PRs Welcome :)
