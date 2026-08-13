// Command hn-mentalhealth-mcp — MCP-SSE-Server + Dashboard für hn-mentalhealth.
//
// Usage:
//
//	hn-mentalhealth-mcp [--db <path>] [--mcp-addr <host:port>] [--dashboard-addr <host:port>]
//
// Environment:
//
//	DB_PATH         — path to SQLite database (default: hn-mentalhealth.db)
//	MCP_ADDR        — MCP SSE listen address (default: localhost:14444)
//	DASHBOARD_ADDR  — Dashboard HTTP listen address (default: localhost:17776)
package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/coolvegan/hn-mentalhealth/internal"
	"github.com/coolvegan/hn-mentalhealth/internal/sqlite"
)

func main() {
	fs := flag.NewFlagSet("hn-mentalhealth-mcp", flag.ExitOnError)
	dbPath := fs.String("db", "hn-mentalhealth.db", "Path to SQLite database")
	mcpAddr := fs.String("mcp-addr", "localhost:14444", "MCP SSE listen address")
	dashAddr := fs.String("dashboard-addr", "localhost:17776", "Dashboard HTTP listen address")
	fs.Parse(os.Args[1:])

	if p := os.Getenv("DB_PATH"); p != "" {
		*dbPath = p
	}
	if a := os.Getenv("MCP_ADDR"); a != "" {
		*mcpAddr = a
	}
	if a := os.Getenv("DASHBOARD_ADDR"); a != "" {
		*dashAddr = a
	}

	db, err := sqlite.New(*dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// MCP-Server (SSE) in eigener Goroutine
	go func() {
		log.Printf("MCP SSE-Server starting on %s", *mcpAddr)
		internal.RunMentalHealthMcp(*mcpAddr, db)
	}()

	// Dashboard (REST + Frontend)
	mux := http.NewServeMux()
	internal.RegisterDashboardRoutes(mux, db)
	log.Printf("Dashboard starting on %s", *dashAddr)
	log.Fatalln(http.ListenAndServe(*dashAddr, mux))
}
