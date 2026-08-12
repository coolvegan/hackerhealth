// Package repository — Shared Repository Interfaces.
package repository

import "github.com/coolvegan/hn-mentalhealth"

// SharedRepository fasst die domänen-übergreifenden Repositories zusammen.
type SharedRepository interface {
	UserRepository
	ConfidenceLevelRepository
	ICD11CodeRepository
	AIKnowledgeRepository
}

// UserRepository verwaltet Personen (HN-Kommentatoren und Patienten).
type UserRepository interface {
	CreateUser(user *hnmentalhealth.User) error
	GetUserByID(id int64) (*hnmentalhealth.User, error)
	GetUserByUsername(username string) (*hnmentalhealth.User, error)
	UpdateUser(user *hnmentalhealth.User) error
	DeleteUser(id int64) error
	ListUsers(offset, limit int) ([]hnmentalhealth.User, error)
}

// ConfidenceLevelRepository verwaltet Konfidenz-Stufen.
type ConfidenceLevelRepository interface {
	CreateConfidenceLevel(cl *hnmentalhealth.ConfidenceLevel) error
	GetConfidenceLevelByID(id int64) (*hnmentalhealth.ConfidenceLevel, error)
	GetConfidenceLevelByName(name string) (*hnmentalhealth.ConfidenceLevel, error)
	UpdateConfidenceLevel(cl *hnmentalhealth.ConfidenceLevel) error
	DeleteConfidenceLevel(id int64) error
	ListConfidenceLevels() ([]hnmentalhealth.ConfidenceLevel, error)
}

// ICD11CodeRepository verwaltet ICD-11-Referenzen.
type ICD11CodeRepository interface {
	CreateICD11Code(code *hnmentalhealth.ICD11Code) error
	GetICD11CodeByID(id int64) (*hnmentalhealth.ICD11Code, error)
	GetICD11CodeByCode(code string) (*hnmentalhealth.ICD11Code, error)
	UpdateICD11Code(code *hnmentalhealth.ICD11Code) error
	DeleteICD11Code(id int64) error
	ListICD11Codes() ([]hnmentalhealth.ICD11Code, error)
}

// AIKnowledgeRepository verwaltet das selbst-beschreibende Daten-Wörterbuch.
type AIKnowledgeRepository interface {
	CreateAIKnowledge(entry *hnmentalhealth.AIKnowledge) error
	GetAIKnowledgeByID(id int64) (*hnmentalhealth.AIKnowledge, error)
	GetAIKnowledgeByKey(tableName, columnName, value string) (*hnmentalhealth.AIKnowledge, error)
	UpdateAIKnowledge(entry *hnmentalhealth.AIKnowledge) error
	DeleteAIKnowledge(id int64) error
	ListAIKnowledgeByTable(tableName string) ([]hnmentalhealth.AIKnowledge, error)
	SearchAIKnowledge(query string, limit int) ([]hnmentalhealth.AIKnowledge, error)
}
