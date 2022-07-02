package topics

// Examples — for events raised before/after state change:
//
// FileDownloading / FileDownloaded
// TemperatureChanging / TemperatureChanged
// MailArriving / MailArrived

// events
const (
	ClientStarted = "client.started"
	ClientStopped = "client.stopped"

	ServerTitleChanged = "server.title_changed"
)

// commands
const (
	ClientCommand           = "client.command"
	ClientCommandLastscores = "client.command.lastscores"
	ClientCommandShowscores = "client.command.showscores"
	StopClient              = "client.stop"

	StreambotConnectToServer = "streambot.connect_to_server"
	StreambotSuggestServer   = "streambot.suggest_server"
	StreambotEnableAuto      = "streambot.enable_auto"
	StreambotDisableAuto     = "streambot.disable_auto"
	StreambotEvaluate        = "streambot.evaluate"
	StreambotSystemUpdate    = "streambot.system_update"
)
