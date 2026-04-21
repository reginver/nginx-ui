package settings

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

var (
	once     sync.Once
	confPath string
)

// ServerSettings holds HTTP server configuration
type ServerSettings struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	RunMode   string `mapstructure:"run_mode"`
	JWTSecret string `mapstructure:"jwt_secret"`
}

// NginxSettings holds nginx binary and config paths
type NginxSettings struct {
	AccessLogPath string `mapstructure:"access_log_path"`
	ErrorLogPath  string `mapstructure:"error_log_path"`
	ConfigDir     string `mapstructure:"config_dir"`
	PIDPath       string `mapstructure:"pid_path"`
	TestConfigCmd string `mapstructure:"test_config_cmd"`
	ReloadCmd     string `mapstructure:"reload_cmd"`
}

// DatabaseSettings holds database connection configuration
type DatabaseSettings struct {
	Name string `mapstructure:"name"`
}

// LogSettings holds application log configuration
type LogSettings struct {
	Level string `mapstructure:"level"`
	Path  string `mapstructure:"path"`
}

// Config is the global application configuration
type Config struct {
	Server   ServerSettings   `mapstructure:"server"`
	Nginx    NginxSettings    `mapstructure:"nginx"`
	Database DatabaseSettings `mapstructure:"database"`
	Log      LogSettings      `mapstructure:"log"`
}

// Conf is the singleton configuration instance
var Conf = &Config{}

// Init loads configuration from the given path, applying defaults where needed.
func Init(path string) {
	once.Do(func() {
		confPath = path
		loadConfig()
	})
}

func loadConfig() {
	v := viper.New()

	// Set defaults
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 9000)
	v.SetDefault("server.run_mode", "release")
	v.SetDefault("nginx.access_log_path", "/var/log/nginx/access.log")
	v.SetDefault("nginx.error_log_path", "/var/log/nginx/error.log")
	v.SetDefault("nginx.config_dir", "/etc/nginx")
	v.SetDefault("nginx.pid_path", "/var/run/nginx.pid")
	v.SetDefault("nginx.test_config_cmd", "nginx -t")
	v.SetDefault("nginx.reload_cmd", "nginx -s reload")
	v.SetDefault("database.name", "database")
	v.SetDefault("log.level", "info")
	v.SetDefault("log.path", "log")

	if confPath != "" {
		v.SetConfigFile(confPath)
	} else {
		v.SetConfigName("app")
		v.SetConfigType("ini")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
	}

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(*os.PathError); ok {
			log.Println("[settings] No config file found, using defaults")
		} else {
			log.Fatalf("[settings] Failed to read config: %v", err)
		}
	}

	if err := v.Unmarshal(Conf); err != nil {
		log.Fatalf("[settings] Failed to unmarshal config: %v", err)
	}

	// Ensure log directory exists
	if Conf.Log.Path != "" {
		if err := os.MkdirAll(filepath.Clean(Conf.Log.Path), 0755); err != nil {
			log.Printf("[settings] Could not create log directory: %v", err)
		}
	}

	log.Printf("[settings] Configuration loaded (run_mode=%s, port=%d)",
		Conf.Server.RunMode, Conf.Server.Port)
}
