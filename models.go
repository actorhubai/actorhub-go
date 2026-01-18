// Package actorhub provides a Go client for the ActorHub.ai API.
package actorhub

import "time"

// TrainingStatus represents the status of an Actor Pack training job.
type TrainingStatus string

const (
	TrainingStatusQueued     TrainingStatus = "QUEUED"
	TrainingStatusProcessing TrainingStatus = "PROCESSING"
	TrainingStatusCompleted  TrainingStatus = "COMPLETED"
	TrainingStatusFailed     TrainingStatus = "FAILED"
)

// ProtectionLevel represents the identity protection tier.
type ProtectionLevel string

const (
	ProtectionLevelFree       ProtectionLevel = "free"
	ProtectionLevelPro        ProtectionLevel = "pro"
	ProtectionLevelEnterprise ProtectionLevel = "enterprise"
)

// LicenseType represents the type of license.
type LicenseType string

const (
	LicenseTypeStandard  LicenseType = "standard"
	LicenseTypeExtended  LicenseType = "extended"
	LicenseTypeExclusive LicenseType = "exclusive"
)

// UsageType represents the usage category for licensing.
type UsageType string

const (
	UsageTypePersonal    UsageType = "personal"
	UsageTypeEditorial   UsageType = "editorial"
	UsageTypeCommercial  UsageType = "commercial"
	UsageTypeEducational UsageType = "educational"
)

// FaceBBox represents face bounding box coordinates.
type FaceBBox struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// LicenseOption represents a license option with pricing.
type LicenseOption struct {
	Type           LicenseType `json:"type"`
	PriceUSD       float64     `json:"price_usd"`
	DurationDays   int         `json:"duration_days"`
	MaxImpressions *int        `json:"max_impressions,omitempty"`
}

// VerifyResult represents an individual identity verification result.
type VerifyResult struct {
	Protected         bool            `json:"protected"`
	IdentityID        *string         `json:"identity_id,omitempty"`
	SimilarityScore   *float64        `json:"similarity_score,omitempty"`
	DisplayName       *string         `json:"display_name,omitempty"`
	LicenseRequired   bool            `json:"license_required"`
	BlockedCategories []string        `json:"blocked_categories"`
	LicenseOptions    []LicenseOption `json:"license_options"`
	FaceBBox          *FaceBBox       `json:"face_bbox,omitempty"`
}

// VerifyResponse is the response from identity verification.
type VerifyResponse struct {
	Protected      bool           `json:"protected"`
	FacesDetected  int            `json:"faces_detected"`
	Identities     []VerifyResult `json:"identities"`
	ResponseTimeMs int            `json:"response_time_ms"`
	RequestID      string         `json:"request_id"`
}

// ConsentDetails represents consent permissions for an identity.
type ConsentDetails struct {
	CommercialUse   bool `json:"commercial_use"`
	AITraining      bool `json:"ai_training"`
	VideoGeneration bool `json:"video_generation"`
	Deepfake        bool `json:"deepfake"`
}

// ConsentRestrictions represents consent restrictions.
type ConsentRestrictions struct {
	BlockedCategories []string `json:"blocked_categories"`
	BlockedRegions    []string `json:"blocked_regions"`
	BlockedBrands     []string `json:"blocked_brands"`
}

// ConsentLicenseInfo represents license availability information.
type ConsentLicenseInfo struct {
	Available bool                `json:"available"`
	URL       *string             `json:"url,omitempty"`
	Pricing   map[string]float64  `json:"pricing,omitempty"`
}

// ConsentResult represents an individual consent check result.
type ConsentResult struct {
	Protected       bool                `json:"protected"`
	IdentityID      *string             `json:"identity_id,omitempty"`
	SimilarityScore *float64            `json:"similarity_score,omitempty"`
	Consent         ConsentDetails      `json:"consent"`
	Restrictions    ConsentRestrictions `json:"restrictions"`
	License         ConsentLicenseInfo  `json:"license"`
}

// ConsentCheckResponse is the response from consent check.
type ConsentCheckResponse struct {
	RequestID          string          `json:"request_id"`
	Protected          bool            `json:"protected"`
	FacesDetected      int             `json:"faces_detected"`
	Faces              []ConsentResult `json:"faces"`
	ResponseTimeMs     int             `json:"response_time_ms"`
	RateLimitRemaining *int            `json:"rate_limit_remaining,omitempty"`
}

// IdentityResponse represents identity details.
type IdentityResponse struct {
	ID                 string          `json:"id"`
	DisplayName        string          `json:"display_name"`
	ProfileImageURL    *string         `json:"profile_image_url,omitempty"`
	Status             string          `json:"status"`
	ProtectionLevel    ProtectionLevel `json:"protection_level"`
	ProtectionMode     string          `json:"protection_mode"`
	TotalVerifications int             `json:"total_verifications"`
	TotalLicenses      int             `json:"total_licenses"`
	TotalRevenue       float64         `json:"total_revenue"`
	AllowCommercial    bool            `json:"allow_commercial"`
	AllowAITraining    bool            `json:"allow_ai_training"`
	CreatedAt          *time.Time      `json:"created_at,omitempty"`
}

