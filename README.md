# streambot [![Test](https://github.com/vikpe/streambot/actions/workflows/test.yml/badge.svg)](https://github.com/vikpe/streambot/actions/workflows/test.yml) [![codecov](https://codecov.io/gh/vikpe/streambot/branch/main/graph/badge.svg)](https://codecov.io/gh/vikpe/streambot)

> Automated QuakeWorld streaming on Twitch.

## Overview

### Components

![image](https://user-images.githubusercontent.com/1616817/178285267-eade607d-8660-4b4d-9522-ab3772dde229.png)

* **Message Proxy**: Central point for communication.
* **Quake Manager**: Interaction with ezQuake
    * **Server monitor**: Server events (map change etc)
    * **Process monitor**: ezQuake events (started, stopped)

* **Twitch Manager**: Interaction with Twitch channel (e.g. set title).
* **Twitch Bot**: Interaction with chat users.
* **Discord Bot**: Interaction with Discord.

### Stack

* Written in Go (Golang)
* Communication/messages are sent using [ZeroMQ](https://zeromq.org/) sockets. One proxy, multiple
  subscribrs/publishers.

## Requirements

* **[ezQuake](https://github.com/ezQuake/ezquake-source/releases)**
* **[Twitch access tokens](https://twitchtokengenerator.com/)**
* **ZeroMQ**: `apt-get install libzmq3-dev`

## Test

```shell
go test ./... --cover
```

## Build

```shell
./scripts/build.sh
```
