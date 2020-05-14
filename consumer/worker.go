package main

import "alpha/domain"

type Worker struct {
	Measurement string
	probes      chan domain.Probe
}
