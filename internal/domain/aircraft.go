package domain

// WeightCategory represents aircraft loading
type WeightCategory string

const (
	WeightLight   WeightCategory = "light"
	WeightTypical WeightCategory = "typical"
	WeightHeavy   WeightCategory = "heavy"
)

// RunwayRequirements contains runway length requirements by weight
type RunwayRequirements struct {
	Light   int `json:"light"`
	Typical int `json:"typical"`
	Heavy   int `json:"heavy"`
}

// Aircraft represents an aircraft type with its performance limits
type Aircraft struct {
	TypeCode           string              `json:"type_code"`
	Name               string              `json:"name"`
	Category           string              `json:"category"` // SEP, MEP, JET, etc.
	MinRunwayFt        int                 `json:"min_runway_ft"`
	MinRunwayByWeight  *RunwayRequirements `json:"min_runway_by_weight,omitempty"`
	MaxCrosswindKt     int                 `json:"max_crosswind_kt"`
	MaxTailwindKt      int                 `json:"max_tailwind_kt"`
	VFROnly            bool                `json:"vfr_only"`
	MinVisibilitySM    float64             `json:"min_visibility_sm"`
	MinCeilingFt       int                 `json:"min_ceiling_ft"`
}

