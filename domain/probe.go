package domain

import (
	"time"
)

type Measurement struct {
	Unit        string  `json:"unit"`
	Measurement string  `json:"measurement"`
	Probes      []Probe `json:"probes"`
}
type Probe struct {
	Date   time.Time              `json:"date"`
	Values map[string]interface{} `json:"values"`
}
