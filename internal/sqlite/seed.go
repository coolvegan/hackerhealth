// Package sqlite — Seed-Daten für die initiale Datenbank.
package sqlite

import "github.com/coolvegan/hn-mentalhealth"

// Seed befüllt die Datenbank mit den initialen Kategorien.
// Idempotent: INSERT OR IGNORE, kann mehrfach aufgerufen werden.
func (db *DB) Seed() error {
	if err := seedConfidenceLevels(db); err != nil {
		return err
	}
	if err := seedICD11Codes(db); err != nil {
		return err
	}
	if err := seedPatterns(db); err != nil {
		return err
	}
	if err := seedEmotions(db); err != nil {
		return err
	}
	if err := seedCharacterizations(db); err != nil {
		return err
	}
	if err := seedInteractions(db); err != nil {
		return err
	}
	if err := seedAIKnowledge(db); err != nil {
		return err
	}
	return nil
}

func seedConfidenceLevels(db *DB) error {
	levels := []hnmentalhealth.ConfidenceLevel{
		{Name: "secure", AIHint: "Hohe Konfidenz: Muster ist klar erkennbar, durch mehrere Evidenz-Quellen gestützt."},
		{Name: "suspected", AIHint: "Mittlere Konfidenz: Muster ist plausibel, aber nicht durch mehrere Quellen bestätigt."},
		{Name: "observation", AIHint: "Niedrige Konfidenz: Einzelbeobachtung, kein wiederholtes Muster. Nicht als Diagnose verwenden."},
		{Name: "healthy", AIHint: "Gesund/adaptiv: Der Kommentar zeigt gesunde Coping-Strategien oder konstruktive Interaktion."},
	}
	for i := range levels {
		if err := db.Shared.CreateConfidenceLevel(&levels[i]); err != nil {
			// INSERT OR IGNORE wäre sauberer, aber das Interface nutzt INSERT.
			// Wir ignorieren UNIQUE-Constraint-Fehler.
			_ = err
		}
	}
	return nil
}

func seedICD11Codes(db *DB) error {
	codes := []hnmentalhealth.ICD11Code{
		{Code: "6B40", Label: "Posttraumatic Stress Disorder", IsICD11: true, Note: "PTBS — Wiedererleben, Vermeidung, Übererregung nach traumatischem Ereignis."},
		{Code: "6B43", Label: "Adjustment Disorder", IsICD11: true, Note: "Anpassungsstörung — emotionale/verhaltensbezogene Symptome als Reaktion auf identifizierbaren Stressor."},
		{Code: "6B23", Label: "Generalized Anxiety Disorder", IsICD11: true, Note: "Generalisierte Angststörung — anhaltende, excessive Angst und Sorge."},
		{Code: "6B70", Label: "Depressive Episode", IsICD11: true, Note: "Depressive Episode — gedrückte Stimmung, Interessenverlust, Antriebsminderung."},
		{Code: "QE84", Label: "Burnout", IsICD11: true, Note: "Burnout — Faktor, der den Gesundheitsstatus beeinflusst (keine eigenständige Diagnose in ICD-11)."},
		{Code: "non-icd:moral_injury", Label: "Moral Injury", IsICD11: false, Note: "Moralische Verletzung — Schuld/Scham durch Handlungen, die eigene Werte verletzen. Kein ICD-11-Code."},
		{Code: "non-icd:attachment_trauma", Label: "Attachment/Developmental Trauma", IsICD11: false, Note: "Bindungs-/Entwicklungs-Trauma — frühe, chronische Beziehungsverletzungen. Kein eigenständiger ICD-11-Code."},
	}
	for i := range codes {
		_ = db.Shared.CreateICD11Code(&codes[i])
	}
	return nil
}

