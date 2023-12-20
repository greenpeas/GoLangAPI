package pg

import (
	"errors"
	"fmt"
	"log"
	app_interface "seal/internal/app/interface"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"os"
	"time"
)

func Create(path, name string, appLog app_interface.Logger) {
	log.Println("Migration create run...")

	now := time.Now()
	sec := now.Unix()

	fileNameUp := fmt.Sprintf("%s/%d_%s.up.sql", path, sec, name)
	if err := createFile(fileNameUp); err != nil {
		appLog.Fatal(err.Error())
		return
	}

	fileNameDown := fmt.Sprintf("%s/%d_%s.down.sql", path, sec, name)
	if err := createFile(fileNameDown); err != nil {
		appLog.Fatal(err.Error())
		return
	}

	log.Println("File created " + fileNameUp)
	log.Println("File created " + fileNameDown)
}

func createFile(fileName string) error {
	f, err := os.Create(fileName)

	if err != nil {
		return err
	}

	defer f.Close()

	return nil
}

func RunMigrateUp(source, dsn string, appLog app_interface.Logger) {
	log.Println("Migration up run...")

	m, err := migrate.New(source, dsn+"?sslmode=disable")

	if err != nil {
		appLog.Fatal(err.Error())
	}

	if err := m.Up(); errors.Is(err, migrate.ErrNoChange) {
		ver, _, _ := m.Version()
		message := fmt.Sprintf("No new migration. Current version: %d", ver)
		log.Println(message)
		appLog.Info(message)
	} else if err != nil {
		appLog.Fatal(err.Error())
	} else {
		ver, _, _ := m.Version()
		message := fmt.Sprintf("Migrate up complete. Current version %d ", ver)
		log.Println(message)
		appLog.Info(message)
	}
}

func RunMigrateDown(source, dsn string, appLog app_interface.Logger) {
	log.Println("Migration down run...")

	m, err := migrate.New(source, dsn+"?sslmode=disable")

	if err != nil {
		appLog.Fatal(err.Error())
	}

	if err := m.Steps(-1); errors.Is(err, migrate.ErrNoChange) {
		ver, _, _ := m.Version()
		message := fmt.Sprintf("No down migration. Current version: %d", ver)
		log.Println(message)
		appLog.Info(message)
	} else if err != nil {
		appLog.Fatal(err.Error())
	} else {
		ver, _, _ := m.Version()
		message := fmt.Sprintf("Migrate down complete. Current version %d ", ver)
		log.Println(message)
		appLog.Info(message)
	}

}
