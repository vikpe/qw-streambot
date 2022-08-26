# QW streambot [![Test](https://github.com/vikpe/qw-streambot/actions/workflows/test.yml/badge.svg)](https://github.com/vikpe/qw-streambot/actions/workflows/test.yml) [![codecov](https://codecov.io/gh/vikpe/qw-streambot/branch/main/graph/badge.svg)](https://codecov.io/gh/vikpe/qw-streambot) [![Go Report Card](https://goreportcard.com/badge/github.com/vikpe/qw-streambot)](https://goreportcard.com/report/github.com/vikpe/qw-streambot)

> Automated QuakeWorld client controlled via on Twitch chat.

## Example

Visit [twitch.tv/vikpe](https://www.twitch.tv/vikpe) to see it in action.

## Stack

* Written in Go (Golang)
* [ezQuake](https://github.com/ezQuake/ezquake-source/releases) - QuakeWorld client
* [ZeroMQ](https://zeromq.org/) - Communication/messages (single proxy and multiple subscribers/publishers)
* [serverstat](https://github.com/vikpe/serverstat) - Get server information

## Overview

![image](https://user-images.githubusercontent.com/1616817/186941072-cc99679d-b1d0-41f7-bdba-913bb733e140.png)

* **Message Proxy**: Central point for communication.
* **Quake Manager**: Interaction with ezQuake
    * Server monitor (thread): Server events (map change etc)
    * Process monitor (thread): ezQuake events (started, stopped)

* **Twitch Manager**: Interaction with Twitch channel (e.g. set title).
* **Twitch Bot**: Interaction with Twitch chat.

### Quake Manager - evaluation loop

* Run every 10 seconds
* Join "best server" available. Servers are ranked using a
  custom [scoring algorithm](https://github.com/vikpe/serverstat/blob/main/qserver/mvdsv/qscore/qscore.go).
* Only change server in between matches or if current server has enabled a custom game mode (e.g. `race`).

![image](https://user-images.githubusercontent.com/1616817/186943108-6c87bb9a-72cf-4e20-b288-824a7d292543.png)

## Requirements

* Streaming software, e.g. [Open Broadcaster Sofware (OBS)](https://obsproject.com/)
* Twitch account for the stream
* Twitch account for the chatbot
* [Twitch access tokens](https://twitchtokengenerator.com/) (for chatbot and twitch channel)
* ZeroMQ: `apt-get install libzmq3-dev`
* Create `.env` (see `.env.example`)

## Development

### Directory structure

Uses the [Standard Go Project Layout](https://github.com/golang-standards/project-layout).

```bash
cmd/       # Main applications
internal/  # Private application and library code
scripts/   # Various build, install operations
```

### Build

**Build specific app**

Example: build proxy

```shell
cd cmd/proxy
go build
```

**Build all apps**

```shell
./scripts/build.sh
```

### Run

**Single app**

Example: start the proxy.

```shell
./cmd/proxy/proxy 
```

**App controller scripts**

Runs app forever (restarts on error/sigint with short timeout in between).

```shell
bash scripts/controllers/proxy.sh
bash scripts/controllers/quake_manager.sh
bash scripts/controllers/twitch_manager.sh
bash scripts/controllers/twitchbot.sh
bash scripts/controllers/ezquake.sh
```

### Test

```shell
go test ./... --cover
```

## Production

Build all apps and run all app controller scripts.

```shell
./scripts/build.sh && ./scripts/start.sh
```

## Credits

Thanks to everyone that has provided feedback and improvement suggestions (andeh, bps, circle, hangtime, milton, splash,
wimpeeh) among others.
