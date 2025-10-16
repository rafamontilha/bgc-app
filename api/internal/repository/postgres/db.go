package postgres

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "time"

    _ "github.com/lib/pq"

    "bgc-app/internal/config"
)

func MustConnect(cfg *config.AppConfig) *sql.DB {
    // Priorizar DATABASE_URL se existir
    dsn := os.Getenv("DATABASE_URL")
    
    // Se não tiver DATABASE_URL, construir do jeito antigo
    if dsn == "" {
        dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
            cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
    }

    var conn *sql.DB
    var err error
    for i := 0; i < 30; i++ {
        conn, err = sql.Open("postgres", dsn)
        if err == nil {
            if pingErr := conn.Ping(); pingErr == nil {
                log.Printf("Connected to Postgres successfully")
                return conn
            } else {
                err = pingErr
            }
        }
        log.Printf("Waiting for Postgres... (%d/30): %v", i+1, err)
        time.Sleep(2 * time.Second)
    }
    log.Fatalf("Failed to connect to DB: %v", err)
    return nil
}
