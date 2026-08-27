package main

import (
	"context"
	"datalab_api/internal/auth"
	"datalab_api/internal/config"
	"datalab_api/internal/database"
	"flag"
	"fmt"
	"log"
)

func main() {
	name := flag.String("name", "", "label for this API key (required)")
	rpm := flag.Int("rpm", 60, "requests allowed per minute")
	quota := flag.Int("quota", 1000, "requests allowed per calendar month (ignored if -internal)")
	internal := flag.Bool("internal", false, "mark as first-party/internal key — exempt from monthly quota")
	flag.Parse()

	if *name == "" {
		log.Fatal("-name is required")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	db := database.NewPostgresDB(ctx, cfg.DatabaseURL)
	defer db.Close()

	repo := auth.NewRepository(db)

	plaintext, err := auth.GenerateKey()
	if err != nil {
		log.Fatalf("failed to generate key: %v", err)
	}
	if _, err := repo.Create(ctx, *name, auth.HashKey(plaintext), *rpm, *quota, *internal); err != nil {
		log.Fatalf("failed to save key: %v", err)
	}

	fmt.Println("API key created — save this now, it will not be shown again:")
	fmt.Println()
	fmt.Println("  " + plaintext)
	fmt.Println()
	if *internal {
		fmt.Printf("Name: %s | Rate: %d req/min | Internal (no monthly quota)\n", *name, *rpm)
	} else {
		fmt.Printf("Name: %s | Rate: %d req/min | Quota: %d req/month\n", *name, *rpm, *quota)
	}
}
