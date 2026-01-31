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


