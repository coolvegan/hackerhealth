package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	hnmentalhealth "github.com/coolvegan/hn-mentalhealth"
	"github.com/coolvegan/hn-mentalhealth/internal/sqlite"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RunMentalHealthMcp startet den MCP-SSE-Server für hn-mentalhealth.
// Architektur folgt go-hackernews: shared state (db) wird per Closure
// an die Tool-Handler übergeben, SSE-Transport, Tools geben JSON zurück.
func RunMentalHealthMcp(addr string, db *sqlite.DB) {
	s := server.NewMCPServer(
		"hn-mentalhealth",
		"0.1.0",
		server.WithTaskCapabilities(true, true, true),
		server.WithMaxConcurrentTasks(10),
	)

	registerSharedTools(s, db)
	registerHNTools(s, db)
	registerHealthTools(s, db)

	sseServer := server.NewSSEServer(s, server.WithBaseURL(fmt.Sprintf("http://%s", addr)))
	log.Printf("MCP SSE-Server is running on %s\n", addr)
	if err := sseServer.Start(addr); err != nil {
		log.Fatalf("Server-Error: %v", err)
	}
}

// --- Shared Tools ---

func registerSharedTools(s *server.MCPServer, db *sqlite.DB) {
	// user_create
	s.AddTool(mcp.NewTool("user_create",
		mcp.WithDescription("Create a new user (HN commenter or patient)"),
		mcp.WithString("username", mcp.Required(), mcp.Description("Unique username")),
		mcp.WithString("first_seen_at", mcp.Description("ISO date of first sighting")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		username, _ := req.RequireString("username")
		firstSeen := req.GetString("first_seen_at", "")
		u := &hnmentalhealth.User{Username: username, FirstSeenAt: firstSeen}
		if err := db.Shared.CreateUser(u); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		b, _ := json.Marshal(u)
		return mcp.NewToolResultText(string(b)), nil
	})

	// user_get
	s.AddTool(mcp.NewTool("user_get",
		mcp.WithDescription("Get a user by ID or username"),
		mcp.WithNumber("id", mcp.Description("User ID")),
		mcp.WithString("username", mcp.Description("Username")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if id := req.GetInt("id", 0); id > 0 {
			u, err := db.Shared.GetUserByID(int64(id))
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			b, _ := json.Marshal(u)
			return mcp.NewToolResultText(string(b)), nil
		}
		if username := req.GetString("username", ""); username != "" {
			u, err := db.Shared.GetUserByUsername(username)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			b, _ := json.Marshal(u)
			return mcp.NewToolResultText(string(b)), nil
		}
		return mcp.NewToolResultError("provide id or username"), nil
	})

	// user_list
	s.AddTool(mcp.NewTool("user_list",
		mcp.WithDescription("List users with pagination"),
		mcp.WithNumber("offset", mcp.Description("Offset (default 0)")),
		mcp.WithNumber("limit", mcp.Description("Limit (default 50)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		offset := req.GetInt("offset", 0)
		limit := req.GetInt("limit", 50)
		users, err := db.Shared.ListUsers(offset, limit)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		b, _ := json.Marshal(users)
		return mcp.NewToolResultText(string(b)), nil
	})

	// confidence_level_list
	s.AddTool(mcp.NewTool("confidence_level_list",
		mcp.WithDescription("List all confidence levels (secure, suspected, observation, healthy)"),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		levels, err := db.Shared.ListConfidenceLevels()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		b, _ := json.Marshal(levels)
		return mcp.NewToolResultText(string(b)), nil
	})

	// icd11_list
	s.AddTool(mcp.NewTool("icd11_list",
		mcp.WithDescription("List all ICD-11 codes and non-ICD constructs"),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		codes, err := db.Shared.ListICD11Codes()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		b, _ := json.Marshal(codes)
		return mcp.NewToolResultText(string(b)), nil
	})

	// icd11_get
	s.AddTool(mcp.NewTool("icd11_get",
		mcp.WithDescription("Get an ICD-11 code by its code string (e.g. '6B70')"),
		mcp.WithString("code", mcp.Required(), mcp.Description("ICD-11 code, e.g. '6B40', 'QE84', 'non-icd:moral_injury'")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		code, _ := req.RequireString("code")
		c, err := db.Shared.GetICD11CodeByCode(code)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		b, _ := json.Marshal(c)
		return mcp.NewToolResultText(string(b)), nil
	})

	// ai_knowledge_search
	s.AddTool(mcp.NewTool("ai_knowledge_search",
		mcp.WithDescription("Search the AI knowledge dictionary"),
		mcp.WithString("query", mcp.Required(), mcp.Description("Search query")),
		mcp.WithNumber("limit", mcp.Description("Max results (default 20)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query, _ := req.RequireString("query")
		limit := req.GetInt("limit", 20)
		entries, err := db.Shared.SearchAIKnowledge(query, limit)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		b, _ := json.Marshal(entries)
		return mcp.NewToolResultText(string(b)), nil
	})

	// ai_knowledge_create
	s.AddTool(mcp.NewTool("ai_knowledge_create",
		mcp.WithDescription("Create a new AI knowledge entry (self-describing dictionary)"),
		mcp.WithString("table_name", mcp.Required(), mcp.Description("Table name")),
		mcp.WithString("column_name", mcp.Description("Column name (optional)")),
		mcp.WithString("value", mcp.Description("Value (optional)")),
		mcp.WithString("explanation", mcp.Description("Explanation for the AI")),
		mcp.WithString("ai_note", mcp.Description("AI's own note")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		tableName, _ := req.RequireString("table_name")
		e := &hnmentalhealth.AIKnowledge{
			TableName:   tableName,
			ColumnName:  req.GetString("column_name", ""),
			Value:       req.GetString("value", ""),
			Explanation: req.GetString("explanation", ""),
			AINote:      req.GetString("ai_note", ""),
		}
		if err := db.Shared.CreateAIKnowledge(e); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		b, _ := json.Marshal(e)
		return mcp.NewToolResultText(string(b)), nil
	})

	// confidence_level_create
	s.AddTool(mcp.NewTool("confidence_level_create",
		mcp.WithDescription("Create a new confidence level (secure, suspected, observation, healthy)"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Unique confidence level name")),
		mcp.WithString("ai_hint", mcp.Description("Hint for the AI model")),
		mcp.WithString("ai_note", mcp.Description("AI's own note")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, _ := req.RequireString("name")
		cl := &hnmentalhealth.ConfidenceLevel{
			Name:   name,
			AIHint: req.GetString("ai_hint", ""),
			AINote: req.GetString("ai_note", ""),
		}
		if err := db.Shared.CreateConfidenceLevel(cl); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		b, _ := json.Marshal(cl)
		return mcp.NewToolResultText(string(b)), nil
	})

	// icd11_create
	s.AddTool(mcp.NewTool("icd11_create",
		mcp.WithDescription("Create an ICD-11 code or a non-ICD construct (e.g. '6B70', 'QE84', 'non-icd:moral_injury')"),
		mcp.WithString("code", mcp.Required(), mcp.Description("Unique code string, e.g. '6B70', 'QE84', 'non-icd:moral_injury'")),
		mcp.WithString("label", mcp.Required(), mcp.Description("Human-readable label, e.g. 'Depressive Episode'")),
		mcp.WithBoolean("is_icd11", mcp.Description("true = real ICD-11 category, false = non-ICD construct. Default true")),
		mcp.WithString("note", mcp.Description("Optional note")),
		mcp.WithString("ai_hint", mcp.Description("Hint for the AI model")),
		mcp.WithString("ai_note", mcp.Description("AI's own note")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		code, _ := req.RequireString("code")
		label, _ := req.RequireString("label")
		c := &hnmentalhealth.ICD11Code{
			Code:    code,
			Label:   label,
			IsICD11: req.GetBool("is_icd11", true),
			Note:    req.GetString("note", ""),
			AIHint:  req.GetString("ai_hint", ""),
			AINote:  req.GetString("ai_note", ""),
		}
		if err := db.Shared.CreateICD11Code(c); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		b, _ := json.Marshal(c)
		return mcp.NewToolResultText(string(b)), nil
	})
}

// optIntPtr liefert einen *int64 nur dann, wenn der Schlüssel im Request
// vorhanden und > 0 ist. Fehlt der Schlüssel oder ist er 0, wird nil
// zurückgegeben, sodass die DB NULL statt 0 speichert (verhindert
// FOREIGN KEY constraint failed bei optionalen FK-Feldern).
func optIntPtr(req mcp.CallToolRequest, key string) *int64 {
	if v := req.GetInt(key, 0); v > 0 {
		id := int64(v)
		return &id
	}
	return nil
}

// --- HN Tools ---

func registerHNTools(s *server.MCPServer, db *sqlite.DB) {
	// pattern_list
	s.AddTool(mcp.NewTool("pattern_list",
		mcp.WithDescription("List all psychological patterns (workism, burnout, scham, ...)"),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		patterns, err := db.HN.ListPatterns()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		b, _ := json.Marshal(patterns)
		return mcp.NewToolResultText(string(b)), nil
	})

	// pattern_get
	s.AddTool(mcp.NewTool("pattern_get",
		mcp.WithDescription("Get a pattern by name"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Pattern name, e.g. 'workism', 'scham'")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, _ := req.RequireString("name")
		p, err := db.HN.GetPatternByName(name)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		b, _ := json.Marshal(p)
		return mcp.NewToolResultText(string(b)), nil
	})

	// pattern_create
	s.AddTool(mcp.NewTool("pattern_create",
		mcp.WithDescription("Create a new psychological pattern. Use when the AI discovers a pattern not yet in the database."),
		mcp.WithString("name", mcp.Required(), mcp.Description("Unique pattern name")),
		mcp.WithString("description", mcp.Description("What this pattern means")),
		mcp.WithString("ai_hint", mcp.Description("Hint for the AI model (goes into the prompt)")),
		mcp.WithString("ai_note", mcp.Description("AI's own note about this pattern")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, _ := req.RequireString("name")
		p := &hnmentalhealth.Pattern{
			Name:        name,
			Description: req.GetString("description", ""),
			AIHint:      req.GetString("ai_hint", ""),
			AINote:      req.GetString("ai_note", ""),
		}
		if err := db.HN.CreatePattern(p); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		b, _ := json.Marshal(p)
		return mcp.NewToolResultText(string(b)), nil
	})

	// emotion_list
	s.AddTool(mcp.NewTool("emotion_list",
		mcp.WithDescription("List all emotions (joy, anger, sad, fear, shame, ...)"),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		emotions, err := db.HN.ListEmotions()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		b, _ := json.Marshal(emotions)
		return mcp.NewToolResultText(string(b)), nil
	})

	// emotion_create
	s.AddTool(mcp.NewTool("emotion_create",
		mcp.WithDescription("Create a new emotion. Use when the AI discovers an emotion not yet in the database."),
		mcp.WithString("name", mcp.Required(), mcp.Description("Unique emotion name")),
		mcp.WithString("ai_hint", mcp.Description("Hint for the AI model")),
		mcp.WithString("ai_note", mcp.Description("AI's own note")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, _ := req.RequireString("name")
		e := &hnmentalhealth.Emotion{
			Name:   name,
			AIHint: req.GetString("ai_hint", ""),
			AINote: req.GetString("ai_note", ""),
		}
		if err := db.HN.CreateEmotion(e); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		b, _ := json.Marshal(e)
		return mcp.NewToolResultText(string(b)), nil
	})

	// characterization_create
	s.AddTool(mcp.NewTool("characterization_create",
		mcp.WithDescription("Create a new thread characterization (values_debate, grief_group, productive_debate, ...)"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Unique characterization name")),
		mcp.WithString("ai_hint", mcp.Description("Hint for the AI model")),
		mcp.WithString("ai_note", mcp.Description("AI's own note")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, _ := req.RequireString("name")
		c := &hnmentalhealth.Characterization{
			Name:   name,
			AIHint: req.GetString("ai_hint", ""),
			AINote: req.GetString("ai_note", ""),
		}
		if err := db.HN.CreateCharacterization(c); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		b, _ := json.Marshal(c)
		return mcp.NewToolResultText(string(b)), nil
	})

	// interaction_create
	s.AddTool(mcp.NewTool("interaction_create",
		mcp.WithDescription("Create a new interaction type (productive, escalation, polarized, flat)"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Unique interaction name")),
		mcp.WithString("ai_hint", mcp.Description("Hint for the AI model")),
		mcp.WithString("ai_note", mcp.Description("AI's own note")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, _ := req.RequireString("name")
		i := &hnmentalhealth.Interaction{
			Name:   name,
			AIHint: req.GetString("ai_hint", ""),
			AINote: req.GetString("ai_note", ""),
		}
		if err := db.HN.CreateInteraction(i); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		b, _ := json.Marshal(i)
		return mcp.NewToolResultText(string(b)), nil
	})

	// characterization_list
	s.AddTool(mcp.NewTool("characterization_list",
		mcp.WithDescription("List all thread characterizations (values_debate, grief_group, ...)"),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		chars, err := db.HN.ListCharacterizations()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		b, _ := json.Marshal(chars)
		return mcp.NewToolResultText(string(b)), nil
	})

	// interaction_list
	s.AddTool(mcp.NewTool("interaction_list",
		mcp.WithDescription("List all interaction types (productive, escalation, polarized, flat)"),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		interactions, err := db.HN.ListInteractions()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		b, _ := json.Marshal(interactions)
		return mcp.NewToolResultText(string(b)), nil
	})

	// thread_create
	s.AddTool(mcp.NewTool("thread_create",
		mcp.WithDescription("Create a new HN thread (story)"),
		mcp.WithNumber("story_id", mcp.Required(), mcp.Description("HN story ID")),
		mcp.WithString("title", mcp.Description("Story title")),
		mcp.WithNumber("score", mcp.Description("Story score")),
		mcp.WithString("date", mcp.Description("ISO date")),
		mcp.WithString("url", mcp.Description("Story URL")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		storyID := req.GetInt("story_id", 0)
		t := &hnmentalhealth.Thread{
			StoryID: int64(storyID),
			Title:   req.GetString("title", ""),
			Score:   req.GetInt("score", 0),
			Date:    req.GetString("date", ""),
			URL:     req.GetString("url", ""),
		}
		if err := db.HN.CreateThread(t); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		b, _ := json.Marshal(t)
		return mcp.NewToolResultText(string(b)), nil
	})

	// thread_get
	s.AddTool(mcp.NewTool("thread_get",
		mcp.WithDescription("Get a thread by ID or story_id"),
		mcp.WithNumber("id", mcp.Description("Internal thread ID")),
		mcp.WithNumber("story_id", mcp.Description("HN story ID")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if id := req.GetInt("id", 0); id > 0 {
			t, err := db.HN.GetThreadByID(int64(id))
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			b, _ := json.Marshal(t)
			return mcp.NewToolResultText(string(b)), nil
		}
		if storyID := req.GetInt("story_id", 0); storyID > 0 {
			t, err := db.HN.GetThreadByStoryID(int64(storyID))
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			b, _ := json.Marshal(t)
			return mcp.NewToolResultText(string(b)), nil
		}
		return mcp.NewToolResultError("provide id or story_id"), nil
	})

	// thread_quality_set
	s.AddTool(mcp.NewTool("thread_quality_set",
		mcp.WithDescription("Set the qualitative assessment of a thread"),
		mcp.WithNumber("thread_id", mcp.Required(), mcp.Description("Thread ID")),
		mcp.WithNumber("characterization_id", mcp.Description("Characterization ID")),
		mcp.WithNumber("interaction_id", mcp.Description("Interaction ID")),
		mcp.WithString("trauma_potential", mcp.Description("low, medium, or high")),
		mcp.WithNumber("healthy_count", mcp.Description("Number of healthy findings")),
		mcp.WithNumber("notable_count", mcp.Description("Number of notable findings")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		threadID := req.GetInt("thread_id", 0)
		charID := int64(req.GetInt("characterization_id", 0))
		interID := int64(req.GetInt("interaction_id", 0))
		tq := &hnmentalhealth.ThreadQuality{
			ThreadID:           int64(threadID),
			CharacterizationID: &charID,
			InteractionID:      &interID,
			TraumaPotential:    req.GetString("trauma_potential", ""),
			HealthyCount:       req.GetInt("healthy_count", 0),
			NotableCount:       req.GetInt("notable_count", 0),
		}
		// Try update first, create if not exists
		if err := db.HN.UpdateThreadQuality(tq); err != nil {
			if err := db.HN.CreateThreadQuality(tq); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
		}
		b, _ := json.Marshal(tq)
		return mcp.NewToolResultText(string(b)), nil
	})

	// thread_emotion_set
	s.AddTool(mcp.NewTool("thread_emotion_set",
		mcp.WithDescription("Set an emotion count for a thread"),
		mcp.WithNumber("thread_id", mcp.Required(), mcp.Description("Thread ID")),
		mcp.WithNumber("emotion_id", mcp.Required(), mcp.Description("Emotion ID")),
		mcp.WithNumber("count", mcp.Required(), mcp.Description("Emotion count")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		te := &hnmentalhealth.ThreadEmotion{
			ThreadID:  int64(req.GetInt("thread_id", 0)),
			EmotionID: int64(req.GetInt("emotion_id", 0)),
			Count:     req.GetInt("count", 0),
		}
		if err := db.HN.SetThreadEmotion(te); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		b, _ := json.Marshal(te)
		return mcp.NewToolResultText(string(b)), nil
	})

	// hn_finding_create
	s.AddTool(mcp.NewTool("hn_finding_create",
		mcp.WithDescription("Create a new HN finding — a psychological observation about a user in a thread"),
		mcp.WithNumber("user_id", mcp.Required(), mcp.Description("User ID")),
		mcp.WithNumber("thread_id", mcp.Required(), mcp.Description("Thread ID")),
		mcp.WithNumber("pattern_id", mcp.Description("Pattern ID (from pattern_list/pattern_get)")),
		mcp.WithNumber("confidence_id", mcp.Description("Confidence level ID (from confidence_level_list)")),
		mcp.WithNumber("icd11_id", mcp.Description("ICD-11 code ID (from icd11_list)")),
		mcp.WithString("evidence", mcp.Description("Comment excerpt as evidence")),
		mcp.WithNumber("comment_score", mcp.Description("Score of the comment")),
		mcp.WithString("ai_hint", mcp.Description("AI hint for context")),
		mcp.WithString("ai_note", mcp.Description("AI's own note")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		f := &hnmentalhealth.HNFinding{
			UserID:       int64(req.GetInt("user_id", 0)),
			ThreadID:     int64(req.GetInt("thread_id", 0)),
			PatternID:    optIntPtr(req, "pattern_id"),
			ConfidenceID: optIntPtr(req, "confidence_id"),
			ICD11ID:      optIntPtr(req, "icd11_id"),
			Evidence:     req.GetString("evidence", ""),
			CommentScore: req.GetInt("comment_score", 0),
			AIHint:       req.GetString("ai_hint", ""),
			AINote:       req.GetString("ai_note", ""),
		}
		if err := db.HN.CreateHNFinding(f); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		b, _ := json.Marshal(f)
		return mcp.NewToolResultText(string(b)), nil
	})

	// hn_finding_list_by_user
	s.AddTool(mcp.NewTool("hn_finding_list_by_user",
		mcp.WithDescription("List all HN findings for a user"),
		mcp.WithNumber("user_id", mcp.Required(), mcp.Description("User ID")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		userID := req.GetInt("user_id", 0)
		findings, err := db.HN.ListHNFindingsByUser(int64(userID))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		b, _ := json.Marshal(findings)
		return mcp.NewToolResultText(string(b)), nil
	})

	// hn_finding_list_by_thread
	s.AddTool(mcp.NewTool("hn_finding_list_by_thread",
		mcp.WithDescription("List all HN findings in a thread"),
		mcp.WithNumber("thread_id", mcp.Required(), mcp.Description("Thread ID")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		threadID := req.GetInt("thread_id", 0)
		findings, err := db.HN.ListHNFindingsByThread(int64(threadID))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		b, _ := json.Marshal(findings)
		return mcp.NewToolResultText(string(b)), nil
	})
}

// --- Health Tools ---

func registerHealthTools(s *server.MCPServer, db *sqlite.DB) {
	// gene_snippet_create
	s.AddTool(mcp.NewTool("gene_snippet_create",
		mcp.WithDescription("Create a gene snippet for a patient"),
		mcp.WithNumber("patient_id", mcp.Required(), mcp.Description("Patient (user) ID")),
		mcp.WithString("gene", mcp.Required(), mcp.Description("Gene name, e.g. 'MTHFR', 'COMT'")),
		mcp.WithString("variant", mcp.Description("Variant, e.g. 'C677T'")),
		mcp.WithString("genotype", mcp.Description("Genotype, e.g. 'C/T'")),
		mcp.WithString("ai_hint", mcp.Description("Study references: what this gene does")),
		mcp.WithString("ai_note", mcp.Description("AI's own note")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		gene, _ := req.RequireString("gene")
		gs := &hnmentalhealth.GeneSnippet{
			PatientID: int64(req.GetInt("patient_id", 0)),
			Gene:      gene,
			Variant:   req.GetString("variant", ""),
			Genotype:  req.GetString("genotype", ""),
			AIHint:    req.GetString("ai_hint", ""),
			AINote:    req.GetString("ai_note", ""),
		}
		if err := db.Health.CreateGeneSnippet(gs); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		b, _ := json.Marshal(gs)
		return mcp.NewToolResultText(string(b)), nil
	})

	// gene_snippet_list_by_patient
	s.AddTool(mcp.NewTool("gene_snippet_list_by_patient",
		mcp.WithDescription("List all gene snippets for a patient"),
		mcp.WithNumber("patient_id", mcp.Required(), mcp.Description("Patient ID")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		patientID := req.GetInt("patient_id", 0)
		snippets, err := db.Health.ListGeneSnippetsByPatient(int64(patientID))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		b, _ := json.Marshal(snippets)
		return mcp.NewToolResultText(string(b)), nil
	})

	// gene_pattern_list
	s.AddTool(mcp.NewTool("gene_pattern_list",
		mcp.WithDescription("List all gene patterns (aggregations of gene snippets)"),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		patterns, err := db.Health.ListGenePatterns()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		b, _ := json.Marshal(patterns)
		return mcp.NewToolResultText(string(b)), nil
	})

	// gene_pattern_create
	s.AddTool(mcp.NewTool("gene_pattern_create",
		mcp.WithDescription("Create a new gene pattern. Use when the AI discovers a new gene combination."),
		mcp.WithString("name", mcp.Required(), mcp.Description("Unique pattern name, e.g. 'Methylierungs-Störung'")),
		mcp.WithString("description", mcp.Description("What this pattern means")),
		mcp.WithString("ai_hint", mcp.Description("Study references for the gene interplay")),
		mcp.WithString("ai_note", mcp.Description("AI's own note")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, _ := req.RequireString("name")
		gp := &hnmentalhealth.GenePattern{
			Name:        name,
			Description: req.GetString("description", ""),
			AIHint:      req.GetString("ai_hint", ""),
			AINote:      req.GetString("ai_note", ""),
		}
		if err := db.Health.CreateGenePattern(gp); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		b, _ := json.Marshal(gp)
		return mcp.NewToolResultText(string(b)), nil
	})

	// pattern_snippet_add
	s.AddTool(mcp.NewTool("pattern_snippet_add",
		mcp.WithDescription("Add a gene snippet to a gene pattern"),
		mcp.WithNumber("pattern_id", mcp.Required(), mcp.Description("Gene pattern ID")),
		mcp.WithNumber("snippet_id", mcp.Required(), mcp.Description("Gene snippet ID")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		patternID := int64(req.GetInt("pattern_id", 0))
		snippetID := int64(req.GetInt("snippet_id", 0))
		if err := db.Health.AddSnippetToPattern(patternID, snippetID); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf(`{"ok":true,"pattern_id":%d,"snippet_id":%d}`, patternID, snippetID)), nil
	})

	// patient_pattern_add
	s.AddTool(mcp.NewTool("patient_pattern_add",
		mcp.WithDescription("Assign a gene pattern to a patient"),
		mcp.WithNumber("patient_id", mcp.Required(), mcp.Description("Patient ID")),
		mcp.WithNumber("pattern_id", mcp.Required(), mcp.Description("Gene pattern ID")),
		mcp.WithString("ai_note", mcp.Description("AI note specific to this patient-pattern combination")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		patientID := int64(req.GetInt("patient_id", 0))
		patternID := int64(req.GetInt("pattern_id", 0))
		if err := db.Health.AddPatternToPatient(patientID, patternID); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if note := req.GetString("ai_note", ""); note != "" {
			db.Health.UpdatePatientPatternNote(patientID, patternID, note)
		}
		return mcp.NewToolResultText(fmt.Sprintf(`{"ok":true,"patient_id":%d,"pattern_id":%d}`, patientID, patternID)), nil
	})

	// biomarker_create
	s.AddTool(mcp.NewTool("biomarker_create",
		mcp.WithDescription("Create a biomarker (blood/lab value) for a patient"),
		mcp.WithNumber("patient_id", mcp.Required(), mcp.Description("Patient ID")),
		mcp.WithString("name", mcp.Required(), mcp.Description("Biomarker name, e.g. 'Folsäure', 'B12'")),
		mcp.WithNumber("value", mcp.Required(), mcp.Description("Measured value")),
		mcp.WithString("unit", mcp.Description("Unit, e.g. 'ng/ml'")),
		mcp.WithNumber("reference_low", mcp.Description("Lower reference value")),
		mcp.WithNumber("reference_high", mcp.Description("Upper reference value")),
		mcp.WithString("measured_at", mcp.Description("ISO date of measurement")),
		mcp.WithString("ai_hint", mcp.Description("AI hint")),
		mcp.WithString("ai_note", mcp.Description("AI note")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, _ := req.RequireString("name")
		b := &hnmentalhealth.Biomarker{
			PatientID:     int64(req.GetInt("patient_id", 0)),
			Name:          name,
			Value:         float64(req.GetInt("value", 0)),
			Unit:          req.GetString("unit", ""),
			ReferenceLow:  float64(req.GetInt("reference_low", 0)),
			ReferenceHigh: float64(req.GetInt("reference_high", 0)),
			MeasuredAt:    req.GetString("measured_at", ""),
			AIHint:        req.GetString("ai_hint", ""),
			AINote:        req.GetString("ai_note", ""),
		}
		if err := db.Health.CreateBiomarker(b); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		jsonB, _ := json.Marshal(b)
		return mcp.NewToolResultText(string(jsonB)), nil
	})

	// biomarker_list_out_of_range
	s.AddTool(mcp.NewTool("biomarker_list_out_of_range",
		mcp.WithDescription("List all biomarkers outside reference range for a patient"),
		mcp.WithNumber("patient_id", mcp.Required(), mcp.Description("Patient ID")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		patientID := req.GetInt("patient_id", 0)
		biomarkers, err := db.Health.ListBiomarkersOutOfRange(int64(patientID))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		b, _ := json.Marshal(biomarkers)
		return mcp.NewToolResultText(string(b)), nil
	})

	// health_finding_create
	s.AddTool(mcp.NewTool("health_finding_create",
		mcp.WithDescription("Create a health finding — a metabolic/psychological observation about a patient"),
		mcp.WithNumber("patient_id", mcp.Required(), mcp.Description("Patient ID")),
		mcp.WithNumber("pattern_id", mcp.Description("Gene pattern ID")),
		mcp.WithNumber("confidence_id", mcp.Description("Confidence level ID")),
		mcp.WithNumber("icd11_id", mcp.Description("ICD-11 code ID")),
		mcp.WithString("etiology", mcp.Description("metabolic, psychological, mixed, or unknown")),
		mcp.WithString("evidence_chain", mcp.Description("JSON array of evidence sources")),
		mcp.WithString("evidence", mcp.Description("Free-text description")),
		mcp.WithString("ai_hint", mcp.Description("AI hint")),
		mcp.WithString("ai_note", mcp.Description("AI note")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		f := &hnmentalhealth.HealthFinding{
			PatientID:     int64(req.GetInt("patient_id", 0)),
			PatternID:     optIntPtr(req, "pattern_id"),
			ConfidenceID:  optIntPtr(req, "confidence_id"),
			ICD11ID:       optIntPtr(req, "icd11_id"),
			Etiology:      req.GetString("etiology", ""),
			EvidenceChain: req.GetString("evidence_chain", ""),
			Evidence:      req.GetString("evidence", ""),
			AIHint:        req.GetString("ai_hint", ""),
			AINote:        req.GetString("ai_note", ""),
		}
		if err := db.Health.CreateHealthFinding(f); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		b, _ := json.Marshal(f)
		return mcp.NewToolResultText(string(b)), nil
	})

	// health_finding_list_by_patient
	s.AddTool(mcp.NewTool("health_finding_list_by_patient",
		mcp.WithDescription("List all health findings for a patient"),
		mcp.WithNumber("patient_id", mcp.Required(), mcp.Description("Patient ID")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		patientID := req.GetInt("patient_id", 0)
		findings, err := db.Health.ListHealthFindingsByPatient(int64(patientID))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		b, _ := json.Marshal(findings)
		return mcp.NewToolResultText(string(b)), nil
	})
}
