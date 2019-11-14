package models

type BatteryDischarge struct {
	CurrentLoad                float64
	CurrentMeasured            float64
	TemperatureMeasured        float64
	VoltageLoad                float64
	VoltageMeasured            float64
	Time                       float64
	Capacity                   float64
	PercentRemainingUsefulLife float64
}
