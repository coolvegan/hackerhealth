// Package sqlite — Shared Store (users, confidence_levels, icd11_codes, ai_knowledge).
package sqlite

import (
	"database/sql"

	"github.com/coolvegan/hn-mentalhealth"
)

// SharedStore implementiert repository.SharedRepository.
type SharedStore struct{ conn *sql.DB }

// --- UserRepository ---

func (s *SharedStore) CreateUser(u *hnmentalhealth.User) error {
	r, err := s.conn.Exec(
		`INSERT INTO users (username, first_seen_at, thread_count) VALUES (?, ?, ?)`,
		u.Username, nullString(u.FirstSeenAt), u.ThreadCount,
	)
	if err != nil {
		return err
	}
	u.ID, _ = r.LastInsertId()
	return nil
}

func (s *SharedStore) GetUserByID(id int64) (*hnmentalhealth.User, error) {
	u := &hnmentalhealth.User{}
	var firstSeen sql.NullString
	err := s.conn.QueryRow(
		`SELECT id, username, first_seen_at, thread_count FROM users WHERE id = ?`, id,
	).Scan(&u.ID, &u.Username, &firstSeen, &u.ThreadCount)
	if err != nil {
		return nil, err
	}
	u.FirstSeenAt = stringFromNull(firstSeen)
	return u, nil
}

func (s *SharedStore) GetUserByUsername(username string) (*hnmentalhealth.User, error) {
	u := &hnmentalhealth.User{}
	var firstSeen sql.NullString
	err := s.conn.QueryRow(
		`SELECT id, username, first_seen_at, thread_count FROM users WHERE username = ?`, username,
	).Scan(&u.ID, &u.Username, &firstSeen, &u.ThreadCount)
	if err != nil {
		return nil, err
	}
	u.FirstSeenAt = stringFromNull(firstSeen)
	return u, nil
}

func (s *SharedStore) UpdateUser(u *hnmentalhealth.User) error {
	_, err := s.conn.Exec(
		`UPDATE users SET username = ?, first_seen_at = ?, thread_count = ? WHERE id = ?`,
		u.Username, nullString(u.FirstSeenAt), u.ThreadCount, u.ID,
	)
	return err
}

func (s *SharedStore) DeleteUser(id int64) error {
	_, err := s.conn.Exec(`DELETE FROM users WHERE id = ?`, id)
	return err
}

