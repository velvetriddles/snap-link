package main

import (
	"fmt"

	"github.com/velvetriddles/snap-link/internal/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
}
