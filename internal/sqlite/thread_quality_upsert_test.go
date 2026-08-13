package sqlite

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/coolvegan/hn-mentalhealth"
)

// Regressionstest für den thread_quality-Upsert-Fix:
// UpdateThreadQuality auf eine nicht existierende Zeile MUSS sql.ErrNoRows
// liefern (statt nil), damit der MCP-Handler auf CreateThreadQuality
// zurückfällt. Danach muss die Zeile existieren und Update nil liefern.
func TestThreadQualityUpsertFallback(t *testing.T) {
	dir := t.TempDir()
	db, err := New(filepath.Join(dir, "verify.db"))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer db.Close()

	th := &hnmentalhealth.Thread{StoryID: 999001, Title: "verify", Score: 1}
	if err := db.HN.CreateThread(th); err != nil {
		t.Fatalf("CreateThread: %v", err)
	}

	tq := &hnmentalhealth.ThreadQuality{
		ThreadID:        th.ID,
		TraumaPotential: "low",
		HealthyCount:    2,
		NotableCount:    1,
	}

	// 1) Update auf nicht-existente Zeile MUSS ErrNoRows liefern.
	if err := db.HN.UpdateThreadQuality(tq); err != sql.ErrNoRows {
		t.Fatalf("Update auf fehlende Zeile: erwartet sql.ErrNoRows, bekam %v", err)
	}

	// 2) Create legt die Zeile an.
	if err := db.HN.CreateThreadQuality(tq); err != nil {
		t.Fatalf("CreateThreadQuality: %v", err)
	}

	// 3) Update auf existierende Zeile MUSS nil liefern (kein ErrNoRows).
	if err := db.HN.UpdateThreadQuality(tq); err != nil {
		t.Fatalf("Update auf existierende Zeile: erwartet nil, bekam %v", err)
	}

	// 4) Zeile ist persistiert.
	got, err := db.HN.GetThreadQualityByThreadID(th.ID)
	if err != nil {
		t.Fatalf("GetThreadQualityByThreadID: %v", err)
	}
	if got == nil || got.TraumaPotential != "low" || got.HealthyCount != 2 {
		t.Fatalf("persistierte Qualität falsch: %+v", got)
	}
}
