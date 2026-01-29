package tests

import (
	"testing"
	"github.com/liyanafin/airportsuitability/internal/domain"
)

func TestAircraftLimits(t *testing.T) {
	aircraft := domain.Aircraft{
		TypeCode:        "C172",
		Name:            "Cessna 172",
		Category:        "SEP",
		MinRunwayFt:     1500,
		MaxCrosswindKt:  15,
		MaxTailwindKt:   10,
		VFROnly:         true,
		MinVisibilitySM: 3,
		MinCeilingFt:    1000,
	}

	t.Run("aircraft has valid type code", func(t *testing.T) {
		if aircraft.TypeCode != "C172" {
			t.Errorf("expected type code C172, got %s", aircraft.TypeCode)
		}
	})

	t.Run("aircraft has crosswind limit", func(t *testing.T) {
		if aircraft.MaxCrosswindKt != 15 {
			t.Errorf("expected max crosswind 15kt, got %d", aircraft.MaxCrosswindKt)
		}
	})

	t.Run("aircraft has minimum runway length", func(t *testing.T) {
		if aircraft.MinRunwayFt != 1500 {
			t.Errorf("expected min runway 1500ft, got %d", aircraft.MinRunwayFt)
		}
	})

	t.Run("VFR aircraft has visibility minimums", func(t *testing.T) {
		if !aircraft.VFROnly {
			t.Error("expected aircraft to be VFR only")
		}
		if aircraft.MinVisibilitySM != 3 {
			t.Errorf("expected min visibility 3SM, got %f", aircraft.MinVisibilitySM)
		}
	})
}

