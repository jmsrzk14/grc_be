package conf

// Bootstrap is the root configuration struct.
type Bootstrap struct {
	Server *Server `yaml:"server"`
	Data   *Data   `yaml:"data"`
	Log    *Log    `yaml:"log"`
}

// Server holds HTTP server configuration.
type Server struct {
	HTTP *HTTPServer `yaml:"http"`
}

// HTTPServer defines bind address and timeout for the HTTP server.
type HTTPServer struct {
	Addr    string `yaml:"addr"`
	Timeout string `yaml:"timeout"`
}

// Data holds database configuration.
type Data struct {
	Database *Database `yaml:"database"`
}

// Database defines the DB driver and DSN source.
type Database struct {
	Driver string `yaml:"driver"`
	Source string `yaml:"source"`
}

// Log holds logging configuration.
type Log struct {
	Level string `yaml:"level"`
}
