// Package sqlite — Integrationstest für die SQLite-Implementierung.
package sqlite

import (
	"testing"

	"github.com/coolvegan/hn-mentalhealth"
)

func TestSQLiteIntegration(t *testing.T) {
	db, err := New(":memory:")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer db.Close()

	// --- Shared: User ---
	u := &hnmentalhealth.User{Username: "testuser", FirstSeenAt: "2026-08-12", ThreadCount: 1}
	if err := db.Shared.CreateUser(u); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	if u.ID == 0 {
		t.Fatal("CreateUser did not set ID")
	}
	got, err := db.Shared.GetUserByID(u.ID)
	if err != nil {
		t.Fatalf("GetUserByID: %v", err)
	}
	if got.Username != "testuser" {
		t.Fatalf("expected testuser, got %s", got.Username)
	}

	// --- Shared: ConfidenceLevel ---
	cl := &hnmentalhealth.ConfidenceLevel{Name: "secure", AIHint: "hohe Konfidenz"}
	if err := db.Shared.CreateConfidenceLevel(cl); err != nil {
		t.Fatalf("CreateConfidenceLevel: %v", err)
	}

	// --- Shared: ICD11Code ---
	icd := &hnmentalhealth.ICD11Code{Code: "6B70", Label: "Depressive Episode", IsICD11: true}
	if err := db.Shared.CreateICD11Code(icd); err != nil {
		t.Fatalf("CreateICD11Code: %v", err)
	}

	// --- HN: Pattern ---
	p := &hnmentalhealth.Pattern{Name: "workism", Description: "Überidentifikation mit Arbeit"}
	if err := db.HN.CreatePattern(p); err != nil {
		t.Fatalf("CreatePattern: %v", err)
	}

	// --- HN: Thread ---
	th := &hnmentalhealth.Thread{StoryID: 12345, Title: "Test Story", Score: 100}
	if err := db.HN.CreateThread(th); err != nil {
		t.Fatalf("CreateThread: %v", err)
	}

	// --- HN: HNFinding ---
	f := &hnmentalhealth.HNFinding{
		UserID:       u.ID,
		ThreadID:     th.ID,
		PatternID:    &p.ID,
		ConfidenceID: &cl.ID,
		ICD11ID:      &icd.ID,
		Evidence:     "Ich arbeite 80h/Woche und bin stolz drauf",
		CommentScore: 42,
	}
	if err := db.HN.CreateHNFinding(f); err != nil {
		t.Fatalf("CreateHNFinding: %v", err)
	}

	// --- HN: List by user ---
	findings, err := db.HN.ListHNFindingsByUser(u.ID)
	if err != nil {
		t.Fatalf("ListHNFindingsByUser: %v", err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}

	// --- Health: GeneSnippet ---
	gs := &hnmentalhealth.GeneSnippet{
		PatientID: u.ID,
		Gene:      "MTHFR",
		Variant:   "C677T",
		Genotype:  "C/T",
		AIHint:    "Studien zeigen reduzierte Folsäure-Verwertung",
	}
	if err := db.Health.CreateGeneSnippet(gs); err != nil {
		t.Fatalf("CreateGeneSnippet: %v", err)
	}

	// --- Health: GenePattern ---
	gp := &hnmentalhealth.GenePattern{
		Name:        "Methylierungs-Störung",
		Description: "Kombination aus MTHFR + MTRR",
	}
	if err := db.Health.CreateGenePattern(gp); err != nil {
		t.Fatalf("CreateGenePattern: %v", err)
	}

	// --- Health: PatternSnippet ---
	if err := db.Health.AddSnippetToPattern(gp.ID, gs.ID); err != nil {
		t.Fatalf("AddSnippetToPattern: %v", err)
	}

	// --- Health: PatientPattern ---
	if err := db.Health.AddPatternToPatient(u.ID, gp.ID); err != nil {
		t.Fatalf("AddPatternToPatient: %v", err)
	}

	// --- Health: Biomarker ---
	b := &hnmentalhealth.Biomarker{
		PatientID:     u.ID,
		Name:          "Folsäure",
		Value:         3.1,
		Unit:          "ng/ml",
		ReferenceLow:  4.0,
		ReferenceHigh: 20.0,
		MeasuredAt:    "2026-08-01",
	}
	if err := db.Health.CreateBiomarker(b); err != nil {
		t.Fatalf("CreateBiomarker: %v", err)
	}

	// --- Health: ListBiomarkersOutOfRange ---
	outOfRange, err := db.Health.ListBiomarkersOutOfRange(u.ID)
	if err != nil {
		t.Fatalf("ListBiomarkersOutOfRange: %v", err)
	}
	if len(outOfRange) != 1 {
		t.Fatalf("expected 1 out-of-range biomarker, got %d", len(outOfRange))
	}

	// --- Health: HealthFinding ---
	hf := &hnmentalhealth.HealthFinding{
		PatientID:     u.ID,
		PatternID:     &gp.ID,
		ConfidenceID:  &cl.ID,
		ICD11ID:       &icd.ID,
		Etiology:      "metabolic",
		EvidenceChain: `["MTHFR C677T", "Folsäure 3.1 ng/ml", "Studie X"]`,
		Evidence:      "Metabolische Depression via Folsäure-Mangel",
	}
	if err := db.Health.CreateHealthFinding(hf); err != nil {
		t.Fatalf("CreateHealthFinding: %v", err)
	}

	// --- Health: List by patient ---
	hfs, err := db.Health.ListHealthFindingsByPatient(u.ID)
	if err != nil {
		t.Fatalf("ListHealthFindingsByPatient: %v", err)
	}
	if len(hfs) != 1 {
		t.Fatalf("expected 1 health finding, got %d", len(hfs))
	}

	// --- AIKnowledge ---
	ak := &hnmentalhealth.AIKnowledge{
		TableName:   "patterns",
		ColumnName:  "name",
		Value:       "workism",
		Explanation: "Überidentifikation mit Arbeit als Coping-Mechanismus",
	}
	if err := db.Shared.CreateAIKnowledge(ak); err != nil {
		t.Fatalf("CreateAIKnowledge: %v", err)
	}

	// --- AIKnowledge: Search ---
	results, err := db.Shared.SearchAIKnowledge("arbeit", 10)
	if err != nil {
		t.Fatalf("SearchAIKnowledge: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("SearchAIKnowledge returned no results")
	}

	t.Logf("All integration tests passed. User ID=%d, Finding ID=%d, HealthFinding ID=%d", u.ID, f.ID, hf.ID)
}
