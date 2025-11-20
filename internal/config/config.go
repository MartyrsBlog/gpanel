package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Auth     AuthConfig     `yaml:"auth"`
	SSL      SSLConfig      `yaml:"ssl"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

type DatabaseConfig struct {
	Type     string `yaml:"type"`
	SQLite   SQLiteConfig   `yaml:"sqlite"`
	MySQL    MySQLConfig    `yaml:"mysql"`
}

type SQLiteConfig struct {
	Path string `yaml:"path"`
}

type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type AuthConfig struct {
	JWTSecret     string `yaml:"jwt_secret"`
	TokenExpire   int    `yaml:"token_expire"`
	DefaultUser   string `yaml:"default_user"`
	DefaultPasswd string `yaml:"default_passwd"`
}

type SSLConfig struct {
	Email      string `yaml:"email"`
	Staging    bool   `yaml:"staging"`
}

var DB *sql.DB

func Load() (*Config, error) {
	configPath := "internal/config/config.yaml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 创建默认配置文件
		defaultConfig := &Config{
			Server: ServerConfig{
				Port: "8080",
				Host: "0.0.0.0",
			},
			Database: DatabaseConfig{
				Type: "sqlite",
				SQLite: SQLiteConfig{
					Path: "./data/gpanel.db",
				},
				MySQL: MySQLConfig{
					Host:     "localhost",
					Port:     "3306",
					Database: "gpanel",
					Username: "root",
					Password: "",
				},
			},
			Auth: AuthConfig{
				JWTSecret:     "gpanel-secret-key-change-in-production",
				TokenExpire:   86400, // 24 hours
				DefaultUser:   "admin",
				DefaultPasswd: "admin123",
			},
			SSL: SSLConfig{
				Email:   "admin@example.com",
				Staging: false,
			},
		}

		data, err := yaml.Marshal(defaultConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal default config: %w", err)
		}

		if err := os.MkdirAll("internal/config", 0755); err != nil {
			return nil, fmt.Errorf("failed to create config directory: %w", err)
		}

		if err := os.WriteFile(configPath, data, 0644); err != nil {
			return nil, fmt.Errorf("failed to write default config: %w", err)
		}

		return defaultConfig, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func InitDB(cfg *Config) error {
	var dsn string
	var driver string

	switch cfg.Database.Type {
	case "sqlite":
		driver = "sqlite3"
		dsn = cfg.Database.SQLite.Path
		
		// 确保数据库目录存在
		if err := os.MkdirAll("data", 0755); err != nil {
			return fmt.Errorf("failed to create data directory: %w", err)
		}
		
	case "mysql":
		driver = "mysql"
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Database.MySQL.Username,
			cfg.Database.MySQL.Password,
			cfg.Database.MySQL.Host,
			cfg.Database.MySQL.Port,
			cfg.Database.MySQL.Database,
		)
	default:
		return fmt.Errorf("unsupported database type: %s", cfg.Database.Type)
	}

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	DB = db

	// 自动迁移表结构
	if err := migrateDB(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}

func migrateDB() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(50) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			email VARCHAR(100),
			role VARCHAR(20) DEFAULT 'user',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS sites (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			domain VARCHAR(255) NOT NULL,
			path VARCHAR(500) NOT NULL,
			php_version VARCHAR(10) DEFAULT '74',
			ssl_enabled BOOLEAN DEFAULT FALSE,
			ssl_cert_path VARCHAR(500),
			ssl_key_path VARCHAR(500),
			nginx_config TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS databases (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(100) NOT NULL,
			type VARCHAR(20) NOT NULL,
			username VARCHAR(50) NOT NULL,
			password VARCHAR(255) NOT NULL,
			host VARCHAR(100) DEFAULT 'localhost',
			port VARCHAR(10) DEFAULT '3306',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			return fmt.Errorf("failed to execute migration query: %w", err)
		}
	}

	// 创建默认管理员用户
	if err := createDefaultUser(); err != nil {
		return fmt.Errorf("failed to create default user: %w", err)
	}

	return nil
}

func createDefaultUser() error {
	cfg, _ := Load()
	
	// 检查是否已存在管理员用户
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'admin'").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		// 创建默认管理员用户 (密码使用简单的哈希，生产环境应使用 bcrypt)
		_, err = DB.Exec(
			"INSERT INTO users (username, password, email, role) VALUES (?, ?, ?, ?)",
			cfg.Auth.DefaultUser,
			cfg.Auth.DefaultPasswd, // 生产环境应使用哈希
			"admin@example.com",
			"admin",
		)
		if err != nil {
			return err
		}
	}

	return nil
}