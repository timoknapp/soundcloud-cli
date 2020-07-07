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
   search, ls    Search for a tracks
   help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)

```

### Run

Search:
```
./sc-cli ls --limit 5 andhim
     ID     | ARTIST |             TITLE              |                         ARTWORK-URL                           
------------+--------+--------------------------------+---------------------------------------------------------------
  177318577 |        | Elderbrook x Andhim - How Many | https://i1.sndcdn.com/artworks-000133571398-7ytejm-large.jpg  
            |        | Times                          |                                                               
  29552513  |        | Theophilus London - Wine &     | https://i1.sndcdn.com/artworks-000034576883-5s9hru-large.jpg  
            |        | Chocolates (andhim rmx)        |                                                               
  90830301  |        | Hausch                         | https://i1.sndcdn.com/artworks-y7iXz1VWh3ic-0-large.jpg       
  112750743 |        | Boy Boy Boy                    | https://i1.sndcdn.com/artworks-FJ0GdILknqpA-0-large.jpg       
  99537297  |        | andhim live at Fusion Festival | https://i1.sndcdn.com/artworks-000052089664-ruhdt8-large.jpg
```

Download:
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

- Download tracks by URL or ID
- Search tracks
- Show metadata for a track
- ID3 Tags
- Support multiple system-architectures

## Planned Features

- Support multiple tracks download (e.g. by playlist URL)


## Contribution

PRs Welcome :)
