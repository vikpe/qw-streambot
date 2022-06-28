package topics

// Examples — for events raised before/after state change:
//
// FileDownloading / FileDownloaded
// TemperatureChanging / TemperatureChanged
// MailArriving / MailArrived

const (
	ClientConnect = "client.connect"
	ClientCommand = "client.command"

	ClientStarted      = "client.started"
	ClientStopped      = "client.stopped"
	ClientConnected    = "client.connected"
	ClientDisconnected = "client.disconnected"

	ServerMapChanged    = "server.map_changed"
	ServerStatusChanged = "server.status_changed"
	ServerScoreChanged  = "server.score_changed"
	ServerTitleChanged  = "server.title_changed"

	SystemHealthCheck = "system.health_check"
	SystemUpdate      = "system.update"
	StopClient        = "client.stop"
	StopBrain         = "brain.stop"
	StopChatbot       = "chatbot.stop"
)
