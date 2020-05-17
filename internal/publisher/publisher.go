package publisher

import (
	"alpha/domain"
	influx "github.com/influxdata/influxdb-client-go"
	log "github.com/sirupsen/logrus"
	"time"
)

type Publisher interface {
	Start() error
}

type PublisherImpl struct {
	InfluxHost  string
	InfluxToken string
	Probes      chan domain.Probe
	influx      influx.InfluxDBClient
	topics      *map[string]bool
}

func (c PublisherImpl) publishToInflux(p domain.Probe) {
	if c.influx == nil {
		c.influx = influx.NewClient(c.InfluxHost, c.InfluxToken)
	}
	writeApi := c.influx.WriteApi("", "probes")
	point := influx.NewPoint(p.Sensor,
		map[string]string{"unit": "delay"},
		map[string]interface{}{p.Unit: p.Value},
		time.Now())
	writeApi.WritePoint(point)

	log.Info("Writing to influx")
}

func (c PublisherImpl) Start() {
	go func() {
		for {
			select {
			case p := <-c.Probes:
				c.publishToInflux(p)
			}
		}
	}()
}
