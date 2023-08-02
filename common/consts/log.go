package consts

// default value for logger (pkg/log)
const (
	LogEnableFileWriter = false                          // enable file writer
	LogDisableColor     = false                          // disable colorful output
	LogMaxSize          = 10                             // megabytes
	LogMaxBackups       = 1                              // copies
	LogMaxAge           = 1                              // days
	LogCompress         = true                           // enable by default
	LogFile             = "/var/log/juicity/juicity.log" // path
)
