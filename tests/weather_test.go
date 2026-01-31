package tests

import (
	"math"
	"testing"
	"github.com/liyanafin/airportsuitability/internal/domain"
)

func TestWeatherBasics(t *testing.T) {
	weather := Weather{
		RawMETAR:      "KJFK 261456Z 09025G35KT 10SM FEW025 BKN050 12/05 A2992",
		WindDirection: 90,
		WindSpeedKt:   25,
		GustKt:        35,
		VisibilitySM:  10,
		CeilingFt:     5000,
	}

	t.Run("has raw METAR", func(t *testing.T) {
		if weather.RawMETAR == "" {
			t.Error("expected raw METAR to be set")
		}
	})

	t.Run("has wind direction", func(t *testing.T) {
		if weather.WindDirection != 90 {
			t.Errorf("expected wind direction 90, got %d", weather.WindDirection)
		}
	})

	t.Run("has wind speed and gusts", func(t *testing.T) {
		if weather.WindSpeedKt != 25 {
			t.Errorf("expected wind speed 25kt, got %d", weather.WindSpeedKt)
		}
		if weather.GustKt != 35 {
			t.Errorf("expected gusts 35kt, got %d", weather.GustKt)
		}
	})
}

func TestWindComponents(t *testing.T) {
	tests := []struct {
		name              string
		windDirection     int
		windSpeed         int
		runwayHeading     int
		expectedHeadwind  int
		expectedCrosswind int
	}{
		{
			name:              "direct headwind",
			windDirection:     360,
			windSpeed:         20,
			runwayHeading:     360,
			expectedHeadwind:  20,
			expectedCrosswind: 0,
		},
		{
			name:              "direct tailwind",
			windDirection:     180,
			windSpeed:         20,
			runwayHeading:     360,
			expectedHeadwind:  -20,
			expectedCrosswind: 0,
		},
		{
			name:              "direct crosswind from right",
			windDirection:     90,
			windSpeed:         20,
			runwayHeading:     360,
			expectedHeadwind:  0,
			expectedCrosswind: 20,
		},
		{
			name:              "direct crosswind from left",
			windDirection:     270,
			windSpeed:         20,
			runwayHeading:     360,
			expectedHeadwind:  0,
			expectedCrosswind: 20,
		},
		{
			name:              "quartering headwind",
			windDirection:     45,
			windSpeed:         20,
			runwayHeading:     360,
			expectedHeadwind:  14, // cos(45) * 20 ≈ 14
			expectedCrosswind: 14, // sin(45) * 20 ≈ 14
		},
		{
			name:              "quartering tailwind",
			windDirection:     135,
			windSpeed:         20,
			runwayHeading:     360,
			expectedHeadwind:  -14, // cos(135) * 20 ≈ -14
			expectedCrosswind: 14,  // sin(135) * 20 ≈ 14
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weather := Weather{
				WindDirection: tt.windDirection,
				WindSpeedKt:   tt.windSpeed,
			}

			headwind, crosswind := weather.CalculateWindComponents(tt.runwayHeading)

			// Allow 1kt tolerance due to rounding
			if math.Abs(float64(headwind-tt.expectedHeadwind)) > 1 {
				t.Errorf("headwind: expected %d, got %d", tt.expectedHeadwind, headwind)
			}
			if math.Abs(float64(crosswind-tt.expectedCrosswind)) > 1 {
				t.Errorf("crosswind: expected %d, got %d", tt.expectedCrosswind, crosswind)
			}
		})
	}
}

func TestVariableWind(t *testing.T) {
	weather := Weather{
		WindDirection: 0, // Variable
		WindSpeedKt:   5,
		VariableWind:  true,
	}

	t.Run("variable wind reports as crosswind for safety", func(t *testing.T) {
		headwind, crosswind := weather.CalculateWindComponents(360)
		// Variable winds should assume worst case (full crosswind)
		if crosswind != 5 {
			t.Errorf("expected crosswind 5kt for variable wind, got %d", crosswind)
		}
		if headwind != 0 {
			t.Errorf("expected headwind 0 for variable wind, got %d", headwind)
		}
	})
}

func TestCalmWind(t *testing.T) {
	weather := Weather{
		WindDirection: 0,
		WindSpeedKt:   0,
		CalmWind:      true,
	}

	headwind, crosswind := weather.CalculateWindComponents(360)
	if headwind != 0 || crosswind != 0 {
		t.Errorf("expected calm wind to give 0/0, got %d/%d", headwind, crosswind)
	}
}

func TestParseMETAR(t *testing.T) {
	tests := []struct {
		name          string
		metar         string
		windDir       int
		windSpeed     int
		gust          int
		visibility    float64
		ceiling       int
		variableWind  bool
		calmWind      bool
	}{
		{
			name:       "standard METAR",
			metar:      "KJFK 261456Z 09025G35KT 10SM FEW025 BKN050 12/05 A2992",
			windDir:    90,
			windSpeed:  25,
			gust:       35,
			visibility: 10,
			ceiling:    5000,
		},
		{
			name:       "METAR no gusts",
			metar:      "KBOS 261456Z 27015KT 10SM SCT040 15/08 A3001",
			windDir:    270,
			windSpeed:  15,
			gust:       0,
			visibility: 10,
			ceiling:    10000, // No ceiling reported, use high default
		},
		{
			name:         "variable wind",
			metar:        "KSFO 261456Z VRB05KT 10SM CLR 18/10 A2998",
			windDir:      0,
			windSpeed:    5,
			gust:         0,
			visibility:   10,
			ceiling:      10000,
			variableWind: true,
		},
		{
			name:       "calm wind",
			metar:      "KLAX 261456Z 00000KT 10SM CLR 20/12 A2995",
			windDir:    0,
			windSpeed:  0,
			gust:       0,
			visibility: 10,
			ceiling:    10000,
			calmWind:   true,
		},
		{
			name:       "reduced visibility",
			metar:      "KORD 261456Z 18010KT 3SM BR OVC005 08/07 A2988",
			windDir:    180,
			windSpeed:  10,
			gust:       0,
			visibility: 3,
			ceiling:    500,
		},
		{
			name:       "fractional visibility",
			metar:      "KSEA 261456Z 36008KT 1/2SM FG VV002 05/05 A3010",
			windDir:    360,
			windSpeed:  8,
			gust:       0,
			visibility: 0.5,
			ceiling:    200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			weather, err := ParseMETAR(tt.metar)
			if err != nil {
				t.Fatalf("failed to parse METAR: %v", err)
			}

			if weather.WindDirection != tt.windDir {
				t.Errorf("wind direction: expected %d, got %d", tt.windDir, weather.WindDirection)
			}
			if weather.WindSpeedKt != tt.windSpeed {
				t.Errorf("wind speed: expected %d, got %d", tt.windSpeed, weather.WindSpeedKt)
			}
			if weather.GustKt != tt.gust {
				t.Errorf("gust: expected %d, got %d", tt.gust, weather.GustKt)
			}
			if weather.VisibilitySM != tt.visibility {
				t.Errorf("visibility: expected %f, got %f", tt.visibility, weather.VisibilitySM)
			}
			if weather.CeilingFt != tt.ceiling {
				t.Errorf("ceiling: expected %d, got %d", tt.ceiling, weather.CeilingFt)
			}
			if weather.VariableWind != tt.variableWind {
				t.Errorf("variable wind: expected %v, got %v", tt.variableWind, weather.VariableWind)
			}
			if weather.CalmWind != tt.calmWind {
				t.Errorf("calm wind: expected %v, got %v", tt.calmWind, weather.CalmWind)
			}
		})
	}
}
