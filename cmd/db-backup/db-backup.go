package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/viper"
	"github.com/vvjke314/mkc-backend/internal/pkg/config"
)

func main() {
	config.GetConfig()
	// Параметры подключения к базе данных
	dbHost := "localhost"
	dbPort := viper.GetString("DATABASE_PORT")
	dbUser := viper.GetString("DATABASE_USERNAME")
	dbName := viper.GetString("DATABASE_NAME")
	dbPassword := viper.GetString("DATABASE_PASSWORD")

	// Путь к файлу бэкапа
	backupPath := fmt.Sprintf("./%s_backup_%s.sql", dbName, time.Now().Format("20060102_150405"))

	// Формирование команды pg_dump
	cmd := exec.Command("pg_dump",
		"-h", dbHost,
		"-p", dbPort,
		"-U", dbUser,
		"-d", dbName,
		"-F", "c",
		"-b",
		"-v",
		"-f", backupPath,
	)

	// Установка переменной окружения для пароля
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", dbPassword))

	// Выполнение команды
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		fmt.Printf("Output: %s\n", output)
		return
	}

	fmt.Printf("Backup successful! File saved to: %s\n", backupPath)
}