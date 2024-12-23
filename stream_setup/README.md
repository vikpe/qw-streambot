# Stream Setup

* **OS**: Ubuntu `22.04.1`
* [ezQuake](https://ezquake.com/) `3.6.3-dev`
* [Open Broadcaster Sofware (OBS)](https://obsproject.com/)
* [PulseAudio](https://www.freedesktop.org/wiki/Software/PulseAudio/)
* [Advanced Scene Switcher](https://obsproject.com/forum/resources/advanced-scene-switcher.395/)

![image](https://github.com/vikpe/qw-streambot/assets/1616817/bb44ea3b-ff0a-4867-9baf-b2c735ac8126)

## ezQuake setup (config, skins, textures etc)

See [qw-streambot-ezquake](https://github.com/vikpe/qw-streambot-ezquake).

## Scripts (Advanced Scene Switcher)

* Start stream whenever stopped
* Stop stream every 24 hours

Import [advanced_scene_switcher.txt](./advanced_scene_switcher.txt)

## Filters

* Color Correction: `gamma 0.4`

## Virtual audio device

Create a dedicated virtual audio device using [PulseAudio](https://wiki.archlinux.org/title/PulseAudio/Examples) for
audio to include in the stream.

```sh
# install PulseAdio
sudo apt install pavucontrol

# create device "Stream-Sink"
pacmd load-module module-null-sink sink_name=streambot
pactl load-module module-loopback source=streambot.monitor sink=@DEFAULT_SINK@
```

Assign ezQuake to the Stream Sink using PulseAudio.
![image](https://github.com/vikpe/qw-streambot/assets/1616817/3df7a49c-91fd-43f6-bfc5-947be6656e18)
