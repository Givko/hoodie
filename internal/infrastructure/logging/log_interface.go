package logging

type LogInterface interface {
	LogInfo(message string, fields map[string]interface{})
	LogWarn(message string, fields map[string]interface{})
	LogError(message string, fields map[string]interface{})
	LogErr(err error, message string, fields map[string]interface{})
	SetupLogger()
}
