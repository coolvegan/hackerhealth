// Package hnmentalhealth — Domain-Types für hn-mentalhealth.
// Spiegelt die drei SQL-Schemata (shared, hn, health) als Go-Strukturen.
package hnmentalhealth

// ============================================================
// SHARED (shared.sql)
// ============================================================

// User ist eine Person — HN-Kommentator oder Patient.
type User struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	FirstSeenAt  string `json:"first_seen_at,omitempty"`
	ThreadCount  int    `json:"thread_count"`
}

// ConfidenceLevel ist eine Konfidenz-Stufe (secure, suspected, observation, healthy).
type ConfidenceLevel struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	AIHint string `json:"ai_hint,omitempty"`
	AINote string `json:"ai_note,omitempty"`
}

// ICD11Code ist eine ICD-11-Referenz oder ein nicht-ICD-Konstrukt.
type ICD11Code struct {
	ID      int64  `json:"id"`
	Code    string `json:"code"`
	Label   string `json:"label"`
	IsICD11 bool   `json:"is_icd11"`
	Note    string `json:"note,omitempty"`
	AIHint  string `json:"ai_hint,omitempty"`
	AINote  string `json:"ai_note,omitempty"`
}

// AIKnowledge ist ein Eintrag im selbst-beschreibenden Daten-Wörterbuch.
type AIKnowledge struct {
	ID          int64  `json:"id"`
	TableName   string `json:"table_name"`
	ColumnName  string `json:"column_name,omitempty"`
	Value       string `json:"value,omitempty"`
	Explanation string `json:"explanation,omitempty"`
	AINote      string `json:"ai_note,omitempty"`
}

// ============================================================
// HN (hn.sql)
// ============================================================

// Pattern ist ein erkanntes psychologisches Muster (workism, scham, flucht, ...).
type Pattern struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	AIHint      string `json:"ai_hint,omitempty"`
	AINote      string `json:"ai_note,omitempty"`
}

// Emotion ist ein Gefühl (joy, anger, sad, fear, shame, ...).
type Emotion struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	AIHint string `json:"ai_hint,omitempty"`
	AINote string `json:"ai_note,omitempty"`
}

// Characterization ist eine Thread-Charakterisierung (values_debate, grief_group, ...).
type Characterization struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	AIHint string `json:"ai_hint,omitempty"`
	AINote string `json:"ai_note,omitempty"`
}

// Interaction ist ein Interaktions-Typ (productive, escalation, polarized, flat).
type Interaction struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	AIHint string `json:"ai_hint,omitempty"`
	AINote string `json:"ai_note,omitempty"`
}

// Thread ist eine HN-Story.
type Thread struct {
	ID      int64  `json:"id"`
	StoryID int64  `json:"story_id"`
	Title   string `json:"title,omitempty"`
	Score   int    `json:"score"`
	Date    string `json:"date,omitempty"`
	URL     string `json:"url,omitempty"`
}

// ThreadQuality ist die qualitative Bewertung eines Threads.
type ThreadQuality struct {
	ThreadID           int64  `json:"thread_id"`
	CharacterizationID *int64 `json:"characterization_id,omitempty"`
	InteractionID      *int64 `json:"interaction_id,omitempty"`
	TraumaPotential    string `json:"trauma_potential,omitempty"`
	HealthyCount       int    `json:"healthy_count"`
	NotableCount       int    `json:"notable_count"`
}

// ThreadEmotion ist ein Emotions-Zähler für einen Thread.
type ThreadEmotion struct {
	ThreadID  int64 `json:"thread_id"`
	EmotionID int64 `json:"emotion_id"`
	Count     int   `json:"count"`
}

