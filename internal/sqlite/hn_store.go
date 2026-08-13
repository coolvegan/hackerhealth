// Package sqlite — HN Store (patterns, emotions, characterizations, interactions,
// threads, thread_quality, thread_emotions, hn_findings).
package sqlite

import (
	"database/sql"

	"github.com/coolvegan/hn-mentalhealth"
)

// HNStore implementiert repository.HNRepository.
type HNStore struct{ conn *sql.DB }

// --- PatternRepository ---

func (s *HNStore) CreatePattern(p *hnmentalhealth.Pattern) error {
	r, err := s.conn.Exec(
		`INSERT INTO patterns (name, description, ai_hint, ai_note) VALUES (?, ?, ?, ?)`,
		p.Name, nullString(p.Description), nullString(p.AIHint), nullString(p.AINote),
	)
	if err != nil {
		return err
	}
	p.ID, _ = r.LastInsertId()
	return nil
}

func (s *HNStore) GetPatternByID(id int64) (*hnmentalhealth.Pattern, error) {
	p := &hnmentalhealth.Pattern{}
	var desc, hint, note sql.NullString
	err := s.conn.QueryRow(
		`SELECT id, name, description, ai_hint, ai_note FROM patterns WHERE id = ?`, id,
	).Scan(&p.ID, &p.Name, &desc, &hint, &note)
	if err != nil {
		return nil, err
	}
	p.Description = stringFromNull(desc)
	p.AIHint = stringFromNull(hint)
	p.AINote = stringFromNull(note)
	return p, nil
}

func (s *HNStore) GetPatternByName(name string) (*hnmentalhealth.Pattern, error) {
	p := &hnmentalhealth.Pattern{}
	var desc, hint, note sql.NullString
	err := s.conn.QueryRow(
		`SELECT id, name, description, ai_hint, ai_note FROM patterns WHERE name = ?`, name,
	).Scan(&p.ID, &p.Name, &desc, &hint, &note)
	if err != nil {
		return nil, err
	}
	p.Description = stringFromNull(desc)
	p.AIHint = stringFromNull(hint)
	p.AINote = stringFromNull(note)
	return p, nil
}

func (s *HNStore) UpdatePattern(p *hnmentalhealth.Pattern) error {
	_, err := s.conn.Exec(
		`UPDATE patterns SET name = ?, description = ?, ai_hint = ?, ai_note = ? WHERE id = ?`,
		p.Name, nullString(p.Description), nullString(p.AIHint), nullString(p.AINote), p.ID,
	)
	return err
}

func (s *HNStore) DeletePattern(id int64) error {
	_, err := s.conn.Exec(`DELETE FROM patterns WHERE id = ?`, id)
	return err
}

