package main

const (
	logMsgConfigLoaded      = "configuration loaded"
	logMsgInitializeDeps    = "initialize dependencies"
	logMsgServerFailed      = "server failed"
	logMsgInitObservability = "initialize observability"
	logFieldApp             = "app"
	logFieldEnv             = "env"
	logFieldAddr            = "addr"
	logFieldError           = "error"
	stderrLoadConfigPref    = "load config: "
	stderrInitLoggerPref    = "init logger: "

	// MsgConfigLoaded is the message for when the configuration is loaded.
	MsgConfigLoaded = "configuration loaded"
)

const (
	ErrLoadConfig = "failed to load configuration"
)
