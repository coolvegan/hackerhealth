// Command hn-mentalhealth — CLI für die hn-mentalhealth Datenbank.
//
// Usage:
//
//	hn-mentalhealth init [--db <path>]   Erstellt DB + Seed-Daten
//	hn-mentalhealth seed [--db <path>]   Nur Seed-Daten (DB muss existieren)
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/coolvegan/hn-mentalhealth/internal/sqlite"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <command> [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  init   Create database and seed initial data\n")
		fmt.Fprintf(os.Stderr, "  seed   Seed existing database with initial data\n")
		os.Exit(1)
	}

	cmd := os.Args[1]
	fs := flag.NewFlagSet(cmd, flag.ExitOnError)
	dbPath := fs.String("db", "hn-mentalhealth.db", "Path to SQLite database")
	fs.Parse(os.Args[2:])

	switch cmd {
	case "init":
		if err := initDB(*dbPath); err != nil {
			fmt.Fprintf(os.Stderr, "init failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Database initialized: %s\n", *dbPath)
	case "seed":
		if err := seedDB(*dbPath); err != nil {
			fmt.Fprintf(os.Stderr, "seed failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Database seeded: %s\n", *dbPath)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		os.Exit(1)
	}
}

func initDB(path string) error {
	db, err := sqlite.New(path)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer db.Close()

	if err := db.Seed(); err != nil {
		return fmt.Errorf("seed: %w", err)
	}
	return nil
}

func seedDB(path string) error {
	db, err := sqlite.New(path)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer db.Close()

	if err := db.Seed(); err != nil {
		return fmt.Errorf("seed: %w", err)
	}
	return nil
}
