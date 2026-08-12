-- hn-mentalhealth — HN-Schema (HN-spezifische Tabellen)
-- Nutzt shared.sql für users, confidence_levels, icd11_codes, ai_knowledge.

-- ============================================================
-- HN-LOOKUP-TABELLEN (erweiterbar)
-- ============================================================

CREATE TABLE patterns (
    id          INTEGER PRIMARY KEY,
    name        TEXT NOT NULL UNIQUE,   -- z.B. 'workism', 'scham', 'flucht'
    description TEXT,
    ai_hint     TEXT,
    ai_note     TEXT
);

CREATE TABLE emotions (
    id      INTEGER PRIMARY KEY,
    name    TEXT NOT NULL UNIQUE,       -- z.B. 'joy', 'anger', 'sad', 'fear', 'shame'
    ai_hint TEXT,
    ai_note TEXT
);

CREATE TABLE characterizations (
    id      INTEGER PRIMARY KEY,
    name    TEXT NOT NULL UNIQUE,       -- z.B. 'values_debate', 'grief_group', 'nostalgic_tech'
    ai_hint TEXT,
    ai_note TEXT
);

CREATE TABLE interactions (
    id      INTEGER PRIMARY KEY,
    name    TEXT NOT NULL UNIQUE,       -- z.B. 'productive', 'escalation', 'polarized', 'flat'
    ai_hint TEXT,
    ai_note TEXT
);

-- ============================================================
-- HN-KERN-ENTITÄTEN
-- ============================================================

CREATE TABLE threads (
    id       INTEGER PRIMARY KEY,
    story_id INTEGER NOT NULL UNIQUE,
    title    TEXT,
    score    INTEGER,
    date     TEXT,
    url      TEXT
);

CREATE TABLE thread_quality (
    thread_id           INTEGER PRIMARY KEY REFERENCES threads(id),
    characterization_id INTEGER REFERENCES characterizations(id),
    interaction_id      INTEGER REFERENCES interactions(id),
    trauma_potential    TEXT,           -- 'low' | 'medium' | 'high'
    healthy_count       INTEGER DEFAULT 0,
    notable_count       INTEGER DEFAULT 0
);

CREATE TABLE thread_emotions (
    thread_id  INTEGER REFERENCES threads(id),
    emotion_id INTEGER REFERENCES emotions(id),
    count      INTEGER DEFAULT 0,
    PRIMARY KEY (thread_id, emotion_id)
);

-- HN-Findings: eine Auffälligkeit eines Users in einem HN-Thread.
-- Referenziert shared.users, shared.confidence_levels, shared.icd11_codes.
CREATE TABLE hn_findings (
    id            INTEGER PRIMARY KEY,
    user_id       INTEGER NOT NULL REFERENCES users(id),
    thread_id     INTEGER NOT NULL REFERENCES threads(id),
    pattern_id    INTEGER REFERENCES patterns(id),
    confidence_id INTEGER REFERENCES confidence_levels(id),
    icd11_id      INTEGER REFERENCES icd11_codes(id),
    evidence      TEXT,
    comment_score INTEGER DEFAULT 0,
    ai_hint       TEXT,
    ai_note       TEXT
);

-- ============================================================
-- INDIZES
-- ============================================================
CREATE INDEX idx_hn_findings_user    ON hn_findings(user_id);
CREATE INDEX idx_hn_findings_thread  ON hn_findings(thread_id);
CREATE INDEX idx_hn_findings_pattern ON hn_findings(pattern_id);
CREATE INDEX idx_thread_emotions_tid ON thread_emotions(thread_id);
