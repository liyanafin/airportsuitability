package domain

type Airport struct {
	ICAO        string   `json:"icao"`
	Name        string   `json:"name"`
	ElevationFt int      `json:"elevation_ft"`
	Runways     []Runway `json:"runways"`
}

// Runway represents a runway at an airport
type Runway struct {
	ID       string `json:"id"`
	Heading  int    `json:"heading"` // (1-360)
	LengthFt int    `json:"length_ft"`
	WidthFt  int    `json:"width_ft"`
	Surface  string `json:"surface"` // ASPH, CONC, TURF, etc.
	Closed   bool   `json:"closed"`
}

func (a *Airport) GetAllRunways() []Runway {
	return a.Runways
}