func seedPatterns(db *DB) error {
	patterns := []hnmentalhealth.Pattern{
		{Name: "workism", Description: "Überidentifikation mit Arbeit als primäre Identitätsquelle und Coping-Mechanismus."},
		{Name: "burnout", Description: "Emotionale Erschöpfung, Zynismus, reduzierte Leistungsfähigkeit durch chronischen Arbeitsstress."},
		{Name: "ai_angst", Description: "Angst vor Ersetzung durch KI, existenzielle Bedrohung der beruflichen Identität."},
		{Name: "sinnkrise", Description: "Verlust von Sinn und Bedeutung in der Arbeit oder im Leben."},
		{Name: "einsamkeit", Description: "Soziale Isolation, Mangel an bedeutungsvollen Beziehungen."},
		{Name: "vater_stimme", Description: "Internalisierte kritische Elternstimme, Perfektionismus, nie gut genug."},
		{Name: "flucht", Description: "Vermeidungsverhalten, Weglaufen vor Verantwortung oder schwierigen Situationen."},
		{Name: "scham", Description: "Toxische Scham, Gefühl fundamentaler Fehlerhaftigkeit."},
		{Name: "kontrollbeduerfnis", Description: "Übermäßiges Kontrollbedürfnis als Angstbewältigung."},
		{Name: "gesund", Description: "Gesunde/adaptive Muster: Selbstreflexion, emotionale Regulation, konstruktive Interaktion."},
	}
	for i := range patterns {
		_ = db.HN.CreatePattern(&patterns[i])
	}
	return nil
}

func seedEmotions(db *DB) error {
	emotions := []hnmentalhealth.Emotion{
		{Name: "joy", AIHint: "Freude, Begeisterung, positive Erregung."},
		{Name: "anger", AIHint: "Wut, Frustration, Empörung."},
		{Name: "sad", AIHint: "Trauer, Melancholie, Verlust."},
		{Name: "fear", AIHint: "Angst, Sorge, Bedrohung."},
		{Name: "shame", AIHint: "Scham, Peinlichkeit, Minderwertigkeit."},
	}
	for i := range emotions {
		_ = db.HN.CreateEmotion(&emotions[i])
	}
	return nil
}

func seedCharacterizations(db *DB) error {
	chars := []hnmentalhealth.Characterization{
		{Name: "values_debate", AIHint: "Gesunde Werte-Debatte: Teilnehmer tauschen Perspektiven aus, respektvoll, selbstreflektiert."},
		{Name: "grief_group", AIHint: "Trauergruppe: Geteilte Verlusterfahrung, gegenseitige Unterstützung, emotionale Offenheit."},
		{Name: "nostalgic_tech", AIHint: "Nostalgischer Technik-Talk: Erinnerungen an alte Technologien, gemeinsame positive Erfahrungen."},
		{Name: "competence_community", AIHint: "Kompetenz-Gemeinschaft: Fachaustausch auf Augenhöhe, konstruktives Feedback."},
		{Name: "existential_crisis", AIHint: "Existenzkrise: Sinnfragen, Lebensbilanz, fundamentale Verunsicherung."},
		{Name: "technical_discussion", AIHint: "Sachliche technische Diskussion: Faktenbasiert, lösungsorientiert, wenig emotionale Beteiligung."},
		{Name: "polarized", AIHint: "Polarisierung: Us-vs-Them-Dynamik, Identitätsverteidigung, wenig echte Diskussion."},
	}
	for i := range chars {
		_ = db.HN.CreateCharacterization(&chars[i])
	}
	return nil
}

func seedInteractions(db *DB) error {
	interactions := []hnmentalhealth.Interaction{
		{Name: "productive", AIHint: "Produktive Debatte: Selbstöffnung, Perspektivwechsel, konstruktiver Dissens."},
		{Name: "escalation", AIHint: "Eskalation: Persönliche Angriffe, aufschaukelnde Emotionen, Verletzungen."},
		{Name: "polarized", AIHint: "Polarisiert: Us-vs-Them, Verachtung, keine echte Kommunikation."},
		{Name: "flat", AIHint: "Flach: Keine Antwortketten, oberflächliche Interaktion, kein Engagement."},
	}
	for i := range interactions {
		_ = db.HN.CreateInteraction(&interactions[i])
	}
	return nil
}

