// Package sqlite — Health Store (gene_snippets, gene_patterns, pattern_snippets,
// patient_patterns, biomarkers, health_conditions, life_context, health_findings).
package sqlite

import (
	"database/sql"

	"github.com/coolvegan/hn-mentalhealth"
)

// HealthStore implementiert repository.HealthRepository.
type HealthStore struct{ conn *sql.DB }

// --- GeneSnippetRepository ---

func (s *HealthStore) CreateGeneSnippet(gs *hnmentalhealth.GeneSnippet) error {
	r, err := s.conn.Exec(
		`INSERT INTO gene_snippets (patient_id, gene, variant, genotype, ai_hint, ai_note) VALUES (?, ?, ?, ?, ?, ?)`,
		gs.PatientID, gs.Gene, nullString(gs.Variant), nullString(gs.Genotype),
		nullString(gs.AIHint), nullString(gs.AINote),
	)
	if err != nil {
		return err
	}
	gs.ID, _ = r.LastInsertId()
	return nil
}

func (s *HealthStore) GetGeneSnippetByID(id int64) (*hnmentalhealth.GeneSnippet, error) {
	gs := &hnmentalhealth.GeneSnippet{}
	var variant, genotype, hint, note sql.NullString
	err := s.conn.QueryRow(
		`SELECT id, patient_id, gene, variant, genotype, ai_hint, ai_note FROM gene_snippets WHERE id = ?`, id,
	).Scan(&gs.ID, &gs.PatientID, &gs.Gene, &variant, &genotype, &hint, &note)
	if err != nil {
		return nil, err
	}
	gs.Variant = stringFromNull(variant)
	gs.Genotype = stringFromNull(genotype)
	gs.AIHint = stringFromNull(hint)
	gs.AINote = stringFromNull(note)
	return gs, nil
}

func (s *HealthStore) UpdateGeneSnippet(gs *hnmentalhealth.GeneSnippet) error {
	_, err := s.conn.Exec(
		`UPDATE gene_snippets SET patient_id = ?, gene = ?, variant = ?, genotype = ?, ai_hint = ?, ai_note = ? WHERE id = ?`,
		gs.PatientID, gs.Gene, nullString(gs.Variant), nullString(gs.Genotype),
		nullString(gs.AIHint), nullString(gs.AINote), gs.ID,
	)
	return err
}

func (s *HealthStore) DeleteGeneSnippet(id int64) error {
	_, err := s.conn.Exec(`DELETE FROM gene_snippets WHERE id = ?`, id)
	return err
}

