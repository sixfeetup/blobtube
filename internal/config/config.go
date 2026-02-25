package config

import (
	"os"
	"strconv"
)

type Config struct {
	HTTPSAddr   string
	HTTPAddr    string
	TLSCertFile string
	TLSKeyFile  string
	StaticDir   string

	LogLevel string
	DevMode  bool

	YtDLPPath string
}

func FromEnv() Config {
	httpsPort := envInt("PORT", 8443)
	httpPort := envInt("HTTP_PORT", 8080)

	return Config{
		HTTPSAddr:   ":" + strconv.Itoa(httpsPort),
		HTTPAddr:    ":" + strconv.Itoa(httpPort),
		TLSCertFile: envString("TLS_CERT_FILE", "./certs/server.crt"),
		TLSKeyFile:  envString("TLS_KEY_FILE", "./certs/server.key"),
		StaticDir:   envString("STATIC_DIR", "./web"),
		LogLevel:    envString("LOG_LEVEL", "info"),
		DevMode:     envBool("DEV_MODE", false),
		YtDLPPath:   envString("YTDLP_PATH", "yt-dlp"),
	}
}

func envString(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func envBool(key string, def bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}

func envInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}
