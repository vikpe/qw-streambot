package topics

// Examples — for events raised before/after state change:
//
// FileDownloading / FileDownloaded
// TemperatureChanging / TemperatureChanged
// MailArriving / MailArrived

// events
const (
	ClientStarted      = "client.started"
	ClientStopped      = "client.stopped"
	ClientConnected    = "client.connected"
	ClientDisconnected = "client.disconnected"

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

	ConnectToServer = "streambot.connect_to_server"
	SuggestServer   = "streambot.suggest_server"
	EnableAuto      = "streambot.enable_auto"
	DisableAuto     = "streambot.disable_auto"
	Evaluate        = "streambot.evaluate"
	SystemUpdate    = "streambot.system_update"
)
