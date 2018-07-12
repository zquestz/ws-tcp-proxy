package server

// Config is the server package config.
type Config struct {
	Port         int
	Cert         string
	Key          string
	TextMode     bool
	Address      string
	TCPTLS       bool
	TCPTLSCert   string
	TCPTLSKey    string
	TCPTLSRootCA string
}
