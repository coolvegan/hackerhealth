# hn-mentalhealth — Stand 2026-08-12

## Was gebaut wurde

Drei getrennte SQL-Schemata für ein AI-gestütztes System, das psychologische Auffälligkeiten in HN-Kommentaren erkennt und mit biologischen Gesundheitsdaten kombiniert.

### `shared.sql` — Domänen-übergreifend
- `users` — Personen (HN-Kommentatoren und Patienten)
- `confidence_levels` — Konfidenz-Stufen (secure/suspected/observation/healthy)
- `icd11_codes` — ICD-11-Referenzen + nicht-ICD-Konstrukte (Moral Injury, Bindungs-Trauma)
- `ai_knowledge` — selbst-beschreibendes Daten-Wörterbuch für die KI

### `hn.sql` — HN-spezifisch
- Lookup-Tabellen: `patterns`, `emotions`, `characterizations`, `interactions`
- Thread-Daten: `threads`, `thread_quality`, `thread_emotions`
- `hn_findings` — Auffälligkeiten in HN-Kommentaren

### `health.sql` — Gesundheit/Gene
- `gene_snippets` → `gene_patterns` (zweistufig, via `pattern_snippets`)
- `patient_patterns` — Patienten ↔ Muster mit eigener KI-Notiz
- `biomarkers`, `health_conditions`, `life_context` — Gesundheitsdaten
- `health_findings` — Auffälligkeiten aus Gesundheitsdaten, mit `etiology` und `evidence_chain`

## Warum so gebaut

**Kategorien sind Tabellenzeilen, nicht Code-Konstanten.** Neue Muster, Gefühle, Charakterisierungen, ICD-Codes werden per `INSERT` hinzugefügt — ohne Code-Änderung oder Rebuild. Das AI-Modell kann die Struktur zur Laufzeit erweitern.

**AI-bewusst:** Jede Lookup-Tabelle und jedes Finding hat `ai_hint` (Hinweis für die KI, geht via MCP in den Prompt) und `ai_note` (Notiz, die die KI selbst setzt). Das `ai_knowledge`-Wörterbuch dokumentiert die gesamte Struktur für die KI.

**Domänen getrennt:** HN und Health sind eigene Schemata mit eigenen Findings-Tabellen. `users`, `icd11_codes`, `confidence_levels` und `ai_knowledge` werden geteilt. So kann jede Domäne unabhängig wachsen, ohne die andere zu brechen.

**Health-Findings haben eine Ätiologie-Dimension:** `etiology` (metabolic/psychological/mixed/unknown) + `evidence_chain` (JSON-Array der Quellen). Das trennt Symptom-Klassifikation (ICD-11) von Ursachen-Hypothese — der entscheidende Punkt für personalisierte Einsichten.

## Design-Entscheidungen

1. **Dynamisch, nicht statisch** — Kategorien als Zeilen, nicht als Go-Konstanten
2. **Zwei Findings-Tabellen** — `hn_findings` und `health_findings` sind getrennt, teilen aber users/icd11/confidence
3. **Gen-Muster zweistufig** — Snippets (einzelne Gene) → Patterns (Aggregationen), viele-zu-viele
4. **ai_hint/ai_note auf allen Ebenen** — Lookup-Tabellen, Findings, Gen-Snippets, Gen-Muster, Patient-Muster
5. **ai_knowledge als zentrales Wörterbuch** — eine Tabelle, die alle anderen Tabellen/Spalten/Werte für die KI erklärt

## Was noch nicht gebaut ist

- MCP-Server (Go, mark3labs/mcp-go)
- Parser (Markdown-Berichte → SQL)
- Dashboard
- Prototyp-Testfall (fiktiver MTHFR-Patient mit Kommentar)
- Seed-Daten (die bestehenden Kategorien)
- Embedding-Dedup-Schicht (später)

## Nächste Schritte (wenn wieder aufgenommen)

1. Seed-Daten einfügen (Patterns, Emotions, ICD-11-Codes, ai_knowledge-Einträge)
2. MCP-Server-Skelett bauen (mark3labs/mcp-go, SSE)
3. Parser: 233 Markdown-Berichte → `hn_findings`
4. Prototyp-Testfall: fiktiver Patient + Kommentar → Erkenntnis-Loop durchspielen
