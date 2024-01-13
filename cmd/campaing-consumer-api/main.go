package main

import (
	"campaing-comsumer-service/internal/config"
	"campaing-comsumer-service/internal/db"
	"context"
	"fmt"

	"github.com/lockp111/go-easyzap"
)

func main() {
	ctx := context.Background()

	cfg := config.GetConfig()

	db, err := db.NewDb(ctx, cfg)
	if err != nil {
		easyzap.Panic(err)

	}
	defer db.Close()

	fmt.Println(cfg)
}
