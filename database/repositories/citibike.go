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
			citibike_station.latitude,
			citibike_station.longitude,
			citibike_station_statistics_by_day.number_arrivals_previous_90_days,
			citibike_station_statistics_by_day.number_departures_previous_90_days,
			TO_CHAR(citibike_station_statistics_by_day.m_date, 'DAY', 'NLS_DATE_LANGUAGE=''numeric date language'''),
			TO_CHAR(citibike_station_statistics_by_day.m_date, 'MONTH', 'NLS_DATE_LANGUAGE=''numeric date language'''),
			citibike_station_statistics_by_day.number_arrivals
		FROM citibike_station
		JOIN citibike_station_statistics_by_day ON citibike_station.station_id = citibike_station_statistics_by_day.station_id
		WHERE citibike_station_statistics_by_day.m_date >= to_date('29-10-2013', 'dd-mm-yyyy')
		AND citibike_station_statistics_by_day.number_arrivals_previous_90_days > 0
		AND citibike_station_statistics_by_day.number_departures_previous_90_days > 0
		WHERE ROWNUM <= :1`,
		top,
	)
	if err != nil {
		return ds, err
	}

	defer rows.Close()
	for rows.Next() {
		var d models.BatteryDischarge
		err := rows.Scan(
			&d.CurrentLoad,
			&d.CurrentMeasured,
			&d.TemperatureMeasured,
			&d.VoltageLoad,
			&d.VoltageMeasured,
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
