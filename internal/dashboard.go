// Package internal — Dashboard REST-API für hn-mentalhealth.
// Architektur folgt go-hackernews: HTTP-Handler lesen aus der DB,
// geben JSON zurück. Das Frontend (index.html) ruft diese Endpoints auf.
package internal

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/coolvegan/hn-mentalhealth/internal/sqlite"
)

// RegisterDashboardRoutes registriert die REST-API-Routen für das Dashboard.
// Wird vom main-Binary aufgerufen, bevor der HTTP-Server startet.
func RegisterDashboardRoutes(mux *http.ServeMux, db *sqlite.DB) {
	// Statisches Frontend
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "dashboard/index.html")
	})

	// --- API: Threads mit Qualitäts-Metriken ---
	mux.HandleFunc("/api/threads", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		threads, err := db.HN.ListThreads(0, 500)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		type ThreadWithQuality struct {
			ID              int64  `json:"id"`
			StoryID         int64  `json:"story_id"`
			Title           string `json:"title"`
			Score           int    `json:"score"`
			Date            string `json:"date"`
			TraumaPotential string `json:"trauma_potential"`
			HealthyCount    int    `json:"healthy_count"`
			NotableCount    int    `json:"notable_count"`
			TotalFindings   int    `json:"total_findings"`
		}
		result := make([]ThreadWithQuality, 0, len(threads))
		for _, t := range threads {
			twq := ThreadWithQuality{
				ID:      t.ID,
				StoryID: t.StoryID,
				Title:   t.Title,
				Score:   t.Score,
				Date:    t.Date,
			}
			if tq, err := db.HN.GetThreadQualityByThreadID(t.ID); err == nil {
				twq.TraumaPotential = tq.TraumaPotential
				twq.HealthyCount = tq.HealthyCount
				twq.NotableCount = tq.NotableCount
			}
			findings, _ := db.HN.ListHNFindingsByThread(t.ID)
			twq.TotalFindings = len(findings)
			result = append(result, twq)
		}
		json.NewEncoder(w).Encode(result)
	})

	// --- API: Thread-Detail (Findings + User) ---
	mux.HandleFunc("/api/thread/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", 400)
			return
		}
		thread, err := db.HN.GetThreadByID(id)
		if err != nil {
			http.Error(w, err.Error(), 404)
			return
		}
		findings, _ := db.HN.ListHNFindingsByThread(id)
		type FindingWithUser struct {
			ID           int64  `json:"id"`
			UserID       int64  `json:"user_id"`
			Username     string `json:"username"`
			PatternName  string `json:"pattern_name"`
			Confidence   string `json:"confidence"`
			ICD11Code    string `json:"icd11_code,omitempty"`
			ICD11Label   string `json:"icd11_label,omitempty"`
			Evidence     string `json:"evidence"`
			CommentScore int    `json:"comment_score"`
		}
		fwu := make([]FindingWithUser, 0, len(findings))
		for _, f := range findings {
			fw := FindingWithUser{
				ID:           f.ID,
				UserID:       f.UserID,
				Evidence:     f.Evidence,
				CommentScore: f.CommentScore,
			}
			if u, err := db.Shared.GetUserByID(f.UserID); err == nil {
				fw.Username = u.Username
			}
			if f.PatternID != nil {
				if p, err := db.HN.GetPatternByID(*f.PatternID); err == nil {
					fw.PatternName = p.Name
				}
			}
			if f.ConfidenceID != nil {
				if cl, err := db.Shared.GetConfidenceLevelByID(*f.ConfidenceID); err == nil {
					fw.Confidence = cl.Name
				}
			}
			if f.ICD11ID != nil {
				if icd, err := db.Shared.GetICD11CodeByID(*f.ICD11ID); err == nil {
					fw.ICD11Code = icd.Code
					fw.ICD11Label = icd.Label
				}
			}
			fwu = append(fwu, fw)
		}
		type ThreadDetail struct {
			Thread   interface{}       `json:"thread"`
			Quality  interface{}       `json:"quality"`
			Findings []FindingWithUser `json:"findings"`
		}
		td := ThreadDetail{Thread: thread, Findings: fwu}
		if tq, err := db.HN.GetThreadQualityByThreadID(id); err == nil {
			td.Quality = tq
		}
		json.NewEncoder(w).Encode(td)
	})

	// --- API: Toxische User (nach Anzahl auffälliger Findings) ---
	mux.HandleFunc("/api/users/toxic", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Alle HN-Findings durchgehen, nach User gruppieren, nicht-healthy zählen
		users, err := db.Shared.ListUsers(0, 1000)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		type ToxicUser struct {
			UserID        int64            `json:"user_id"`
			Username      string           `json:"username"`
			TotalFindings int              `json:"total_findings"`
			NotableCount  int              `json:"notable_count"`
			Patterns      map[string]int   `json:"patterns"`
			ICD11Refs     []string         `json:"icd11_refs"`
		}
		result := make([]ToxicUser, 0)
		for _, u := range users {
			findings, _ := db.HN.ListHNFindingsByUser(u.ID)
			if len(findings) == 0 {
				continue
			}
			tu := ToxicUser{
				UserID:        u.ID,
				Username:      u.Username,
				TotalFindings: len(findings),
				Patterns:      make(map[string]int),
			}
			for _, f := range findings {
				if f.ConfidenceID != nil {
					if cl, err := db.Shared.GetConfidenceLevelByID(*f.ConfidenceID); err == nil && cl.Name != "healthy" {
						tu.NotableCount++
					}
				}
				if f.PatternID != nil {
					if p, err := db.HN.GetPatternByID(*f.PatternID); err == nil {
						tu.Patterns[p.Name]++
					}
				}
				if f.ICD11ID != nil {
					if icd, err := db.Shared.GetICD11CodeByID(*f.ICD11ID); err == nil {
						tu.ICD11Refs = append(tu.ICD11Refs, icd.Code)
					}
				}
			}
			if tu.NotableCount > 0 {
				result = append(result, tu)
			}
		}
		json.NewEncoder(w).Encode(result)
	})

	// --- API: Stats (Gesamtverteilung) ---
	mux.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		threads, _ := db.HN.ListThreads(0, 500)
		patterns, _ := db.HN.ListPatterns()
		confidenceLevels, _ := db.Shared.ListConfidenceLevels()

		healthyThreads := 0
		notableThreads := 0
		for _, t := range threads {
			if tq, err := db.HN.GetThreadQualityByThreadID(t.ID); err == nil {
				if tq.NotableCount > 0 {
					notableThreads++
				} else {
					healthyThreads++
				}
			}
		}

		// Pattern-Verteilung über alle Findings
		patternDist := make(map[string]int)
		confDist := make(map[string]int)
		for _, t := range threads {
			findings, _ := db.HN.ListHNFindingsByThread(t.ID)
			for _, f := range findings {
				if f.PatternID != nil {
					for _, p := range patterns {
						if p.ID == *f.PatternID {
							patternDist[p.Name]++
							break
						}
					}
				}
				if f.ConfidenceID != nil {
					for _, cl := range confidenceLevels {
						if cl.ID == *f.ConfidenceID {
							confDist[cl.Name]++
							break
						}
					}
				}
			}
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"total_threads":     len(threads),
			"healthy_threads":   healthyThreads,
			"notable_threads":   notableThreads,
			"pattern_distribution": patternDist,
			"confidence_distribution": confDist,
		})
	})
}
