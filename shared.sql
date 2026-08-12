-- hn-mentalhealth — Shared Schema (Domänen-übergreifend)
-- Wird von HN- und Health-Schema gemeinsam genutzt.
-- users, icd11_codes, confidence_levels, ai_knowledge

PRAGMA foreign_keys = ON;

-- ============================================================
-- USERS (Personen — HN-Kommentatoren UND Patienten)
-- ============================================================
CREATE TABLE users (
    id            INTEGER PRIMARY KEY,
    username      TEXT NOT NULL UNIQUE,
    first_seen_at TEXT,
    thread_count  INTEGER DEFAULT 0
);

-- ============================================================
-- KONFIDENZ-STUFEN (erweiterbar)
-- ============================================================
CREATE TABLE confidence_levels (
    id      INTEGER PRIMARY KEY,
    name    TEXT NOT NULL UNIQUE,   -- z.B. 'secure', 'suspected', 'observation', 'healthy'
    ai_hint TEXT,
    ai_note TEXT
);

-- ============================================================
-- ICD-11-REFERENZEN (erweiterbar)
-- ============================================================
CREATE TABLE icd11_codes (
    id       INTEGER PRIMARY KEY,
    code     TEXT NOT NULL UNIQUE,       -- z.B. '6B40', 'QE84', 'non-icd:moral_injury'
    label    TEXT NOT NULL,
    is_icd11 INTEGER NOT NULL DEFAULT 1, -- 1 = echte ICD-11-Kategorie, 0 = nicht-ICD-Konstrukt
    note     TEXT,
    ai_hint  TEXT,
    ai_note  TEXT
);

-- ============================================================
-- AI-WISSENS-REPOSITORY (selbst-beschreibendes Daten-Wörterbuch)
-- ============================================================
CREATE TABLE ai_knowledge (
    id          INTEGER PRIMARY KEY,
    table_name  TEXT NOT NULL,
    column_name TEXT,
    value       TEXT,
    explanation TEXT,
    ai_note     TEXT,
    UNIQUE (table_name, column_name, value)
);

CREATE INDEX idx_ai_knowledge_table ON ai_knowledge(table_name);
CREATE INDEX idx_ai_knowledge_value ON ai_knowledge(value);
