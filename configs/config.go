package configs

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"strings"
	"time"
)

var cfg *Config

func Init(path string) error {
	if path == "" {
		path = "configs"
	}

	initViper(path)
	loadConfigs()
	//setTimeZone(cfg.Server.TimeZone)
	//initDb()
	if cfg == nil {
		return fmt.Errorf("config is not initialized")
	}
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("validate config: \n%v\n", err)
	}

	return nil
}

func GetConfig() Config {
	if cfg == nil {
		loadConfigs()
	}
	return *cfg
}

func initViper(path string) error {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("cannot read config file: %s", err)
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return nil
}

func loadConfigs() {
	_ = viper.Unmarshal(&cfg)
}

func SetTimeZone(timeZone string) {
	location, err := time.LoadLocation(timeZone)
	if err != nil {
		log.Fatalf("Failed to load location: %v", err)
	}
	time.Local = location
	log.Println("Timezone set to:", timeZone)
}
func FormatTime(t time.Time) string {
	return t.Format("2006/01/02 15:04:05")
}

//func initDb() {
//	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
//		cfg.DB.Host, cfg.DB.Port, cfg.Secrets.DbUsername, cfg.Secrets.DbPassword, cfg.DB.Database)
//	db, err := sql.Open("postgres", dataSourceName)
//	if err != nil {
//		log.Fatal("Error connecting to the database:", err)
//	}
//	// Verify the connection
//	if err := db.Ping(); err != nil {
//		log.Fatal("Database is unreachable:", err)
//	}
//
//	log.Println("Connected to PostgreSQL database!")
//}
