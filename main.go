package main

import (
	"fmt"
	"log/slog"
	"os"
)

func main() {
	storage, err := NewPostgresStorage()
	if err != nil {
		slog.Error("can't create storage", err)
		os.Exit(1)
	}

	if err := storage.Init(); err != nil {
		slog.Error("can't init storage", err)
		os.Exit(1)
	}

	fmt.Printf("%+v\n", storage)

	server := NewAPIServer(":3000", storage)

	server.Run()
}
