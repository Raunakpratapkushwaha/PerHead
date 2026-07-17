package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/api/routers"
	"github.com/Raunakpratapkushwaha/Batwara/backend/internal/config"
	"github.com/Raunakpratapkushwaha/Batwara/backend/pkg/token"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("[DEBUG] Step 1: Loading configuration...")
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load configuration: %v", err)
	}
	fmt.Println("[DEBUG] Config loaded successfully. Target Port:", cfg.Port)

	fmt.Println("[DEBUG] Step 2: Dialing Postgres Database...")
	db, err := sqlx.Connect("postgres", cfg.DBSource)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}
	defer db.Close()
	fmt.Println("[DEBUG] Database connection established.")

	fmt.Println("[DEBUG] Step 3: Initializing Token Maker...")
	tokenMaker, err := token.NewTokenMaker(cfg.JWTAccessSecret, cfg.JWTRefreshSecret)
	if err != nil {
		log.Fatalf("cannot initialize token maker: %v", err)
	}

	fmt.Println("[DEBUG] Step 4: Setting up Gin Router...")
	router := routers.SetupRouter(db, tokenMaker, cfg)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	fmt.Printf("🔥 Starting Paisa Ka Batwara API Server on port %s\n", cfg.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failure: %v", err)
	}
}
