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
NAME:
   SoundCloud CLI - A simple CLI to interact with tracks on SoundCloud

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
   download, dl  Download a track
   meta, m       Show metadata for a track
   help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)

```

### Run
See how to use:
```
./sc-cli download -h
NAME:
   main download - Download a track

USAGE:
   main download [command options] [arguments...]

OPTIONS:
   --path value     Path where the files will be stored (default: "download")
   --quality value  Quality of the track (default: "mp3")
   --help, -h       show help (default: false)


# By ID >
./sc-cli download $TRACK_ID
# By URL >
./sc-cli download $TRACK_URL
```

## Features

- ID3 Tags
- Support multiple system-architectures

## Planned Features

- Support multiple tracks download (e.g. by playlist URL)
- Search tracks (shows results, e.g. ID, URL, etc.)


## Contribution

PRs Welcome :)
