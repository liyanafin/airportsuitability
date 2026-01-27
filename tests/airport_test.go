package tests

import (
	"testing"

	"github.com/liyanafin/airportsuitability/internal/domain"
)

func TestAirport_Basics(t *testing.T) {
	airport := domain.Airport{
		ICAO:        "KJFK",
		Name:        "John F Kennedy Intl",
		ElevationFt: 13,
		Runways: []domain.Runway{
			{ID: "04L", Heading: 40, LengthFt: 8400, WidthFt: 150, Surface: "ASPH"},
			{ID: "22R", Heading: 220, LengthFt: 8400, WidthFt: 150, Surface: "ASPH"},
			{ID: "13L", Heading: 130, LengthFt: 10000, WidthFt: 150, Surface: "ASPH"},
		},
	}

	t.Run("airport has ICAO code", func(t *testing.T) {
		if airport.ICAO != "KJFK" {
			t.Errorf("expected ICAO KJFK, got %s", airport.ICAO)
		}
	})

	t.Run("airport has name", func(t *testing.T) {
		if airport.Name != "John F Kennedy Intl" {
			t.Errorf("expected name John F Kennedy Intl, got %s", airport.Name)
		}
	})

	t.Run("airport has elevation", func(t *testing.T) {
		if airport.ElevationFt != 13 {
			t.Errorf("expected elevation 13ft, got %d", airport.ElevationFt)
		}
	})

	t.Run("airport has runways", func(t *testing.T) {
		if len(airport.Runways) != 3 {
			t.Errorf("expected 3 runways, got %d", len(airport.Runways))
		}
	})
}
