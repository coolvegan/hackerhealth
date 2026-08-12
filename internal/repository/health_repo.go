// Package repository — Health Repository Interfaces.
package repository

import "github.com/coolvegan/hn-mentalhealth"

// HealthRepository fasst alle Health-spezifischen Repositories zusammen.
type HealthRepository interface {
	GeneSnippetRepository
	GenePatternRepository
	PatternSnippetRepository
	PatientPatternRepository
	BiomarkerRepository
	HealthConditionRepository
	LifeContextRepository
	HealthFindingRepository
}

type GeneSnippetRepository interface {
	CreateGeneSnippet(gs *hnmentalhealth.GeneSnippet) error
	GetGeneSnippetByID(id int64) (*hnmentalhealth.GeneSnippet, error)
	UpdateGeneSnippet(gs *hnmentalhealth.GeneSnippet) error
	DeleteGeneSnippet(id int64) error
	ListGeneSnippetsByPatient(patientID int64) ([]hnmentalhealth.GeneSnippet, error)
}

type GenePatternRepository interface {
	CreateGenePattern(gp *hnmentalhealth.GenePattern) error
	GetGenePatternByID(id int64) (*hnmentalhealth.GenePattern, error)
	GetGenePatternByName(name string) (*hnmentalhealth.GenePattern, error)
	UpdateGenePattern(gp *hnmentalhealth.GenePattern) error
	DeleteGenePattern(id int64) error
	ListGenePatterns() ([]hnmentalhealth.GenePattern, error)
}

type PatternSnippetRepository interface {
	AddSnippetToPattern(patternID, snippetID int64) error
	RemoveSnippetFromPattern(patternID, snippetID int64) error
	ListSnippetsByPattern(patternID int64) ([]hnmentalhealth.GeneSnippet, error)
	ListPatternsBySnippet(snippetID int64) ([]hnmentalhealth.GenePattern, error)
}

type PatientPatternRepository interface {
	AddPatternToPatient(patientID, patternID int64) error
	RemovePatternFromPatient(patientID, patternID int64) error
	UpdatePatientPatternNote(patientID, patternID int64, note string) error
	ListPatternsByPatient(patientID int64) ([]hnmentalhealth.GenePattern, error)
	ListPatientsByPattern(patternID int64) ([]hnmentalhealth.User, error)
}

type BiomarkerRepository interface {
	CreateBiomarker(b *hnmentalhealth.Biomarker) error
	GetBiomarkerByID(id int64) (*hnmentalhealth.Biomarker, error)
	UpdateBiomarker(b *hnmentalhealth.Biomarker) error
	DeleteBiomarker(id int64) error
	ListBiomarkersByPatient(patientID int64) ([]hnmentalhealth.Biomarker, error)
	ListBiomarkersOutOfRange(patientID int64) ([]hnmentalhealth.Biomarker, error)
}

type HealthConditionRepository interface {
	CreateHealthCondition(hc *hnmentalhealth.HealthCondition) error
	GetHealthConditionByID(id int64) (*hnmentalhealth.HealthCondition, error)
	UpdateHealthCondition(hc *hnmentalhealth.HealthCondition) error
	DeleteHealthCondition(id int64) error
	ListHealthConditionsByPatient(patientID int64) ([]hnmentalhealth.HealthCondition, error)
}

type LifeContextRepository interface {
	CreateLifeContext(lc *hnmentalhealth.LifeContext) error
	GetLifeContextByID(id int64) (*hnmentalhealth.LifeContext, error)
	UpdateLifeContext(lc *hnmentalhealth.LifeContext) error
	DeleteLifeContext(id int64) error
	ListLifeContextByPatient(patientID int64) ([]hnmentalhealth.LifeContext, error)
}

type HealthFindingRepository interface {
	CreateHealthFinding(f *hnmentalhealth.HealthFinding) error
	GetHealthFindingByID(id int64) (*hnmentalhealth.HealthFinding, error)
	UpdateHealthFinding(f *hnmentalhealth.HealthFinding) error
	DeleteHealthFinding(id int64) error
	ListHealthFindingsByPatient(patientID int64) ([]hnmentalhealth.HealthFinding, error)
	ListHealthFindingsByPattern(patternID int64) ([]hnmentalhealth.HealthFinding, error)
	ListHealthFindingsByEtiology(etiology string) ([]hnmentalhealth.HealthFinding, error)
}
