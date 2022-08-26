# streambot [![Test](https://github.com/vikpe/streambot/actions/workflows/test.yml/badge.svg)](https://github.com/vikpe/streambot/actions/workflows/test.yml) [![codecov](https://codecov.io/gh/vikpe/streambot/branch/main/graph/badge.svg)](https://codecov.io/gh/vikpe/streambot)

> Automated QuakeWorld client controlled via on Twitch chat.

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

![image](https://user-images.githubusercontent.com/1616817/178297376-f4f79a29-94c6-4dce-bb50-95183ef8dfb6.png)

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