// MarketplaceListingResponse represents marketplace listing details.
type MarketplaceListingResponse struct {
	ID              string     `json:"id"`
	IdentityID      string     `json:"identity_id"`
	Title           string     `json:"title"`
	Description     *string    `json:"description,omitempty"`
	Category        string     `json:"category"`
	Tags            []string   `json:"tags"`
	BasePriceUSD    float64    `json:"base_price_usd"`
	DisplayName     string     `json:"display_name"`
	ProfileImageURL *string    `json:"profile_image_url,omitempty"`
	Featured        bool       `json:"featured"`
	ViewCount       int        `json:"view_count"`
	LicenseCount    int        `json:"license_count"`
	Rating          *float64   `json:"rating,omitempty"`
	CreatedAt       *time.Time `json:"created_at,omitempty"`
}

// LicenseResponse represents license details.
type LicenseResponse struct {
	ID                 string      `json:"id"`
	IdentityID         string      `json:"identity_id"`
	IdentityName       string      `json:"identity_name"`
	LicenseType        LicenseType `json:"license_type"`
	UsageType          UsageType   `json:"usage_type"`
	Status             string      `json:"status"`
	ProjectName        string      `json:"project_name"`
	ProjectDescription *string     `json:"project_description,omitempty"`
	AllowedPlatforms   []string    `json:"allowed_platforms"`
	MaxImpressions     *int        `json:"max_impressions,omitempty"`
	MaxOutputs         *int        `json:"max_outputs,omitempty"`
	PriceUSD           float64     `json:"price_usd"`
	StartsAt           *time.Time  `json:"starts_at,omitempty"`
	ExpiresAt          *time.Time  `json:"expires_at,omitempty"`
	CreatedAt          *time.Time  `json:"created_at,omitempty"`
}

// ActorPackComponents represents Actor Pack component availability.
type ActorPackComponents struct {
	Face   bool `json:"face"`
	Voice  bool `json:"voice"`
	Motion bool `json:"motion"`
}

// ActorPackResponse represents Actor Pack details.
type ActorPackResponse struct {
	ID                   string              `json:"id"`
	IdentityID           string              `json:"identity_id"`
	Name                 string              `json:"name"`
	Description          *string             `json:"description,omitempty"`
	TrainingStatus       TrainingStatus      `json:"training_status"`
	TrainingProgress     int                 `json:"training_progress"`
	TrainingImagesCount  int                 `json:"training_images_count"`
	TrainingAudioSeconds int                 `json:"training_audio_seconds"`
	Components           ActorPackComponents `json:"components"`
	LoRAModelURL         *string             `json:"lora_model_url,omitempty"`
	TotalDownloads       int                 `json:"total_downloads"`
	IsAvailable          bool                `json:"is_available"`
	TrainingError        *string             `json:"training_error,omitempty"`
	CreatedAt            *time.Time          `json:"created_at,omitempty"`
}

// PurchaseResponse is the license purchase response.
type PurchaseResponse struct {
	CheckoutURL    string                 `json:"checkout_url"`
	SessionID      string                 `json:"session_id"`
	PriceUSD       float64                `json:"price_usd"`
	LicenseDetails map[string]interface{} `json:"license_details"`
}

// VerifyRequest represents the request for identity verification.
type VerifyRequest struct {
	ImageURL              string `json:"image_url,omitempty"`
	ImageBase64           string `json:"image_base64,omitempty"`
	IncludeLicenseOptions bool   `json:"include_license_options,omitempty"`
}

// ConsentCheckRequest represents the request for consent check.
type ConsentCheckRequest struct {
	ImageURL      string    `json:"image_url,omitempty"`
	ImageBase64   string    `json:"image_base64,omitempty"`
	FaceEmbedding []float64 `json:"face_embedding,omitempty"`
	Platform      string    `json:"platform"`
	IntendedUse   string    `json:"intended_use"`
	Region        string    `json:"region,omitempty"`
}

// MarketplaceListRequest represents the request for marketplace listing.
type MarketplaceListRequest struct {
	Query    string   `json:"query,omitempty"`
	Category string   `json:"category,omitempty"`
	Tags     []string `json:"tags,omitempty"`
	Featured *bool    `json:"featured,omitempty"`
	MinPrice *float64 `json:"min_price,omitempty"`
	MaxPrice *float64 `json:"max_price,omitempty"`
	SortBy   string   `json:"sort_by,omitempty"`
	Page     int      `json:"page,omitempty"`
	Limit    int      `json:"limit,omitempty"`
}

// PurchaseLicenseRequest represents the request for license purchase.
type PurchaseLicenseRequest struct {
	IdentityID         string   `json:"identity_id"`
	LicenseType        string   `json:"license_type"`
	UsageType          string   `json:"usage_type"`
	ProjectName        string   `json:"project_name"`
	ProjectDescription string   `json:"project_description"`
	DurationDays       int      `json:"duration_days,omitempty"`
	AllowedPlatforms   []string `json:"allowed_platforms,omitempty"`
	MaxImpressions     *int     `json:"max_impressions,omitempty"`
	MaxOutputs         *int     `json:"max_outputs,omitempty"`
}