func seedAIKnowledge(db *DB) error {
	entries := []hnmentalhealth.AIKnowledge{
		// Tabellen-Erklärungen
		{TableName: "users", Explanation: "Personen — HN-Kommentatoren und Patienten. Zentrale Identitäts-Tabelle, von HN und Health gemeinsam genutzt."},
		{TableName: "confidence_levels", Explanation: "Konfidenz-Stufen für Findings: secure (hoch), suspected (mittel), observation (niedrig), healthy (gesund)."},
		{TableName: "icd11_codes", Explanation: "ICD-11-Referenzen und nicht-ICD-Konstrukte. is_icd11=1 für echte ICD-11-Codes, 0 für klinische Konzepte ohne ICD-11-Status."},
		{TableName: "ai_knowledge", Explanation: "Selbst-beschreibendes Daten-Wörterbuch. Dokumentiert Tabellen, Spalten und Werte für die KI. Die KI kann hier abfragen und eigene Notizen hinterlegen."},
		{TableName: "patterns", Explanation: "Psychologische Muster, die in HN-Kommentaren erkannt werden (workism, burnout, scham, ...). Erweiterbar per INSERT."},
		{TableName: "emotions", Explanation: "Gefühle/Emotionen (joy, anger, sad, fear, shame). Werden pro Thread gezählt (thread_emotions)."},
		{TableName: "characterizations", Explanation: "Thread-Charakterisierungen: die qualitative Einordnung eines Threads (values_debate, grief_group, ...)."},
		{TableName: "interactions", Explanation: "Interaktions-Typen: wie Teilnehmer miteinander umgehen (productive, escalation, polarized, flat)."},
		{TableName: "threads", Explanation: "HN-Storys: rohe Metadaten (story_id, title, score, date, url)."},
		{TableName: "thread_quality", Explanation: "Qualitative Bewertung eines Threads: Charakterisierung, Interaktion, Trauma-Potential, healthy/notable Counts."},
		{TableName: "thread_emotions", Explanation: "Emotions-Zähler pro Thread: viele-zu-viele zwischen threads und emotions, mit count."},
		{TableName: "hn_findings", Explanation: "HN-Findings: eine Auffälligkeit eines Users in einem HN-Thread. Kernstruktur der HN-Analyse."},
		{TableName: "gene_snippets", Explanation: "Einzelne Gen-Snippets eines Patienten (Gen, Variante, Genotyp). Stufe 1 der Gen-Struktur."},
		{TableName: "gene_patterns", Explanation: "Gen-Muster: Aggregation mehrerer Gen-Snippets zu einem Stoffwechselweg/Krankheitsbild. Stufe 2."},
		{TableName: "pattern_snippets", Explanation: "Verknüpfung: welche Gen-Snippets gehören zu welchem Gen-Muster (viele-zu-viele)."},
		{TableName: "patient_patterns", Explanation: "Verknüpfung: welche Patienten haben welche Gen-Muster, mit eigener KI-Notiz pro Kombination."},
		{TableName: "biomarkers", Explanation: "Blutwerte/Laborwerte eines Patienten (Name, Wert, Einheit, Referenzbereich)."},
		{TableName: "health_conditions", Explanation: "Diagnosen/Gesundheitszustände eines Patienten (Diabetes, KPU/HPU, ...)."},
		{TableName: "life_context", Explanation: "Lebenskontext/Sorgen eines Patienten (Arbeitslosigkeit, Trennung, finanzielle Sorgen, ...)."},
		{TableName: "health_findings", Explanation: "Health-Findings: eine Auffälligkeit aus Gesundheitsdaten. Mit etiology (metabolic/psychological/mixed/unknown) und evidence_chain."},

		// Spalten-Erklärungen
		{TableName: "findings", ColumnName: "ai_hint", Explanation: "Hinweis für das AI-Modell: Kontext, der via MCP in den Prompt geht. Wird gelesen, nicht geschrieben."},
		{TableName: "findings", ColumnName: "ai_note", Explanation: "Notiz des AI-Modells: Erkenntnisse, die die KI selbst via MCP/SQL setzt. Wird geschrieben, nicht gelesen."},
		{TableName: "health_findings", ColumnName: "etiology", Explanation: "Ätiologie: metabolic (Stoffwechsel), psychological (psychologisch), mixed (gemischt), unknown (unbekannt)."},
		{TableName: "health_findings", ColumnName: "evidence_chain", Explanation: "JSON-Array der Evidenz-Quellen, z.B. ['MTHFR C677T', 'Folsäure 3.1 ng/ml', 'Studie X']."},
	}
	for i := range entries {
		_ = db.Shared.CreateAIKnowledge(&entries[i])
	}
	return nil
}
