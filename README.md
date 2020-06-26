# 🔊☁️ soundcloud-cli
Simple CLI to fetch the SoundCloud streams

## Getting started

### Build
You can either build it yourself or just download the latest build from [here](https://github.com/timoknapp/soundcloud-cli/actions?query=workflow%3ASoundCloud-CLI). If you want build it yourself, you can use the Makefile with you desired architeture:

```
# build for linux
make build-linux

# build for Raspberry
make build-linux-arm

# build for Mac
make build-mac

# build for Windows
make build-windows
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

## Planned Features

- Add ID3 Tags
- Extend CICD pipeline to build multiple architectures
- Support multiple tracks download


## Contribution

PRs Welcome :)