func (s *HNStore) ListPatterns() ([]hnmentalhealth.Pattern, error) {
	rows, err := s.conn.Query(`SELECT id, name, description, ai_hint, ai_note FROM patterns ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var patterns []hnmentalhealth.Pattern
	for rows.Next() {
		var p hnmentalhealth.Pattern
		var desc, hint, note sql.NullString
		if err := rows.Scan(&p.ID, &p.Name, &desc, &hint, &note); err != nil {
			return nil, err
		}
		p.Description = stringFromNull(desc)
		p.AIHint = stringFromNull(hint)
		p.AINote = stringFromNull(note)
		patterns = append(patterns, p)
	}
	return patterns, rows.Err()
}

// --- EmotionRepository ---

func (s *HNStore) CreateEmotion(e *hnmentalhealth.Emotion) error {
	r, err := s.conn.Exec(
		`INSERT INTO emotions (name, ai_hint, ai_note) VALUES (?, ?, ?)`,
		e.Name, nullString(e.AIHint), nullString(e.AINote),
	)
	if err != nil {
		return err
	}
	e.ID, _ = r.LastInsertId()
	return nil
}

func (s *HNStore) GetEmotionByID(id int64) (*hnmentalhealth.Emotion, error) {
	e := &hnmentalhealth.Emotion{}
	var hint, note sql.NullString
	err := s.conn.QueryRow(
		`SELECT id, name, ai_hint, ai_note FROM emotions WHERE id = ?`, id,
	).Scan(&e.ID, &e.Name, &hint, &note)
	if err != nil {
		return nil, err
	}
	e.AIHint = stringFromNull(hint)
	e.AINote = stringFromNull(note)
	return e, nil
}

func (s *HNStore) GetEmotionByName(name string) (*hnmentalhealth.Emotion, error) {
	e := &hnmentalhealth.Emotion{}
	var hint, note sql.NullString
	err := s.conn.QueryRow(
		`SELECT id, name, ai_hint, ai_note FROM emotions WHERE name = ?`, name,
	).Scan(&e.ID, &e.Name, &hint, &note)
	if err != nil {
		return nil, err
	}
	e.AIHint = stringFromNull(hint)
	e.AINote = stringFromNull(note)
	return e, nil
}

func (s *HNStore) UpdateEmotion(e *hnmentalhealth.Emotion) error {
	_, err := s.conn.Exec(
		`UPDATE emotions SET name = ?, ai_hint = ?, ai_note = ? WHERE id = ?`,
		e.Name, nullString(e.AIHint), nullString(e.AINote), e.ID,
	)
	return err
}

func (s *HNStore) DeleteEmotion(id int64) error {
	_, err := s.conn.Exec(`DELETE FROM emotions WHERE id = ?`, id)
	return err
}

func (s *HNStore) ListEmotions() ([]hnmentalhealth.Emotion, error) {
	rows, err := s.conn.Query(`SELECT id, name, ai_hint, ai_note FROM emotions ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var emotions []hnmentalhealth.Emotion
	for rows.Next() {
		var e hnmentalhealth.Emotion
		var hint, note sql.NullString
		if err := rows.Scan(&e.ID, &e.Name, &hint, &note); err != nil {
			return nil, err
		}
		e.AIHint = stringFromNull(hint)
		e.AINote = stringFromNull(note)
		emotions = append(emotions, e)
	}
	return emotions, rows.Err()
}

// --- CharacterizationRepository ---

func (s *HNStore) CreateCharacterization(c *hnmentalhealth.Characterization) error {
	r, err := s.conn.Exec(
		`INSERT INTO characterizations (name, ai_hint, ai_note) VALUES (?, ?, ?)`,
		c.Name, nullString(c.AIHint), nullString(c.AINote),
	)
	if err != nil {
		return err
	}
	c.ID, _ = r.LastInsertId()
	return nil
}

func (s *HNStore) GetCharacterizationByID(id int64) (*hnmentalhealth.Characterization, error) {
	c := &hnmentalhealth.Characterization{}
	var hint, note sql.NullString
	err := s.conn.QueryRow(
		`SELECT id, name, ai_hint, ai_note FROM characterizations WHERE id = ?`, id,
	).Scan(&c.ID, &c.Name, &hint, &note)
	if err != nil {
		return nil, err
	}
	c.AIHint = stringFromNull(hint)
	c.AINote = stringFromNull(note)
	return c, nil
}

func (s *HNStore) GetCharacterizationByName(name string) (*hnmentalhealth.Characterization, error) {
	c := &hnmentalhealth.Characterization{}
	var hint, note sql.NullString
	err := s.conn.QueryRow(
		`SELECT id, name, ai_hint, ai_note FROM characterizations WHERE name = ?`, name,
	).Scan(&c.ID, &c.Name, &hint, &note)
	if err != nil {
		return nil, err
	}
	c.AIHint = stringFromNull(hint)
	c.AINote = stringFromNull(note)
	return c, nil
}

func (s *HNStore) UpdateCharacterization(c *hnmentalhealth.Characterization) error {
	_, err := s.conn.Exec(
		`UPDATE characterizations SET name = ?, ai_hint = ?, ai_note = ? WHERE id = ?`,
		c.Name, nullString(c.AIHint), nullString(c.AINote), c.ID,
	)
	return err
}

func (s *HNStore) DeleteCharacterization(id int64) error {
	_, err := s.conn.Exec(`DELETE FROM characterizations WHERE id = ?`, id)
	return err
}

func (s *HNStore) ListCharacterizations() ([]hnmentalhealth.Characterization, error) {
	rows, err := s.conn.Query(`SELECT id, name, ai_hint, ai_note FROM characterizations ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var chars []hnmentalhealth.Characterization
	for rows.Next() {
		var c hnmentalhealth.Characterization
		var hint, note sql.NullString
		if err := rows.Scan(&c.ID, &c.Name, &hint, &note); err != nil {
			return nil, err
		}
		c.AIHint = stringFromNull(hint)
		c.AINote = stringFromNull(note)
		chars = append(chars, c)
	}
	return chars, rows.Err()
}

// --- InteractionRepository ---

func (s *HNStore) CreateInteraction(i *hnmentalhealth.Interaction) error {
	r, err := s.conn.Exec(
		`INSERT INTO interactions (name, ai_hint, ai_note) VALUES (?, ?, ?)`,
		i.Name, nullString(i.AIHint), nullString(i.AINote),
	)
	if err != nil {
		return err
	}
	i.ID, _ = r.LastInsertId()
	return nil
}

func (s *HNStore) GetInteractionByID(id int64) (*hnmentalhealth.Interaction, error) {
	i := &hnmentalhealth.Interaction{}
	var hint, note sql.NullString
	err := s.conn.QueryRow(
		`SELECT id, name, ai_hint, ai_note FROM interactions WHERE id = ?`, id,
	).Scan(&i.ID, &i.Name, &hint, &note)
	if err != nil {
		return nil, err
	}
	i.AIHint = stringFromNull(hint)
	i.AINote = stringFromNull(note)
	return i, nil
}

func (s *HNStore) GetInteractionByName(name string) (*hnmentalhealth.Interaction, error) {
	i := &hnmentalhealth.Interaction{}
	var hint, note sql.NullString
	err := s.conn.QueryRow(
		`SELECT id, name, ai_hint, ai_note FROM interactions WHERE name = ?`, name,
	).Scan(&i.ID, &i.Name, &hint, &note)
	if err != nil {
		return nil, err
	}
	i.AIHint = stringFromNull(hint)
	i.AINote = stringFromNull(note)
	return i, nil
}

func (s *HNStore) UpdateInteraction(i *hnmentalhealth.Interaction) error {
	_, err := s.conn.Exec(
		`UPDATE interactions SET name = ?, ai_hint = ?, ai_note = ? WHERE id = ?`,
		i.Name, nullString(i.AIHint), nullString(i.AINote), i.ID,
	)
	return err
}

func (s *HNStore) DeleteInteraction(id int64) error {
	_, err := s.conn.Exec(`DELETE FROM interactions WHERE id = ?`, id)
	return err
}

func (s *HNStore) ListInteractions() ([]hnmentalhealth.Interaction, error) {
	rows, err := s.conn.Query(`SELECT id, name, ai_hint, ai_note FROM interactions ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var interactions []hnmentalhealth.Interaction
	for rows.Next() {
		var i hnmentalhealth.Interaction
		var hint, note sql.NullString
		if err := rows.Scan(&i.ID, &i.Name, &hint, &note); err != nil {
			return nil, err
		}
		i.AIHint = stringFromNull(hint)
		i.AINote = stringFromNull(note)
		interactions = append(interactions, i)
	}
	return interactions, rows.Err()
}

// --- ThreadRepository ---

func (s *HNStore) CreateThread(t *hnmentalhealth.Thread) error {
	r, err := s.conn.Exec(
		`INSERT INTO threads (story_id, title, score, date, url) VALUES (?, ?, ?, ?, ?)`,
		t.StoryID, nullString(t.Title), t.Score, nullString(t.Date), nullString(t.URL),
	)
	if err != nil {
		return err
	}
	t.ID, _ = r.LastInsertId()
	return nil
}

func (s *HNStore) GetThreadByID(id int64) (*hnmentalhealth.Thread, error) {
	t := &hnmentalhealth.Thread{}
	var title, date, url sql.NullString
	err := s.conn.QueryRow(
		`SELECT id, story_id, title, score, date, url FROM threads WHERE id = ?`, id,
	).Scan(&t.ID, &t.StoryID, &title, &t.Score, &date, &url)
	if err != nil {
		return nil, err
	}
	t.Title = stringFromNull(title)
	t.Date = stringFromNull(date)
	t.URL = stringFromNull(url)
	return t, nil
}

func (s *HNStore) GetThreadByStoryID(storyID int64) (*hnmentalhealth.Thread, error) {
	t := &hnmentalhealth.Thread{}
	var title, date, url sql.NullString
	err := s.conn.QueryRow(
		`SELECT id, story_id, title, score, date, url FROM threads WHERE story_id = ?`, storyID,
	).Scan(&t.ID, &t.StoryID, &title, &t.Score, &date, &url)
	if err != nil {
		return nil, err
	}
	t.Title = stringFromNull(title)
	t.Date = stringFromNull(date)
	t.URL = stringFromNull(url)
	return t, nil
}

func (s *HNStore) UpdateThread(t *hnmentalhealth.Thread) error {
	_, err := s.conn.Exec(
		`UPDATE threads SET story_id = ?, title = ?, score = ?, date = ?, url = ? WHERE id = ?`,
		t.StoryID, nullString(t.Title), t.Score, nullString(t.Date), nullString(t.URL), t.ID,
	)
	return err
}

func (s *HNStore) DeleteThread(id int64) error {
	_, err := s.conn.Exec(`DELETE FROM threads WHERE id = ?`, id)
	return err
}

func (s *HNStore) ListThreads(offset, limit int) ([]hnmentalhealth.Thread, error) {
	rows, err := s.conn.Query(
		`SELECT id, story_id, title, score, date, url FROM threads ORDER BY id LIMIT ? OFFSET ?`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var threads []hnmentalhealth.Thread
	for rows.Next() {
		var t hnmentalhealth.Thread
		var title, date, url sql.NullString
		if err := rows.Scan(&t.ID, &t.StoryID, &title, &t.Score, &date, &url); err != nil {
			return nil, err
		}
		t.Title = stringFromNull(title)
		t.Date = stringFromNull(date)
		t.URL = stringFromNull(url)
		threads = append(threads, t)
	}
	return threads, rows.Err()
}

// --- ThreadQualityRepository ---

func (s *HNStore) CreateThreadQuality(tq *hnmentalhealth.ThreadQuality) error {
	_, err := s.conn.Exec(
		`INSERT INTO thread_quality (thread_id, characterization_id, interaction_id, trauma_potential, healthy_count, notable_count) VALUES (?, ?, ?, ?, ?, ?)`,
		tq.ThreadID, nullInt64Ptr(tq.CharacterizationID), nullInt64Ptr(tq.InteractionID),
		nullString(tq.TraumaPotential), tq.HealthyCount, tq.NotableCount,
	)
	return err
}

func (s *HNStore) GetThreadQualityByThreadID(threadID int64) (*hnmentalhealth.ThreadQuality, error) {
	tq := &hnmentalhealth.ThreadQuality{}
	var charID, interID sql.NullInt64
	var trauma sql.NullString
	err := s.conn.QueryRow(
		`SELECT thread_id, characterization_id, interaction_id, trauma_potential, healthy_count, notable_count FROM thread_quality WHERE thread_id = ?`, threadID,
	).Scan(&tq.ThreadID, &charID, &interID, &trauma, &tq.HealthyCount, &tq.NotableCount)
	if err != nil {
		return nil, err
	}
	tq.CharacterizationID = int64PtrFromNull(charID)
	tq.InteractionID = int64PtrFromNull(interID)
	tq.TraumaPotential = stringFromNull(trauma)
	return tq, nil
}

func (s *HNStore) UpdateThreadQuality(tq *hnmentalhealth.ThreadQuality) error {
	res, err := s.conn.Exec(
		`UPDATE thread_quality SET characterization_id = ?, interaction_id = ?, trauma_potential = ?, healthy_count = ?, notable_count = ? WHERE thread_id = ?`,
		nullInt64Ptr(tq.CharacterizationID), nullInt64Ptr(tq.InteractionID),
		nullString(tq.TraumaPotential), tq.HealthyCount, tq.NotableCount, tq.ThreadID,
	)
	if err != nil {
		return err
	}
	// Wenn keine Zeile betroffen war, existiert der Thread-Qualitäts-Eintrag
	// noch nicht — der Aufrufer soll dann Create ausführen.
	if n, err := res.RowsAffected(); err == nil && n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *HNStore) DeleteThreadQuality(threadID int64) error {
	_, err := s.conn.Exec(`DELETE FROM thread_quality WHERE thread_id = ?`, threadID)
	return err
}

// --- ThreadEmotionRepository ---

func (s *HNStore) SetThreadEmotion(te *hnmentalhealth.ThreadEmotion) error {
	_, err := s.conn.Exec(
		`INSERT OR REPLACE INTO thread_emotions (thread_id, emotion_id, count) VALUES (?, ?, ?)`,
		te.ThreadID, te.EmotionID, te.Count,
	)
	return err
}

func (s *HNStore) GetThreadEmotions(threadID int64) ([]hnmentalhealth.ThreadEmotion, error) {
	rows, err := s.conn.Query(
		`SELECT thread_id, emotion_id, count FROM thread_emotions WHERE thread_id = ?`, threadID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var emotions []hnmentalhealth.ThreadEmotion
	for rows.Next() {
		var te hnmentalhealth.ThreadEmotion
		if err := rows.Scan(&te.ThreadID, &te.EmotionID, &te.Count); err != nil {
			return nil, err
		}
		emotions = append(emotions, te)
	}
	return emotions, rows.Err()
}

func (s *HNStore) DeleteThreadEmotion(threadID, emotionID int64) error {
	_, err := s.conn.Exec(`DELETE FROM thread_emotions WHERE thread_id = ? AND emotion_id = ?`, threadID, emotionID)
	return err
}

// --- HNFindingRepository ---

func (s *HNStore) CreateHNFinding(f *hnmentalhealth.HNFinding) error {
	r, err := s.conn.Exec(
		`INSERT INTO hn_findings (user_id, thread_id, pattern_id, confidence_id, icd11_id, evidence, comment_score, ai_hint, ai_note) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		f.UserID, f.ThreadID, nullInt64Ptr(f.PatternID), nullInt64Ptr(f.ConfidenceID),
		nullInt64Ptr(f.ICD11ID), nullString(f.Evidence), f.CommentScore,
		nullString(f.AIHint), nullString(f.AINote),
	)
	if err != nil {
		return err
	}
	f.ID, _ = r.LastInsertId()
	return nil
}

func (s *HNStore) GetHNFindingByID(id int64) (*hnmentalhealth.HNFinding, error) {
	f := &hnmentalhealth.HNFinding{}
	var patternID, confID, icd11ID sql.NullInt64
	var evidence, hint, note sql.NullString
	err := s.conn.QueryRow(
		`SELECT id, user_id, thread_id, pattern_id, confidence_id, icd11_id, evidence, comment_score, ai_hint, ai_note FROM hn_findings WHERE id = ?`, id,
	).Scan(&f.ID, &f.UserID, &f.ThreadID, &patternID, &confID, &icd11ID, &evidence, &f.CommentScore, &hint, &note)
	if err != nil {
		return nil, err
	}
	f.PatternID = int64PtrFromNull(patternID)
	f.ConfidenceID = int64PtrFromNull(confID)
	f.ICD11ID = int64PtrFromNull(icd11ID)
	f.Evidence = stringFromNull(evidence)
	f.AIHint = stringFromNull(hint)
	f.AINote = stringFromNull(note)
	return f, nil
}

func (s *HNStore) UpdateHNFinding(f *hnmentalhealth.HNFinding) error {
	_, err := s.conn.Exec(
		`UPDATE hn_findings SET user_id = ?, thread_id = ?, pattern_id = ?, confidence_id = ?, icd11_id = ?, evidence = ?, comment_score = ?, ai_hint = ?, ai_note = ? WHERE id = ?`,
		f.UserID, f.ThreadID, nullInt64Ptr(f.PatternID), nullInt64Ptr(f.ConfidenceID),
		nullInt64Ptr(f.ICD11ID), nullString(f.Evidence), f.CommentScore,
		nullString(f.AIHint), nullString(f.AINote), f.ID,
	)
	return err
}

func (s *HNStore) DeleteHNFinding(id int64) error {
	_, err := s.conn.Exec(`DELETE FROM hn_findings WHERE id = ?`, id)
	return err
}

func (s *HNStore) ListHNFindingsByUser(userID int64) ([]hnmentalhealth.HNFinding, error) {
	rows, err := s.conn.Query(
		`SELECT id, user_id, thread_id, pattern_id, confidence_id, icd11_id, evidence, comment_score, ai_hint, ai_note FROM hn_findings WHERE user_id = ? ORDER BY id`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanHNFindings(rows)
}

func (s *HNStore) ListHNFindingsByThread(threadID int64) ([]hnmentalhealth.HNFinding, error) {
	rows, err := s.conn.Query(
		`SELECT id, user_id, thread_id, pattern_id, confidence_id, icd11_id, evidence, comment_score, ai_hint, ai_note FROM hn_findings WHERE thread_id = ? ORDER BY id`, threadID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanHNFindings(rows)
}

func (s *HNStore) ListHNFindingsByPattern(patternID int64) ([]hnmentalhealth.HNFinding, error) {
	rows, err := s.conn.Query(
		`SELECT id, user_id, thread_id, pattern_id, confidence_id, icd11_id, evidence, comment_score, ai_hint, ai_note FROM hn_findings WHERE pattern_id = ? ORDER BY id`, patternID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanHNFindings(rows)
}

func (s *HNStore) CountHNFindingsByUserPattern(userID, patternID int64) (int, error) {
	var count int
	err := s.conn.QueryRow(
		`SELECT COUNT(DISTINCT thread_id) FROM hn_findings WHERE user_id = ? AND pattern_id = ?`,
		userID, patternID,
	).Scan(&count)
	return count, err
}

func scanHNFindings(rows *sql.Rows) ([]hnmentalhealth.HNFinding, error) {
	var findings []hnmentalhealth.HNFinding
	for rows.Next() {
		var f hnmentalhealth.HNFinding
		var patternID, confID, icd11ID sql.NullInt64
		var evidence, hint, note sql.NullString
		if err := rows.Scan(&f.ID, &f.UserID, &f.ThreadID, &patternID, &confID, &icd11ID, &evidence, &f.CommentScore, &hint, &note); err != nil {
			return nil, err
		}
		f.PatternID = int64PtrFromNull(patternID)
		f.ConfidenceID = int64PtrFromNull(confID)
		f.ICD11ID = int64PtrFromNull(icd11ID)
		f.Evidence = stringFromNull(evidence)
		f.AIHint = stringFromNull(hint)
		f.AINote = stringFromNull(note)
		findings = append(findings, f)
	}
	return findings, rows.Err()
}
