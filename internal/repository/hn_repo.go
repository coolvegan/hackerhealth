// Package repository — HN Repository Interfaces.
package repository

import "github.com/coolvegan/hn-mentalhealth"

// HNRepository fasst alle HN-spezifischen Repositories zusammen.
type HNRepository interface {
	PatternRepository
	EmotionRepository
	CharacterizationRepository
	InteractionRepository
	ThreadRepository
	ThreadQualityRepository
	ThreadEmotionRepository
	HNFindingRepository
}

type PatternRepository interface {
	CreatePattern(p *hnmentalhealth.Pattern) error
	GetPatternByID(id int64) (*hnmentalhealth.Pattern, error)
	GetPatternByName(name string) (*hnmentalhealth.Pattern, error)
	UpdatePattern(p *hnmentalhealth.Pattern) error
	DeletePattern(id int64) error
	ListPatterns() ([]hnmentalhealth.Pattern, error)
}

type EmotionRepository interface {
	CreateEmotion(e *hnmentalhealth.Emotion) error
	GetEmotionByID(id int64) (*hnmentalhealth.Emotion, error)
	GetEmotionByName(name string) (*hnmentalhealth.Emotion, error)
	UpdateEmotion(e *hnmentalhealth.Emotion) error
	DeleteEmotion(id int64) error
	ListEmotions() ([]hnmentalhealth.Emotion, error)
}

type CharacterizationRepository interface {
	CreateCharacterization(c *hnmentalhealth.Characterization) error
	GetCharacterizationByID(id int64) (*hnmentalhealth.Characterization, error)
	GetCharacterizationByName(name string) (*hnmentalhealth.Characterization, error)
	UpdateCharacterization(c *hnmentalhealth.Characterization) error
	DeleteCharacterization(id int64) error
	ListCharacterizations() ([]hnmentalhealth.Characterization, error)
}

type InteractionRepository interface {
	CreateInteraction(i *hnmentalhealth.Interaction) error
	GetInteractionByID(id int64) (*hnmentalhealth.Interaction, error)
	GetInteractionByName(name string) (*hnmentalhealth.Interaction, error)
	UpdateInteraction(i *hnmentalhealth.Interaction) error
	DeleteInteraction(id int64) error
	ListInteractions() ([]hnmentalhealth.Interaction, error)
}

type ThreadRepository interface {
	CreateThread(t *hnmentalhealth.Thread) error
	GetThreadByID(id int64) (*hnmentalhealth.Thread, error)
	GetThreadByStoryID(storyID int64) (*hnmentalhealth.Thread, error)
	UpdateThread(t *hnmentalhealth.Thread) error
	DeleteThread(id int64) error
	ListThreads(offset, limit int) ([]hnmentalhealth.Thread, error)
}

type ThreadQualityRepository interface {
	CreateThreadQuality(tq *hnmentalhealth.ThreadQuality) error
	GetThreadQualityByThreadID(threadID int64) (*hnmentalhealth.ThreadQuality, error)
	UpdateThreadQuality(tq *hnmentalhealth.ThreadQuality) error
	DeleteThreadQuality(threadID int64) error
}

type ThreadEmotionRepository interface {
	SetThreadEmotion(te *hnmentalhealth.ThreadEmotion) error
	GetThreadEmotions(threadID int64) ([]hnmentalhealth.ThreadEmotion, error)
	DeleteThreadEmotion(threadID, emotionID int64) error
}

type HNFindingRepository interface {
	CreateHNFinding(f *hnmentalhealth.HNFinding) error
	GetHNFindingByID(id int64) (*hnmentalhealth.HNFinding, error)
	UpdateHNFinding(f *hnmentalhealth.HNFinding) error
	DeleteHNFinding(id int64) error
	ListHNFindingsByUser(userID int64) ([]hnmentalhealth.HNFinding, error)
	ListHNFindingsByThread(threadID int64) ([]hnmentalhealth.HNFinding, error)
	ListHNFindingsByPattern(patternID int64) ([]hnmentalhealth.HNFinding, error)
	CountHNFindingsByUserPattern(userID, patternID int64) (int, error)
}
