package symbols

const (
	AppName    = "phoo"
	AppVersion = "v2.0.0"

	FlagNameHTTPListenAddr    = "listen"
	FlagNameDocumentRoot      = "root"
	FlagNamePHPFPMBinary      = "php-fpm"
	FlagNamePHPINI            = "php-ini"
	FlagNamePHPUser           = "php-user"
	FlagNamePHPGroup          = "php-group"
	FlagNameWorkersCount      = "workers"
	FlagNameWorkerMaxRequests = "worker-max-requests"
	FlagNameRequestTimeout    = "timeout"
	FlagNameDefaultScript     = "router"
	FlagNameEnableLogs        = "logs"
	FlagNameEnvFilename       = "env-file"

	EnvKeyDocumentRoot      = "PHOO_DOCUMENT_ROOT"
	EnvKeyListen            = "PHOO_LISTEN"
	EnvKeyFPMBin            = "PHOO_FPM_BIN"
	EnvKeyWorkersCount      = "PHOO_WORKERS_COUNT"
	EnvKeyWorkerMaxRequests = "PHOO_WORKER_MAX_REQUESTS"
	EnvKeyRequestTimeout    = "PHOO_REQUEST_TIMEOUT"
	EnvKeyRouter            = "PHOO_ROUTER"
	EnvKeyEnableLogs        = "PHOO_LOGS"
	EnvKeyPHPINI            = "PHOO_PHP_INI"
	EnvKeyPHPUser           = "PHOO_PHP_USER"
	EnvKeyPHPGroup          = "PHOO_PHP_GROUP"
)
