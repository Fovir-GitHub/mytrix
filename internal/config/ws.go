package config

type WSConfig struct {
	// RetryInterval defines the retry interval (in second) between websocket requests.
	RetryInterval int `env:"WS_RETRY_INTERVAL" envDefault:"3"`
	// RecvBufferSize defines the channel size of websocket client.
	RecvBufferSize int `env:"WS_RECV_BUFFER_SIZE" envDefault:"64"`
}
