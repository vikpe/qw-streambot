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

	ServerMapChanged    = "server.map_changed"
	ServerStatusChanged = "server.status_changed"
	ServerScoreChanged  = "server.score_changed"
	ServerTitleChanged  = "server.title_changed"
)

// commands
const (
	ClientCommand = "client.command"
	StopClient    = "client.stop"
	StopBrain     = "brain.stop"
	StopChatbot   = "chatbot.stop"

	StreambotConnectToServer = "streambot.connect_to_server"
	StreambotSuggestServer   = "streambot.suggest_server"
	StreambotEnableAuto      = "streambot.enable_auto"
	StreambotDisableAuto     = "streambot.disable_auto"
	StreambotEvaluate        = "streambot.evaluate"
	StreambotSystemUpdate    = "streambot.system_update"
)
