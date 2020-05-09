package domain

import (
	"time"
)

type Probes struct {
	Unit        string  `json:"unit"`
	Measurement string  `json:"string"`
	Probes      []Probe `json:"probes"`
}
type Probe struct {
	Date  time.Time   `json:"date"`
	Value interface{} `json:"value"`
}
