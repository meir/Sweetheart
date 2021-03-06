package main

import (
	"os"

	"github.com/meir/Sweetheart/internal/app/sweetheart"
	"github.com/meir/Sweetheart/internal/pkg/logging"
)

func init() {
	logging.Println("Running Sweetheart Version Hash", os.Getenv("VERSION"))
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			logging.Fatal("Sweetheart had a hiccup, recovered to get error", r)
		}
	}()

	sweetheart.Sweetheart()
}
