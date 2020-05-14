package domain

import (
	"time"
)

type Measurement struct {
	Probes []Probe `json:"probes"`
}

type Probe struct {
	Sensor string    `json:"sensor"`
	Unit   string    `json:"unit"`
	Date   time.Time `json:"date"`
	Value  float32   `json:"value"`
}
