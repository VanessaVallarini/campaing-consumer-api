package main

import (
	"campaing-comsumer-service/internal/config"
	"fmt"
)

func main() {
	cfg := config.GetConfig()
	fmt.Println(cfg)
}
