// Package sqlite — SQLite-Implementierung der Repository-Interfaces.
package sqlite

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
	"github.com/coolvegan/hn-mentalhealth/internal/repository"
)

// Compile-time interface checks.
var (
	_ repository.SharedRepository = (*SharedStore)(nil)
	_ repository.HNRepository     = (*HNStore)(nil)
	_ repository.HealthRepository = (*HealthStore)(nil)
)

// DB hält die SQLite-Verbindung und alle drei Domänen-Stores.
type DB struct {
	conn   *sql.DB
	Shared *SharedStore
	HN     *HNStore
	Health *HealthStore
}

// New öffnet eine SQLite-Datenbank und führt die Schema-Migration aus.
func New(path string) (*DB, error) {
	conn, err := sql.Open("sqlite", path+"?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("sqlite open: %w", err)
	}
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("sqlite ping: %w", err)
	}
	if err := migrate(conn); err != nil {
		return nil, fmt.Errorf("sqlite migrate: %w", err)
	}
	return &DB{
		conn:   conn,
		Shared: &SharedStore{conn: conn},
		HN:     &HNStore{conn: conn},
		Health: &HealthStore{conn: conn},
	}, nil
}

// Close schließt die Datenbankverbindung.
func (db *DB) Close() error {
	return db.conn.Close()
}

