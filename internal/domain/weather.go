package domain

import (
	"math"
	"regexp"
	"strconv"
	"strings"
)

// Weather represents parsed weather data from a METAR
type Weather struct {
	RawMETAR      string  `json:"metar"`
	WindDirection int     `json:"wind_direction"`
	WindSpeedKt   int     `json:"wind_speed_kt"`
	GustKt        int     `json:"gust_kt"`
	VisibilitySM  float64 `json:"visibility_sm"`
	CeilingFt     int     `json:"ceiling_ft"`
	VariableWind  bool    `json:"variable_wind"`
	CalmWind      bool    `json:"calm_wind"`
}

// CalculateWindComponents calculates headwind and crosswind components
// given a runway heading. Positive headwind is beneficial, negative is tailwind.
// Crosswind is always returned as positive (absolute value).
func (w *Weather) CalculateWindComponents(runwayHeading int) (headwind, crosswind int) {
	if w.CalmWind || w.WindSpeedKt == 0 {
		return 0, 0
	}

	if w.VariableWind {
		// Variable winds - assume worst case (full crosswind)
		return 0, w.WindSpeedKt
	}

	// Calculate angle between wind direction and runway heading
	angleDiff := float64(w.WindDirection - runwayHeading)
	angleRad := angleDiff * math.Pi / 180

	// Headwind component (positive = headwind, negative = tailwind)
	headwindFloat := float64(w.WindSpeedKt) * math.Cos(angleRad)
	headwind = int(math.Round(headwindFloat))

	// Crosswind component (absolute value)
	crosswindFloat := math.Abs(float64(w.WindSpeedKt) * math.Sin(angleRad))
	crosswind = int(math.Round(crosswindFloat))

	return headwind, crosswind
}
// EffectiveWindSpeed returns the wind speed to use for calculations
// (uses gusts if present, otherwise steady wind)
func (w *Weather) EffectiveWindSpeed() int {
	if w.GustKt > 0 {
		return w.GustKt
	}
	return w.WindSpeedKt
}

// ParseMETAR parses a raw METAR string into a Weather struct
func ParseMETAR(metar string) (*Weather, error) {
	weather := &Weather{
		RawMETAR:   metar,
		CeilingFt:  10000, // Default high ceiling if none reported
	}

	// Parse wind
	parseWind(metar, weather)

	// Parse visibility
	parseVisibility(metar, weather)

	// Parse ceiling
	parseCeiling(metar, weather)

	return weather, nil
}

func parseWind(metar string, weather *Weather) {
	// Check for calm wind
	if strings.Contains(metar, "00000KT") {
		weather.CalmWind = true
		weather.WindDirection = 0
		weather.WindSpeedKt = 0
		return
	}

	// Check for variable wind (VRB)
	vrbPattern := regexp.MustCompile(`VRB(\d{2,3})KT`)
	if matches := vrbPattern.FindStringSubmatch(metar); matches != nil {
		weather.VariableWind = true
		weather.WindSpeedKt, _ = strconv.Atoi(matches[1])
		return
	}

	// Standard wind pattern: DDDSSGGGKt or DDDSSKT
	windPattern := regexp.MustCompile(`(\d{3})(\d{2,3})(G(\d{2,3}))?KT`)
	if matches := windPattern.FindStringSubmatch(metar); matches != nil {
		weather.WindDirection, _ = strconv.Atoi(matches[1])
		weather.WindSpeedKt, _ = strconv.Atoi(matches[2])
		if matches[4] != "" {
			weather.GustKt, _ = strconv.Atoi(matches[4])
		}
	}
}

func parseVisibility(metar string, weather *Weather) {
	parts := strings.Fields(metar)

	for _, part := range parts {
		// Check for visibility in statute miles (e.g., 10SM, 3SM, 1/2SM)
		if strings.HasSuffix(part, "SM") {
			visStr := strings.TrimSuffix(part, "SM")

			// Handle fractional visibility
			if strings.Contains(visStr, "/") {
				fractionParts := strings.Split(visStr, "/")
				if len(fractionParts) == 2 {
					num, _ := strconv.ParseFloat(fractionParts[0], 64)
					den, _ := strconv.ParseFloat(fractionParts[1], 64)
					if den != 0 {
						weather.VisibilitySM = num / den
					}
				}
			} else {
				// Handle mixed numbers like "1 1/2SM" by checking previous part
				vis, err := strconv.ParseFloat(visStr, 64)
				if err == nil {
					weather.VisibilitySM = vis
				}
			}
			break
		}
	}
}

func parseCeiling(metar string, weather *Weather) {
	// Ceiling is the lowest BKN (broken) or OVC (overcast) layer
	// Also check for VV (vertical visibility) for obscured sky

	// Pattern for cloud layers
	cloudPattern := regexp.MustCompile(`(BKN|OVC|VV)(\d{3})`)
	matches := cloudPattern.FindAllStringSubmatch(metar, -1)

	lowestCeiling := 10000 // Default high ceiling
	for _, match := range matches {
		if len(match) >= 3 {
			height, _ := strconv.Atoi(match[2])
			height *= 100 // Convert from hundreds of feet

			if height < lowestCeiling {
				lowestCeiling = height
			}
		}
	}

	weather.CeilingFt = lowestCeiling
}