func (s *HealthStore) ListGeneSnippetsByPatient(patientID int64) ([]hnmentalhealth.GeneSnippet, error) {
	rows, err := s.conn.Query(
		`SELECT id, patient_id, gene, variant, genotype, ai_hint, ai_note FROM gene_snippets WHERE patient_id = ? ORDER BY id`, patientID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var snippets []hnmentalhealth.GeneSnippet
	for rows.Next() {
		var gs hnmentalhealth.GeneSnippet
		var variant, genotype, hint, note sql.NullString
		if err := rows.Scan(&gs.ID, &gs.PatientID, &gs.Gene, &variant, &genotype, &hint, &note); err != nil {
			return nil, err
		}
		gs.Variant = stringFromNull(variant)
		gs.Genotype = stringFromNull(genotype)
		gs.AIHint = stringFromNull(hint)
		gs.AINote = stringFromNull(note)
		snippets = append(snippets, gs)
	}
	return snippets, rows.Err()
}

// --- GenePatternRepository ---

func (s *HealthStore) CreateGenePattern(gp *hnmentalhealth.GenePattern) error {
	r, err := s.conn.Exec(
		`INSERT INTO gene_patterns (name, description, ai_hint, ai_note) VALUES (?, ?, ?, ?)`,
		gp.Name, nullString(gp.Description), nullString(gp.AIHint), nullString(gp.AINote),
	)
	if err != nil {
		return err
	}
	gp.ID, _ = r.LastInsertId()
	return nil
}

func (s *HealthStore) GetGenePatternByID(id int64) (*hnmentalhealth.GenePattern, error) {
	gp := &hnmentalhealth.GenePattern{}
	var desc, hint, note sql.NullString
	err := s.conn.QueryRow(
		`SELECT id, name, description, ai_hint, ai_note FROM gene_patterns WHERE id = ?`, id,
	).Scan(&gp.ID, &gp.Name, &desc, &hint, &note)
	if err != nil {
		return nil, err
	}
	gp.Description = stringFromNull(desc)
	gp.AIHint = stringFromNull(hint)
	gp.AINote = stringFromNull(note)
	return gp, nil
}

func (s *HealthStore) GetGenePatternByName(name string) (*hnmentalhealth.GenePattern, error) {
	gp := &hnmentalhealth.GenePattern{}
	var desc, hint, note sql.NullString
	err := s.conn.QueryRow(
		`SELECT id, name, description, ai_hint, ai_note FROM gene_patterns WHERE name = ?`, name,
	).Scan(&gp.ID, &gp.Name, &desc, &hint, &note)
	if err != nil {
		return nil, err
	}
	gp.Description = stringFromNull(desc)
	gp.AIHint = stringFromNull(hint)
	gp.AINote = stringFromNull(note)
	return gp, nil
}

func (s *HealthStore) UpdateGenePattern(gp *hnmentalhealth.GenePattern) error {
	_, err := s.conn.Exec(
		`UPDATE gene_patterns SET name = ?, description = ?, ai_hint = ?, ai_note = ? WHERE id = ?`,
		gp.Name, nullString(gp.Description), nullString(gp.AIHint), nullString(gp.AINote), gp.ID,
	)
	return err
}

func (s *HealthStore) DeleteGenePattern(id int64) error {
	_, err := s.conn.Exec(`DELETE FROM gene_patterns WHERE id = ?`, id)
	return err
}

func (s *HealthStore) ListGenePatterns() ([]hnmentalhealth.GenePattern, error) {
	rows, err := s.conn.Query(`SELECT id, name, description, ai_hint, ai_note FROM gene_patterns ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var patterns []hnmentalhealth.GenePattern
	for rows.Next() {
		var gp hnmentalhealth.GenePattern
		var desc, hint, note sql.NullString
		if err := rows.Scan(&gp.ID, &gp.Name, &desc, &hint, &note); err != nil {
			return nil, err
		}
		gp.Description = stringFromNull(desc)
		gp.AIHint = stringFromNull(hint)
		gp.AINote = stringFromNull(note)
		patterns = append(patterns, gp)
	}
	return patterns, rows.Err()
}

// --- PatternSnippetRepository ---

func (s *HealthStore) AddSnippetToPattern(patternID, snippetID int64) error {
	_, err := s.conn.Exec(
		`INSERT OR IGNORE INTO pattern_snippets (pattern_id, snippet_id) VALUES (?, ?)`,
		patternID, snippetID,
	)
	return err
}

func (s *HealthStore) RemoveSnippetFromPattern(patternID, snippetID int64) error {
	_, err := s.conn.Exec(
		`DELETE FROM pattern_snippets WHERE pattern_id = ? AND snippet_id = ?`,
		patternID, snippetID,
	)
	return err
}

func (s *HealthStore) ListSnippetsByPattern(patternID int64) ([]hnmentalhealth.GeneSnippet, error) {
	rows, err := s.conn.Query(
		`SELECT gs.id, gs.patient_id, gs.gene, gs.variant, gs.genotype, gs.ai_hint, gs.ai_note
		 FROM gene_snippets gs
		 JOIN pattern_snippets ps ON gs.id = ps.snippet_id
		 WHERE ps.pattern_id = ?
		 ORDER BY gs.id`, patternID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var snippets []hnmentalhealth.GeneSnippet
	for rows.Next() {
		var gs hnmentalhealth.GeneSnippet
		var variant, genotype, hint, note sql.NullString
		if err := rows.Scan(&gs.ID, &gs.PatientID, &gs.Gene, &variant, &genotype, &hint, &note); err != nil {
			return nil, err
		}
		gs.Variant = stringFromNull(variant)
		gs.Genotype = stringFromNull(genotype)
		gs.AIHint = stringFromNull(hint)
		gs.AINote = stringFromNull(note)
		snippets = append(snippets, gs)
	}
	return snippets, rows.Err()
}

func (s *HealthStore) ListPatternsBySnippet(snippetID int64) ([]hnmentalhealth.GenePattern, error) {
	rows, err := s.conn.Query(
		`SELECT gp.id, gp.name, gp.description, gp.ai_hint, gp.ai_note
		 FROM gene_patterns gp
		 JOIN pattern_snippets ps ON gp.id = ps.pattern_id
		 WHERE ps.snippet_id = ?
		 ORDER BY gp.id`, snippetID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var patterns []hnmentalhealth.GenePattern
	for rows.Next() {
		var gp hnmentalhealth.GenePattern
		var desc, hint, note sql.NullString
		if err := rows.Scan(&gp.ID, &gp.Name, &desc, &hint, &note); err != nil {
			return nil, err
		}
		gp.Description = stringFromNull(desc)
		gp.AIHint = stringFromNull(hint)
		gp.AINote = stringFromNull(note)
		patterns = append(patterns, gp)
	}
	return patterns, rows.Err()
}

// --- PatientPatternRepository ---

func (s *HealthStore) AddPatternToPatient(patientID, patternID int64) error {
	_, err := s.conn.Exec(
		`INSERT OR IGNORE INTO patient_patterns (patient_id, pattern_id) VALUES (?, ?)`,
		patientID, patternID,
	)
	return err
}

func (s *HealthStore) RemovePatternFromPatient(patientID, patternID int64) error {
	_, err := s.conn.Exec(
		`DELETE FROM patient_patterns WHERE patient_id = ? AND pattern_id = ?`,
		patientID, patternID,
	)
	return err
}

func (s *HealthStore) UpdatePatientPatternNote(patientID, patternID int64, note string) error {
	_, err := s.conn.Exec(
		`UPDATE patient_patterns SET ai_note = ? WHERE patient_id = ? AND pattern_id = ?`,
		nullString(note), patientID, patternID,
	)
	return err
}

func (s *HealthStore) ListPatternsByPatient(patientID int64) ([]hnmentalhealth.GenePattern, error) {
	rows, err := s.conn.Query(
		`SELECT gp.id, gp.name, gp.description, gp.ai_hint, gp.ai_note
		 FROM gene_patterns gp
		 JOIN patient_patterns pp ON gp.id = pp.pattern_id
		 WHERE pp.patient_id = ?
		 ORDER BY gp.id`, patientID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var patterns []hnmentalhealth.GenePattern
	for rows.Next() {
		var gp hnmentalhealth.GenePattern
		var desc, hint, note sql.NullString
		if err := rows.Scan(&gp.ID, &gp.Name, &desc, &hint, &note); err != nil {
			return nil, err
		}
		gp.Description = stringFromNull(desc)
		gp.AIHint = stringFromNull(hint)
		gp.AINote = stringFromNull(note)
		patterns = append(patterns, gp)
	}
	return patterns, rows.Err()
}

func (s *HealthStore) ListPatientsByPattern(patternID int64) ([]hnmentalhealth.User, error) {
	rows, err := s.conn.Query(
		`SELECT u.id, u.username, u.first_seen_at, u.thread_count
		 FROM users u
		 JOIN patient_patterns pp ON u.id = pp.patient_id
		 WHERE pp.pattern_id = ?
		 ORDER BY u.id`, patternID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []hnmentalhealth.User
	for rows.Next() {
		var u hnmentalhealth.User
		var firstSeen sql.NullString
		if err := rows.Scan(&u.ID, &u.Username, &firstSeen, &u.ThreadCount); err != nil {
			return nil, err
		}
		u.FirstSeenAt = stringFromNull(firstSeen)
		users = append(users, u)
	}
	return users, rows.Err()
}

// --- BiomarkerRepository ---

func (s *HealthStore) CreateBiomarker(b *hnmentalhealth.Biomarker) error {
	r, err := s.conn.Exec(
		`INSERT INTO biomarkers (patient_id, name, value, unit, reference_low, reference_high, measured_at, ai_hint, ai_note) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		b.PatientID, b.Name, b.Value, nullString(b.Unit), b.ReferenceLow, b.ReferenceHigh,
		nullString(b.MeasuredAt), nullString(b.AIHint), nullString(b.AINote),
	)
	if err != nil {
		return err
	}
	b.ID, _ = r.LastInsertId()
	return nil
}

func (s *HealthStore) GetBiomarkerByID(id int64) (*hnmentalhealth.Biomarker, error) {
	b := &hnmentalhealth.Biomarker{}
	var unit, measured, hint, note sql.NullString
	err := s.conn.QueryRow(
		`SELECT id, patient_id, name, value, unit, reference_low, reference_high, measured_at, ai_hint, ai_note FROM biomarkers WHERE id = ?`, id,
	).Scan(&b.ID, &b.PatientID, &b.Name, &b.Value, &unit, &b.ReferenceLow, &b.ReferenceHigh, &measured, &hint, &note)
	if err != nil {
		return nil, err
	}
	b.Unit = stringFromNull(unit)
	b.MeasuredAt = stringFromNull(measured)
	b.AIHint = stringFromNull(hint)
	b.AINote = stringFromNull(note)
	return b, nil
}

func (s *HealthStore) UpdateBiomarker(b *hnmentalhealth.Biomarker) error {
	_, err := s.conn.Exec(
		`UPDATE biomarkers SET patient_id = ?, name = ?, value = ?, unit = ?, reference_low = ?, reference_high = ?, measured_at = ?, ai_hint = ?, ai_note = ? WHERE id = ?`,
		b.PatientID, b.Name, b.Value, nullString(b.Unit), b.ReferenceLow, b.ReferenceHigh,
		nullString(b.MeasuredAt), nullString(b.AIHint), nullString(b.AINote), b.ID,
	)
	return err
}

func (s *HealthStore) DeleteBiomarker(id int64) error {
	_, err := s.conn.Exec(`DELETE FROM biomarkers WHERE id = ?`, id)
	return err
}

func (s *HealthStore) ListBiomarkersByPatient(patientID int64) ([]hnmentalhealth.Biomarker, error) {
	rows, err := s.conn.Query(
		`SELECT id, patient_id, name, value, unit, reference_low, reference_high, measured_at, ai_hint, ai_note FROM biomarkers WHERE patient_id = ? ORDER BY id`, patientID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanBiomarkers(rows)
}

func (s *HealthStore) ListBiomarkersOutOfRange(patientID int64) ([]hnmentalhealth.Biomarker, error) {
	rows, err := s.conn.Query(
		`SELECT id, patient_id, name, value, unit, reference_low, reference_high, measured_at, ai_hint, ai_note
		 FROM biomarkers
		 WHERE patient_id = ? AND (value < reference_low OR value > reference_high)
		 ORDER BY id`, patientID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanBiomarkers(rows)
}

func scanBiomarkers(rows *sql.Rows) ([]hnmentalhealth.Biomarker, error) {
	var biomarkers []hnmentalhealth.Biomarker
	for rows.Next() {
		var b hnmentalhealth.Biomarker
		var unit, measured, hint, note sql.NullString
		if err := rows.Scan(&b.ID, &b.PatientID, &b.Name, &b.Value, &unit, &b.ReferenceLow, &b.ReferenceHigh, &measured, &hint, &note); err != nil {
			return nil, err
		}
		b.Unit = stringFromNull(unit)
		b.MeasuredAt = stringFromNull(measured)
		b.AIHint = stringFromNull(hint)
		b.AINote = stringFromNull(note)
		biomarkers = append(biomarkers, b)
	}
	return biomarkers, rows.Err()
}

// --- HealthConditionRepository ---

func (s *HealthStore) CreateHealthCondition(hc *hnmentalhealth.HealthCondition) error {
	r, err := s.conn.Exec(
		`INSERT INTO health_conditions (patient_id, condition, diagnosed_at, ai_hint, ai_note) VALUES (?, ?, ?, ?, ?)`,
		hc.PatientID, hc.Condition, nullString(hc.DiagnosedAt), nullString(hc.AIHint), nullString(hc.AINote),
	)
	if err != nil {
		return err
	}
	hc.ID, _ = r.LastInsertId()
	return nil
}

func (s *HealthStore) GetHealthConditionByID(id int64) (*hnmentalhealth.HealthCondition, error) {
	hc := &hnmentalhealth.HealthCondition{}
	var diagnosed, hint, note sql.NullString
	err := s.conn.QueryRow(
		`SELECT id, patient_id, condition, diagnosed_at, ai_hint, ai_note FROM health_conditions WHERE id = ?`, id,
	).Scan(&hc.ID, &hc.PatientID, &hc.Condition, &diagnosed, &hint, &note)
	if err != nil {
		return nil, err
	}
	hc.DiagnosedAt = stringFromNull(diagnosed)
	hc.AIHint = stringFromNull(hint)
	hc.AINote = stringFromNull(note)
	return hc, nil
}

func (s *HealthStore) UpdateHealthCondition(hc *hnmentalhealth.HealthCondition) error {
	_, err := s.conn.Exec(
		`UPDATE health_conditions SET patient_id = ?, condition = ?, diagnosed_at = ?, ai_hint = ?, ai_note = ? WHERE id = ?`,
		hc.PatientID, hc.Condition, nullString(hc.DiagnosedAt), nullString(hc.AIHint), nullString(hc.AINote), hc.ID,
	)
	return err
}

func (s *HealthStore) DeleteHealthCondition(id int64) error {
	_, err := s.conn.Exec(`DELETE FROM health_conditions WHERE id = ?`, id)
	return err
}

func (s *HealthStore) ListHealthConditionsByPatient(patientID int64) ([]hnmentalhealth.HealthCondition, error) {
	rows, err := s.conn.Query(
		`SELECT id, patient_id, condition, diagnosed_at, ai_hint, ai_note FROM health_conditions WHERE patient_id = ? ORDER BY id`, patientID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var conditions []hnmentalhealth.HealthCondition
	for rows.Next() {
		var hc hnmentalhealth.HealthCondition
		var diagnosed, hint, note sql.NullString
		if err := rows.Scan(&hc.ID, &hc.PatientID, &hc.Condition, &diagnosed, &hint, &note); err != nil {
			return nil, err
		}
		hc.DiagnosedAt = stringFromNull(diagnosed)
		hc.AIHint = stringFromNull(hint)
		hc.AINote = stringFromNull(note)
		conditions = append(conditions, hc)
	}
	return conditions, rows.Err()
}

// --- LifeContextRepository ---

func (s *HealthStore) CreateLifeContext(lc *hnmentalhealth.LifeContext) error {
	r, err := s.conn.Exec(
		`INSERT INTO life_context (patient_id, context, noted_at, ai_hint, ai_note) VALUES (?, ?, ?, ?, ?)`,
		lc.PatientID, lc.Context, nullString(lc.NotedAt), nullString(lc.AIHint), nullString(lc.AINote),
	)
	if err != nil {
		return err
	}
	lc.ID, _ = r.LastInsertId()
	return nil
}

func (s *HealthStore) GetLifeContextByID(id int64) (*hnmentalhealth.LifeContext, error) {
	lc := &hnmentalhealth.LifeContext{}
	var noted, hint, note sql.NullString
	err := s.conn.QueryRow(
		`SELECT id, patient_id, context, noted_at, ai_hint, ai_note FROM life_context WHERE id = ?`, id,
	).Scan(&lc.ID, &lc.PatientID, &lc.Context, &noted, &hint, &note)
	if err != nil {
		return nil, err
	}
	lc.NotedAt = stringFromNull(noted)
	lc.AIHint = stringFromNull(hint)
	lc.AINote = stringFromNull(note)
	return lc, nil
}

func (s *HealthStore) UpdateLifeContext(lc *hnmentalhealth.LifeContext) error {
	_, err := s.conn.Exec(
		`UPDATE life_context SET patient_id = ?, context = ?, noted_at = ?, ai_hint = ?, ai_note = ? WHERE id = ?`,
		lc.PatientID, lc.Context, nullString(lc.NotedAt), nullString(lc.AIHint), nullString(lc.AINote), lc.ID,
	)
	return err
}

func (s *HealthStore) DeleteLifeContext(id int64) error {
	_, err := s.conn.Exec(`DELETE FROM life_context WHERE id = ?`, id)
	return err
}

func (s *HealthStore) ListLifeContextByPatient(patientID int64) ([]hnmentalhealth.LifeContext, error) {
	rows, err := s.conn.Query(
		`SELECT id, patient_id, context, noted_at, ai_hint, ai_note FROM life_context WHERE patient_id = ? ORDER BY id`, patientID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var contexts []hnmentalhealth.LifeContext
	for rows.Next() {
		var lc hnmentalhealth.LifeContext
		var noted, hint, note sql.NullString
		if err := rows.Scan(&lc.ID, &lc.PatientID, &lc.Context, &noted, &hint, &note); err != nil {
			return nil, err
		}
		lc.NotedAt = stringFromNull(noted)
		lc.AIHint = stringFromNull(hint)
		lc.AINote = stringFromNull(note)
		contexts = append(contexts, lc)
	}
	return contexts, rows.Err()
}

// --- HealthFindingRepository ---

func (s *HealthStore) CreateHealthFinding(f *hnmentalhealth.HealthFinding) error {
	r, err := s.conn.Exec(
		`INSERT INTO health_findings (patient_id, pattern_id, confidence_id, icd11_id, etiology, evidence_chain, evidence, ai_hint, ai_note) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		f.PatientID, nullInt64Ptr(f.PatternID), nullInt64Ptr(f.ConfidenceID),
		nullInt64Ptr(f.ICD11ID), nullString(f.Etiology), nullString(f.EvidenceChain),
		nullString(f.Evidence), nullString(f.AIHint), nullString(f.AINote),
	)
	if err != nil {
		return err
	}
	f.ID, _ = r.LastInsertId()
	return nil
}

func (s *HealthStore) GetHealthFindingByID(id int64) (*hnmentalhealth.HealthFinding, error) {
	f := &hnmentalhealth.HealthFinding{}
	var patternID, confID, icd11ID sql.NullInt64
	var etiology, chain, evidence, hint, note sql.NullString
	err := s.conn.QueryRow(
		`SELECT id, patient_id, pattern_id, confidence_id, icd11_id, etiology, evidence_chain, evidence, ai_hint, ai_note FROM health_findings WHERE id = ?`, id,
	).Scan(&f.ID, &f.PatientID, &patternID, &confID, &icd11ID, &etiology, &chain, &evidence, &hint, &note)
	if err != nil {
		return nil, err
	}
	f.PatternID = int64PtrFromNull(patternID)
	f.ConfidenceID = int64PtrFromNull(confID)
	f.ICD11ID = int64PtrFromNull(icd11ID)
	f.Etiology = stringFromNull(etiology)
	f.EvidenceChain = stringFromNull(chain)
	f.Evidence = stringFromNull(evidence)
	f.AIHint = stringFromNull(hint)
	f.AINote = stringFromNull(note)
	return f, nil
}

func (s *HealthStore) UpdateHealthFinding(f *hnmentalhealth.HealthFinding) error {
	_, err := s.conn.Exec(
		`UPDATE health_findings SET patient_id = ?, pattern_id = ?, confidence_id = ?, icd11_id = ?, etiology = ?, evidence_chain = ?, evidence = ?, ai_hint = ?, ai_note = ? WHERE id = ?`,
		f.PatientID, nullInt64Ptr(f.PatternID), nullInt64Ptr(f.ConfidenceID),
		nullInt64Ptr(f.ICD11ID), nullString(f.Etiology), nullString(f.EvidenceChain),
		nullString(f.Evidence), nullString(f.AIHint), nullString(f.AINote), f.ID,
	)
	return err
}

func (s *HealthStore) DeleteHealthFinding(id int64) error {
	_, err := s.conn.Exec(`DELETE FROM health_findings WHERE id = ?`, id)
	return err
}

func (s *HealthStore) ListHealthFindingsByPatient(patientID int64) ([]hnmentalhealth.HealthFinding, error) {
	rows, err := s.conn.Query(
		`SELECT id, patient_id, pattern_id, confidence_id, icd11_id, etiology, evidence_chain, evidence, ai_hint, ai_note FROM health_findings WHERE patient_id = ? ORDER BY id`, patientID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanHealthFindings(rows)
}

func (s *HealthStore) ListHealthFindingsByPattern(patternID int64) ([]hnmentalhealth.HealthFinding, error) {
	rows, err := s.conn.Query(
		`SELECT id, patient_id, pattern_id, confidence_id, icd11_id, etiology, evidence_chain, evidence, ai_hint, ai_note FROM health_findings WHERE pattern_id = ? ORDER BY id`, patternID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanHealthFindings(rows)
}

func (s *HealthStore) ListHealthFindingsByEtiology(etiology string) ([]hnmentalhealth.HealthFinding, error) {
	rows, err := s.conn.Query(
		`SELECT id, patient_id, pattern_id, confidence_id, icd11_id, etiology, evidence_chain, evidence, ai_hint, ai_note FROM health_findings WHERE etiology = ? ORDER BY id`, etiology,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanHealthFindings(rows)
}

func scanHealthFindings(rows *sql.Rows) ([]hnmentalhealth.HealthFinding, error) {
	var findings []hnmentalhealth.HealthFinding
	for rows.Next() {
		var f hnmentalhealth.HealthFinding
		var patternID, confID, icd11ID sql.NullInt64
		var etiology, chain, evidence, hint, note sql.NullString
		if err := rows.Scan(&f.ID, &f.PatientID, &patternID, &confID, &icd11ID, &etiology, &chain, &evidence, &hint, &note); err != nil {
			return nil, err
		}
		f.PatternID = int64PtrFromNull(patternID)
		f.ConfidenceID = int64PtrFromNull(confID)
		f.ICD11ID = int64PtrFromNull(icd11ID)
		f.Etiology = stringFromNull(etiology)
		f.EvidenceChain = stringFromNull(chain)
		f.Evidence = stringFromNull(evidence)
		f.AIHint = stringFromNull(hint)
		f.AINote = stringFromNull(note)
		findings = append(findings, f)
	}
	return findings, rows.Err()
}
