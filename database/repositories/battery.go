package repositories

import (
	"gpu-demonstration-api/database"
	"gpu-demonstration-api/models"
)

type BatteryRepo struct{}

func (br *BatteryRepo) SelectSampleDischargeTop(top int) ([]models.BatteryDischarge, error) {
	var ds []models.BatteryDischarge
	rows, err := database.DB.Query(
		`SELECT 
			battery_discharge_data.current_load,
			battery_discharge_data.current_measured,
			battery_discharge_data.temperature_measured,
			battery_discharge_data.voltage_load,
			battery_discharge_data.voltage_measured,
			battery_discharge_data.m_time,
			battery_discharge_data.m_capacity,
			battery_cycle.pct_rul
		FROM battery_battery 
		LEFT JOIN battery_cycle ON battery_battery.battery_id = battery_cycle.battery_id
		LEFT JOIN battery_discharge_data ON battery_discharge_data.cycle_id = battery_cycle.cycle_id
		WHERE ROWNUM <= $1`,
		top,
	)

	defer rows.Close()
	for rows.Next() {
		var d models.BatteryDischarge
		err := rows.Scan(
			&d.CurrentLoad,
			&d.CurrentMeasured,
			&d.TemperatureMeasured,
			&d.VoltageLoad,
			&d.Time,
			&d.Capacity,
			&d.PercentRemainingUsefulLife,
		)
		if err != nil {
			return ds, err
		}
		ds = append(ds, d)
	}
	err = rows.Err()
	if err != nil {
		return ds, err
	}
	return ds, nil
}