// HNFinding ist eine Auffälligkeit eines Users in einem HN-Thread.
type HNFinding struct {
	ID           int64  `json:"id"`
	UserID       int64  `json:"user_id"`
	ThreadID     int64  `json:"thread_id"`
	PatternID    *int64 `json:"pattern_id,omitempty"`
	ConfidenceID *int64 `json:"confidence_id,omitempty"`
	ICD11ID      *int64 `json:"icd11_id,omitempty"`
	Evidence     string `json:"evidence,omitempty"`
	CommentScore int    `json:"comment_score"`
	AIHint       string `json:"ai_hint,omitempty"`
	AINote       string `json:"ai_note,omitempty"`
}

// ============================================================
// HEALTH (health.sql)
// ============================================================

// GeneSnippet ist ein einzelnes Gen eines Patienten.
type GeneSnippet struct {
	ID        int64  `json:"id"`
	PatientID int64  `json:"patient_id"`
	Gene      string `json:"gene"`
	Variant   string `json:"variant,omitempty"`
	Genotype  string `json:"genotype,omitempty"`
	AIHint    string `json:"ai_hint,omitempty"`
	AINote    string `json:"ai_note,omitempty"`
}

// GenePattern ist eine Aggregation mehrerer Gen-Snippets zu einem Muster.
type GenePattern struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	AIHint      string `json:"ai_hint,omitempty"`
	AINote      string `json:"ai_note,omitempty"`
}

// PatternSnippet verknüpft ein Gen-Muster mit einem Gen-Snippet.
type PatternSnippet struct {
	PatternID int64 `json:"pattern_id"`
	SnippetID int64 `json:"snippet_id"`
}

// PatientPattern verknüpft einen Patienten mit einem Gen-Muster.
type PatientPattern struct {
	PatientID int64  `json:"patient_id"`
	PatternID int64  `json:"pattern_id"`
	AINote    string `json:"ai_note,omitempty"`
}

// Biomarker ist ein Blutwert / Laborwert eines Patienten.
type Biomarker struct {
	ID            int64   `json:"id"`
	PatientID     int64   `json:"patient_id"`
	Name          string  `json:"name"`
	Value         float64 `json:"value"`
	Unit          string  `json:"unit,omitempty"`
	ReferenceLow  float64 `json:"reference_low,omitempty"`
	ReferenceHigh float64 `json:"reference_high,omitempty"`
	MeasuredAt    string  `json:"measured_at,omitempty"`
	AIHint        string  `json:"ai_hint,omitempty"`
	AINote        string  `json:"ai_note,omitempty"`
}

// HealthCondition ist eine Diagnose / ein Gesundheitszustand.
type HealthCondition struct {
	ID          int64  `json:"id"`
	PatientID   int64  `json:"patient_id"`
	Condition   string `json:"condition"`
	DiagnosedAt string `json:"diagnosed_at,omitempty"`
	AIHint      string `json:"ai_hint,omitempty"`
	AINote      string `json:"ai_note,omitempty"`
}

// LifeContext ist ein Lebenskontext / eine Sorge eines Patienten.
type LifeContext struct {
	ID        int64  `json:"id"`
	PatientID int64  `json:"patient_id"`
	Context   string `json:"context"`
	NotedAt   string `json:"noted_at,omitempty"`
	AIHint    string `json:"ai_hint,omitempty"`
	AINote    string `json:"ai_note,omitempty"`
}

// HealthFinding ist eine Auffälligkeit aus Gesundheitsdaten.
type HealthFinding struct {
	ID            int64  `json:"id"`
	PatientID     int64  `json:"patient_id"`
	PatternID     *int64 `json:"pattern_id,omitempty"`
	ConfidenceID  *int64 `json:"confidence_id,omitempty"`
	ICD11ID       *int64 `json:"icd11_id,omitempty"`
	Etiology      string `json:"etiology,omitempty"`
	EvidenceChain string `json:"evidence_chain,omitempty"`
	Evidence      string `json:"evidence,omitempty"`
	AIHint        string `json:"ai_hint,omitempty"`
	AINote        string `json:"ai_note,omitempty"`
}
