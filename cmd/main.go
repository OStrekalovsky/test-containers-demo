package main

import (
	"fmt"
	"go.uber.org/zap"

	"test-containers/internal"
)

func main() {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Println("Failed to flush logger", err)
		}
	}(logger)
	sugar := logger.Sugar()

	db, err := internal.NewPostgres("postgres://localhost:5432/postgres?user=ost&password=pass")
	if err != nil {
		sugar.Fatal("Failed to connect to db", err)
	}
	defer db.Close()

	greeting, err := db.ReadHelloWorld()
	if err != nil {
		sugar.Fatal("Failed to read string:", err)
	}
	sugar.Info("Read string=", greeting)
}
