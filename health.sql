-- hn-mentalhealth — Health-Schema (Gesundheit, Gene, Biomarker)
-- Nutzt shared.sql für users, confidence_levels, icd11_codes, ai_knowledge.

-- ============================================================
-- GEN-SNIPPETS & GEN-MUSTER (zweistufig, dynamisch)
-- ============================================================

-- Stufe 1: Einzelne Gen-Snippets eines Patienten.
CREATE TABLE gene_snippets (
    id         INTEGER PRIMARY KEY,
    patient_id INTEGER NOT NULL REFERENCES users(id),
    gene       TEXT NOT NULL,           -- z.B. 'MTHFR', 'COMT', 'MTRR'
    variant    TEXT,                    -- z.B. 'C677T', 'A1298C'
    genotype   TEXT,                    -- z.B. 'C/T', 'T/T'
    ai_hint    TEXT,                    -- Studien-Referenz: was bewirkt dieses Gen?
    ai_note    TEXT                     -- KI-Notiz
);

-- Stufe 2: Gen-Muster — Aggregation mehrerer Gen-Snippets.
CREATE TABLE gene_patterns (
    id          INTEGER PRIMARY KEY,
    name        TEXT NOT NULL UNIQUE,   -- z.B. 'Methylierungs-Störung'
    description TEXT,
    ai_hint     TEXT,                   -- Studien zum Zusammenspiel
    ai_note     TEXT
);

-- Viele-zu-viele: Snippets ↔ Muster
CREATE TABLE pattern_snippets (
    pattern_id INTEGER REFERENCES gene_patterns(id),
    snippet_id INTEGER REFERENCES gene_snippets(id),
    PRIMARY KEY (pattern_id, snippet_id)
);

-- Viele-zu-viele: Patienten ↔ Muster (mit eigener KI-Notiz)
CREATE TABLE patient_patterns (
    patient_id INTEGER REFERENCES users(id),
    pattern_id INTEGER REFERENCES gene_patterns(id),
    ai_note    TEXT,
    PRIMARY KEY (patient_id, pattern_id)
);

-- ============================================================
-- BIOMARKER & GESUNDHEITSDATEN
-- ============================================================

-- Blutwerte / Laborwerte eines Patienten (erweiterbar)
CREATE TABLE biomarkers (
    id              INTEGER PRIMARY KEY,
    patient_id      INTEGER NOT NULL REFERENCES users(id),
    name            TEXT NOT NULL,       -- z.B. 'Folsäure', 'B12', 'Homocystein'
    value           REAL,                -- Messwert
    unit            TEXT,                -- z.B. 'ng/ml', 'pg/ml'
    reference_low   REAL,                -- unterer Referenzwert
    reference_high  REAL,                -- oberer Referenzwert
    measured_at     TEXT,                -- ISO-Datum
    ai_hint         TEXT,
    ai_note         TEXT
);

-- Diagnosen / Gesundheitszustände (erweiterbar)
CREATE TABLE health_conditions (
    id         INTEGER PRIMARY KEY,
    patient_id INTEGER NOT NULL REFERENCES users(id),
    condition  TEXT NOT NULL,            -- z.B. 'Diabetes Typ 2', 'KPU/HPU'
    diagnosed_at TEXT,
    ai_hint    TEXT,
    ai_note    TEXT
);

-- Lebenskontext / Sorgen (erweiterbar)
CREATE TABLE life_context (
    id         INTEGER PRIMARY KEY,
    patient_id INTEGER NOT NULL REFERENCES users(id),
    context    TEXT NOT NULL,            -- z.B. 'Arbeitslosigkeit', 'Trennung', 'finanzielle Sorgen'
    noted_at   TEXT,
    ai_hint    TEXT,
    ai_note    TEXT
);

-- ============================================================
-- HEALTH-FINDINGS (Domänen-spezifisch)
-- ============================================================

-- Health-Findings: eine Auffälligkeit aus Gesundheitsdaten.
-- Referenziert shared.users, shared.confidence_levels, shared.icd11_codes.
-- Etiology: metabolic | psychological | mixed | unknown
-- Evidence-Chain: JSON-Array der Quellen (z.B. ["MTHFR C677T", "Folsäure niedrig", "Studie X"])
CREATE TABLE health_findings (
    id             INTEGER PRIMARY KEY,
    patient_id     INTEGER NOT NULL REFERENCES users(id),
    pattern_id     INTEGER REFERENCES gene_patterns(id),  -- optional
    confidence_id  INTEGER REFERENCES confidence_levels(id),
    icd11_id       INTEGER REFERENCES icd11_codes(id),
    etiology       TEXT,                 -- 'metabolic' | 'psychological' | 'mixed' | 'unknown'
    evidence_chain TEXT,                 -- JSON-Array der Evidenz-Quellen
    evidence       TEXT,                 -- Freitext-Beschreibung
    ai_hint        TEXT,
    ai_note        TEXT
);

-- ============================================================
-- INDIZES
-- ============================================================
CREATE INDEX idx_health_findings_patient  ON health_findings(patient_id);
CREATE INDEX idx_health_findings_pattern  ON health_findings(pattern_id);
CREATE INDEX idx_gene_snippets_patient    ON gene_snippets(patient_id);
CREATE INDEX idx_biomarkers_patient       ON biomarkers(patient_id);
CREATE INDEX idx_health_conditions_patient ON health_conditions(patient_id);
CREATE INDEX idx_life_context_patient     ON life_context(patient_id);
