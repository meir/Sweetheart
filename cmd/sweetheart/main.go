package main

import (
	"os"

	"github.com/meir/Sweetheart/internal/app/sweetheart"
)

func init() {
	println("Running Sweetheart Version Hash", os.Getenv("VERSION"))
}

func main() {
	sweetheart.Sweetheart()
}
