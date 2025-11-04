package main

const (
	envHostKey = "CHAT_GRPC_HOST"
	envPortKey = "CHAT_GRPC_PORT"

	defaultHost = "127.0.0.1"
	defaultPort = "50051"
	defaultRoom = "general"

	promptDisplayName = "Qual nome você quer usar? "
	promptRoom        = "Sala (deixe em branco para general): "
	promptInput       = "> "

	messageConnected        = "✅ Conectado à sala %q como %s"
	messagePromptCommands   = "Digite mensagens e pressione Enter. Use !quit para sair."
	messageServerClosed     = "⚠️ Conexão encerrada pelo servidor."
	messageReceiveError     = "⚠️ Erro ao receber mensagens: %v\n"
	messageSendError        = "⚠️ Erro ao enviar mensagem: %v\n"
	messageLeaveError       = "⚠️ Erro ao sair da sala: %v\n"
	messageInvalidJoinAck   = "⚠️ Resposta inesperada do servidor, encerrando..."
	messageLeaving          = "Saindo da sala..."
	messageDisconnected     = "👋 Até logo!"
	messageNoticeUserJoined = "👤 %s entrou na sala"
	messageNoticeUserLeft   = "👤 %s saiu da sala"
	messageNoticeGeneric    = "💬 %s"
	messageIncomingChat     = "[%s] %s: %s"
	messageSystemError      = "❗ %s"
	messageUnknownEvent     = "❗ Evento desconhecido recebido"

	timeDisplayFormat = "15:04:05"
	commandQuit       = "!quit"
)
