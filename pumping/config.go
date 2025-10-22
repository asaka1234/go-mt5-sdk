package pumping

import "time"

// Config TCP客户端配置
type Config struct {
	ServerAddr        string        // 服务器地址，格式：host:port
	Timeout           time.Duration // 连接超时时间
	Reconnect         bool          // 是否自动重连
	MaxReconnects     int           // 最大重连次数
	ReconnectInterval time.Duration // 重连间隔
	HeartbeatInterval time.Duration // 心跳间隔
	BufferSize        int           // 读写缓冲区大小
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		ServerAddr:        "localhost:8080",
		Timeout:           10 * time.Second,
		Reconnect:         true,
		MaxReconnects:     5,
		ReconnectInterval: 5 * time.Second,
		HeartbeatInterval: 30 * time.Second,
		BufferSize:        4096,
	}
}
