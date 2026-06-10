package main

import (
	"fmt"

	"github.com/brunosilv96/bs-aesthetics-api/pkg"
)

func main() {
	cfg := pkg.Load()

	fmt.Print(cfg.DatabaseURL)
}