func (s *SharedStore) ListUsers(offset, limit int) ([]hnmentalhealth.User, error) {
	rows, err := s.conn.Query(
		`SELECT id, username, first_seen_at, thread_count FROM users ORDER BY id LIMIT ? OFFSET ?`,
		limit, offset,
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

// --- ConfidenceLevelRepository ---

func (s *SharedStore) CreateConfidenceLevel(cl *hnmentalhealth.ConfidenceLevel) error {
	r, err := s.conn.Exec(
		`INSERT INTO confidence_levels (name, ai_hint, ai_note) VALUES (?, ?, ?)`,
		cl.Name, nullString(cl.AIHint), nullString(cl.AINote),
	)
	if err != nil {
		return err
	}
	cl.ID, _ = r.LastInsertId()
	return nil
}

func (s *SharedStore) GetConfidenceLevelByID(id int64) (*hnmentalhealth.ConfidenceLevel, error) {
	cl := &hnmentalhealth.ConfidenceLevel{}
	var hint, note sql.NullString
	err := s.conn.QueryRow(
		`SELECT id, name, ai_hint, ai_note FROM confidence_levels WHERE id = ?`, id,
	).Scan(&cl.ID, &cl.Name, &hint, &note)
	if err != nil {
		return nil, err
	}
	cl.AIHint = stringFromNull(hint)
	cl.AINote = stringFromNull(note)
	return cl, nil
}

func (s *SharedStore) GetConfidenceLevelByName(name string) (*hnmentalhealth.ConfidenceLevel, error) {
	cl := &hnmentalhealth.ConfidenceLevel{}
	var hint, note sql.NullString
	err := s.conn.QueryRow(
		`SELECT id, name, ai_hint, ai_note FROM confidence_levels WHERE name = ?`, name,
	).Scan(&cl.ID, &cl.Name, &hint, &note)
	if err != nil {
		return nil, err
	}
	cl.AIHint = stringFromNull(hint)
	cl.AINote = stringFromNull(note)
	return cl, nil
}

func (s *SharedStore) UpdateConfidenceLevel(cl *hnmentalhealth.ConfidenceLevel) error {
	_, err := s.conn.Exec(
		`UPDATE confidence_levels SET name = ?, ai_hint = ?, ai_note = ? WHERE id = ?`,
		cl.Name, nullString(cl.AIHint), nullString(cl.AINote), cl.ID,
	)
	return err
}

func (s *SharedStore) DeleteConfidenceLevel(id int64) error {
	_, err := s.conn.Exec(`DELETE FROM confidence_levels WHERE id = ?`, id)
	return err
}

func (s *SharedStore) ListConfidenceLevels() ([]hnmentalhealth.ConfidenceLevel, error) {
	rows, err := s.conn.Query(`SELECT id, name, ai_hint, ai_note FROM confidence_levels ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var levels []hnmentalhealth.ConfidenceLevel
	for rows.Next() {
		var cl hnmentalhealth.ConfidenceLevel
		var hint, note sql.NullString
		if err := rows.Scan(&cl.ID, &cl.Name, &hint, &note); err != nil {
			return nil, err
		}
		cl.AIHint = stringFromNull(hint)
		cl.AINote = stringFromNull(note)
		levels = append(levels, cl)
	}
	return levels, rows.Err()
}

// --- ICD11CodeRepository ---

func (s *SharedStore) CreateICD11Code(c *hnmentalhealth.ICD11Code) error {
	r, err := s.conn.Exec(
		`INSERT INTO icd11_codes (code, label, is_icd11, note, ai_hint, ai_note) VALUES (?, ?, ?, ?, ?, ?)`,
		c.Code, c.Label, c.IsICD11, nullString(c.Note), nullString(c.AIHint), nullString(c.AINote),
	)
	if err != nil {
		return err
	}
	c.ID, _ = r.LastInsertId()
	return nil
}

func (s *SharedStore) GetICD11CodeByID(id int64) (*hnmentalhealth.ICD11Code, error) {
	c := &hnmentalhealth.ICD11Code{}
	var note, hint, aiNote sql.NullString
	err := s.conn.QueryRow(
		`SELECT id, code, label, is_icd11, note, ai_hint, ai_note FROM icd11_codes WHERE id = ?`, id,
	).Scan(&c.ID, &c.Code, &c.Label, &c.IsICD11, &note, &hint, &aiNote)
	if err != nil {
		return nil, err
	}
	c.Note = stringFromNull(note)
	c.AIHint = stringFromNull(hint)
	c.AINote = stringFromNull(aiNote)
	return c, nil
}

func (s *SharedStore) GetICD11CodeByCode(code string) (*hnmentalhealth.ICD11Code, error) {
	c := &hnmentalhealth.ICD11Code{}
	var note, hint, aiNote sql.NullString
	err := s.conn.QueryRow(
		`SELECT id, code, label, is_icd11, note, ai_hint, ai_note FROM icd11_codes WHERE code = ?`, code,
	).Scan(&c.ID, &c.Code, &c.Label, &c.IsICD11, &note, &hint, &aiNote)
	if err != nil {
		return nil, err
	}
	c.Note = stringFromNull(note)
	c.AIHint = stringFromNull(hint)
	c.AINote = stringFromNull(aiNote)
	return c, nil
}

func (s *SharedStore) UpdateICD11Code(c *hnmentalhealth.ICD11Code) error {
	_, err := s.conn.Exec(
		`UPDATE icd11_codes SET code = ?, label = ?, is_icd11 = ?, note = ?, ai_hint = ?, ai_note = ? WHERE id = ?`,
		c.Code, c.Label, c.IsICD11, nullString(c.Note), nullString(c.AIHint), nullString(c.AINote), c.ID,
	)
	return err
}

func (s *SharedStore) DeleteICD11Code(id int64) error {
	_, err := s.conn.Exec(`DELETE FROM icd11_codes WHERE id = ?`, id)
	return err
}

func (s *SharedStore) ListICD11Codes() ([]hnmentalhealth.ICD11Code, error) {
	rows, err := s.conn.Query(`SELECT id, code, label, is_icd11, note, ai_hint, ai_note FROM icd11_codes ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var codes []hnmentalhealth.ICD11Code
	for rows.Next() {
		var c hnmentalhealth.ICD11Code
		var note, hint, aiNote sql.NullString
		if err := rows.Scan(&c.ID, &c.Code, &c.Label, &c.IsICD11, &note, &hint, &aiNote); err != nil {
			return nil, err
		}
		c.Note = stringFromNull(note)
		c.AIHint = stringFromNull(hint)
		c.AINote = stringFromNull(aiNote)
		codes = append(codes, c)
	}
	return codes, rows.Err()
}

// --- AIKnowledgeRepository ---

func (s *SharedStore) CreateAIKnowledge(e *hnmentalhealth.AIKnowledge) error {
	r, err := s.conn.Exec(
		`INSERT INTO ai_knowledge (table_name, column_name, value, explanation, ai_note) VALUES (?, ?, ?, ?, ?)`,
		e.TableName, nullString(e.ColumnName), nullString(e.Value), nullString(e.Explanation), nullString(e.AINote),
	)
	if err != nil {
		return err
	}
	e.ID, _ = r.LastInsertId()
	return nil
}

func (s *SharedStore) GetAIKnowledgeByID(id int64) (*hnmentalhealth.AIKnowledge, error) {
	e := &hnmentalhealth.AIKnowledge{}
	var col, val, expl, note sql.NullString
	err := s.conn.QueryRow(
		`SELECT id, table_name, column_name, value, explanation, ai_note FROM ai_knowledge WHERE id = ?`, id,
	).Scan(&e.ID, &e.TableName, &col, &val, &expl, &note)
	if err != nil {
		return nil, err
	}
	e.ColumnName = stringFromNull(col)
	e.Value = stringFromNull(val)
	e.Explanation = stringFromNull(expl)
	e.AINote = stringFromNull(note)
	return e, nil
}

func (s *SharedStore) GetAIKnowledgeByKey(tableName, columnName, value string) (*hnmentalhealth.AIKnowledge, error) {
	e := &hnmentalhealth.AIKnowledge{}
	var col, val, expl, note sql.NullString
	err := s.conn.QueryRow(
		`SELECT id, table_name, column_name, value, explanation, ai_note FROM ai_knowledge WHERE table_name = ? AND column_name IS ? AND value IS ?`,
		tableName, nullString(columnName), nullString(value),
	).Scan(&e.ID, &e.TableName, &col, &val, &expl, &note)
	if err != nil {
		return nil, err
	}
	e.ColumnName = stringFromNull(col)
	e.Value = stringFromNull(val)
	e.Explanation = stringFromNull(expl)
	e.AINote = stringFromNull(note)
	return e, nil
}

func (s *SharedStore) UpdateAIKnowledge(e *hnmentalhealth.AIKnowledge) error {
	_, err := s.conn.Exec(
		`UPDATE ai_knowledge SET table_name = ?, column_name = ?, value = ?, explanation = ?, ai_note = ? WHERE id = ?`,
		e.TableName, nullString(e.ColumnName), nullString(e.Value), nullString(e.Explanation), nullString(e.AINote), e.ID,
	)
	return err
}

func (s *SharedStore) DeleteAIKnowledge(id int64) error {
	_, err := s.conn.Exec(`DELETE FROM ai_knowledge WHERE id = ?`, id)
	return err
}

func (s *SharedStore) ListAIKnowledgeByTable(tableName string) ([]hnmentalhealth.AIKnowledge, error) {
	rows, err := s.conn.Query(
		`SELECT id, table_name, column_name, value, explanation, ai_note FROM ai_knowledge WHERE table_name = ? ORDER BY id`,
		tableName,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAIKnowledge(rows)
}

func (s *SharedStore) SearchAIKnowledge(query string, limit int) ([]hnmentalhealth.AIKnowledge, error) {
	rows, err := s.conn.Query(
		`SELECT id, table_name, column_name, value, explanation, ai_note FROM ai_knowledge
		 WHERE table_name LIKE ? OR column_name LIKE ? OR value LIKE ? OR explanation LIKE ?
		 ORDER BY id LIMIT ?`,
		"%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%", limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAIKnowledge(rows)
}

func scanAIKnowledge(rows *sql.Rows) ([]hnmentalhealth.AIKnowledge, error) {
	var entries []hnmentalhealth.AIKnowledge
	for rows.Next() {
		var e hnmentalhealth.AIKnowledge
		var col, val, expl, note sql.NullString
		if err := rows.Scan(&e.ID, &e.TableName, &col, &val, &expl, &note); err != nil {
			return nil, err
		}
		e.ColumnName = stringFromNull(col)
		e.Value = stringFromNull(val)
		e.Explanation = stringFromNull(expl)
		e.AINote = stringFromNull(note)
		entries = append(entries, e)
	}
	return entries, rows.Err()
}