// migrate führt die Schema-DDL aus allen drei SQL-Dateien aus.
func migrate(conn *sql.DB) error {
	// Die Schemata sind in den .sql-Dateien definiert.
	// Für die SQLite-Implementierung führen wir sie als embedded SQL aus.
	// In der Praxis würden wir embed.FS nutzen; hier der Einfachheit halber inline.
	ddl := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		username TEXT NOT NULL UNIQUE,
		first_seen_at TEXT,
		thread_count INTEGER DEFAULT 0
	);
	CREATE TABLE IF NOT EXISTS confidence_levels (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		ai_hint TEXT,
		ai_note TEXT
	);
	CREATE TABLE IF NOT EXISTS icd11_codes (
		id INTEGER PRIMARY KEY,
		code TEXT NOT NULL UNIQUE,
		label TEXT NOT NULL,
		is_icd11 INTEGER NOT NULL DEFAULT 1,
		note TEXT,
		ai_hint TEXT,
		ai_note TEXT
	);
	CREATE TABLE IF NOT EXISTS ai_knowledge (
		id INTEGER PRIMARY KEY,
		table_name TEXT NOT NULL,
		column_name TEXT,
		value TEXT,
		explanation TEXT,
		ai_note TEXT,
		UNIQUE (table_name, column_name, value)
	);
	CREATE TABLE IF NOT EXISTS patterns (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		description TEXT,
		ai_hint TEXT,
		ai_note TEXT
	);
	CREATE TABLE IF NOT EXISTS emotions (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		ai_hint TEXT,
		ai_note TEXT
	);
	CREATE TABLE IF NOT EXISTS characterizations (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		ai_hint TEXT,
		ai_note TEXT
	);
	CREATE TABLE IF NOT EXISTS interactions (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		ai_hint TEXT,
		ai_note TEXT
	);
	CREATE TABLE IF NOT EXISTS threads (
		id INTEGER PRIMARY KEY,
		story_id INTEGER NOT NULL UNIQUE,
		title TEXT,
		score INTEGER,
		date TEXT,
		url TEXT
	);
	CREATE TABLE IF NOT EXISTS thread_quality (
		thread_id INTEGER PRIMARY KEY REFERENCES threads(id),
		characterization_id INTEGER REFERENCES characterizations(id),
		interaction_id INTEGER REFERENCES interactions(id),
		trauma_potential TEXT,
		healthy_count INTEGER DEFAULT 0,
		notable_count INTEGER DEFAULT 0
	);
	CREATE TABLE IF NOT EXISTS thread_emotions (
		thread_id INTEGER REFERENCES threads(id),
		emotion_id INTEGER REFERENCES emotions(id),
		count INTEGER DEFAULT 0,
		PRIMARY KEY (thread_id, emotion_id)
	);
	CREATE TABLE IF NOT EXISTS hn_findings (
		id INTEGER PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id),
		thread_id INTEGER NOT NULL REFERENCES threads(id),
		pattern_id INTEGER REFERENCES patterns(id),
		confidence_id INTEGER REFERENCES confidence_levels(id),
		icd11_id INTEGER REFERENCES icd11_codes(id),
		evidence TEXT,
		comment_score INTEGER DEFAULT 0,
		ai_hint TEXT,
		ai_note TEXT
	);
	CREATE TABLE IF NOT EXISTS gene_snippets (
		id INTEGER PRIMARY KEY,
		patient_id INTEGER NOT NULL REFERENCES users(id),
		gene TEXT NOT NULL,
		variant TEXT,
		genotype TEXT,
		ai_hint TEXT,
		ai_note TEXT
	);
	CREATE TABLE IF NOT EXISTS gene_patterns (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		description TEXT,
		ai_hint TEXT,
		ai_note TEXT
	);
	CREATE TABLE IF NOT EXISTS pattern_snippets (
		pattern_id INTEGER REFERENCES gene_patterns(id),
		snippet_id INTEGER REFERENCES gene_snippets(id),
		PRIMARY KEY (pattern_id, snippet_id)
	);
	CREATE TABLE IF NOT EXISTS patient_patterns (
		patient_id INTEGER REFERENCES users(id),
		pattern_id INTEGER REFERENCES gene_patterns(id),
		ai_note TEXT,
		PRIMARY KEY (patient_id, pattern_id)
	);
	CREATE TABLE IF NOT EXISTS biomarkers (
		id INTEGER PRIMARY KEY,
		patient_id INTEGER NOT NULL REFERENCES users(id),
		name TEXT NOT NULL,
		value REAL,
		unit TEXT,
		reference_low REAL,
		reference_high REAL,
		measured_at TEXT,
		ai_hint TEXT,
		ai_note TEXT
	);
	CREATE TABLE IF NOT EXISTS health_conditions (
		id INTEGER PRIMARY KEY,
		patient_id INTEGER NOT NULL REFERENCES users(id),
		condition TEXT NOT NULL,
		diagnosed_at TEXT,
		ai_hint TEXT,
		ai_note TEXT
	);
	CREATE TABLE IF NOT EXISTS life_context (
		id INTEGER PRIMARY KEY,
		patient_id INTEGER NOT NULL REFERENCES users(id),
		context TEXT NOT NULL,
		noted_at TEXT,
		ai_hint TEXT,
		ai_note TEXT
	);
	CREATE TABLE IF NOT EXISTS health_findings (
		id INTEGER PRIMARY KEY,
		patient_id INTEGER NOT NULL REFERENCES users(id),
		pattern_id INTEGER REFERENCES gene_patterns(id),
		confidence_id INTEGER REFERENCES confidence_levels(id),
		icd11_id INTEGER REFERENCES icd11_codes(id),
		etiology TEXT,
		evidence_chain TEXT,
		evidence TEXT,
		ai_hint TEXT,
		ai_note TEXT
	);
	CREATE INDEX IF NOT EXISTS idx_ai_knowledge_table ON ai_knowledge(table_name);
	CREATE INDEX IF NOT EXISTS idx_ai_knowledge_value ON ai_knowledge(value);
	CREATE INDEX IF NOT EXISTS idx_hn_findings_user ON hn_findings(user_id);
	CREATE INDEX IF NOT EXISTS idx_hn_findings_thread ON hn_findings(thread_id);
	CREATE INDEX IF NOT EXISTS idx_hn_findings_pattern ON hn_findings(pattern_id);
	CREATE INDEX IF NOT EXISTS idx_thread_emotions_tid ON thread_emotions(thread_id);
	CREATE INDEX IF NOT EXISTS idx_health_findings_patient ON health_findings(patient_id);
	CREATE INDEX IF NOT EXISTS idx_health_findings_pattern ON health_findings(pattern_id);
	CREATE INDEX IF NOT EXISTS idx_gene_snippets_patient ON gene_snippets(patient_id);
	CREATE INDEX IF NOT EXISTS idx_biomarkers_patient ON biomarkers(patient_id);
	CREATE INDEX IF NOT EXISTS idx_health_conditions_patient ON health_conditions(patient_id);
	CREATE INDEX IF NOT EXISTS idx_life_context_patient ON life_context(patient_id);
	`
	_, err := conn.Exec(ddl)
	return err
}

// nullString und nullInt64 sind Helfer für nullable SQL-Felder.
func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

func nullInt64Ptr(p *int64) sql.NullInt64 {
	if p == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *p, Valid: true}
}

func int64PtrFromNull(n sql.NullInt64) *int64 {
	if n.Valid {
		return &n.Int64
	}
	return nil
}

func stringFromNull(n sql.NullString) string {
	if n.Valid {
		return n.String
	}
	return ""
}
