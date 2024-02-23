package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"leetgode/internal/db/postgres"
	"leetgode/internal/leetgode"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type Config struct {
	Database struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Name     string `json:"name"`
	} `json:"database"`
}

func loadConfig() (Config, error) {
	var config Config
	file, err := os.Open("config.json")
	if err != nil {
		return config, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, fmt.Errorf("failed to decode config: %w", err)
	}

	return config, nil
}

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to load config: %w", err))
	}

	url := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=require",
		config.Database.Username,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name)
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect to PostgreSQL: %w", err))
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "internal/db/migrations/"); err != nil {
		panic(err)
	}

	userRepository := postgres.NewUserRepository(db)
	notificationRepository := postgres.NewNotificationRepository(db)
	leetcodeRepository := postgres.NewLeetcodeRepository(db)

	err = godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	token := os.Getenv("BOT_TOKEN")

	lg := leetgode.NewLeetgode(token, userRepository, notificationRepository, leetcodeRepository)
	lg.Run()
}
