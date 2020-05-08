package messages

type data struct {
	Interval int `json:"interval"`
	Value    int `json:"value"`
}

type Probe struct {
	SensorId int  `json:"sensorId"`
	Data     data `json:"data"`
}
