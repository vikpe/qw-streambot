package topic

// Examples — for events raised before/after state change:
//
// FileDownloading / FileDownloaded
// TemperatureChanging / TemperatureChanged
// MailArriving / MailArrived

// events
const (
	EzquakeStarted = "ezquake.started"
	EzquakeStopped = "ezquake.stopped"

	ServerTitleChanged    = "server.title_changed"
	ServerMatchtagChanged = "server.matchtag_changed"
)

// commands
const (
	TwitchbotSay = "twitchbot.say"

	EzquakeCommand = "ezquake.command"
	EzquakeScript  = "ezquake.script"
	EzquakeStop    = "ezquake.stop"

	StreambotSuggestServer = "streambot.suggest_server"
	StreambotEnableAuto    = "streambot.enable_auto"
	StreambotDisableAuto   = "streambot.disable_auto"
	StreambotEvaluate      = "streambot.evaluate"
)
