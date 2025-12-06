package pumping

const (
	// AlertLevel for alert

	AlertLevelWarning = "warn"
	AlertLevelError   = "error"

	// AlertAction for network connection

	AlertActionConnect    = "connect"
	AlertActionDisConnect = "disconnect"
	AlertActionClose      = "close"
)

// AlertHandler Alert消息处理器接口
type AlertHandler interface {
	SendAlert(level, action, message string) // 发送alert消息回调
}
