package repositories

import (
	"gpu-demonstration-api/database"
	"gpu-demonstration-api/models"
)

type CitibikeRepo struct{}

func (cr *CitibikeRepo) SelectSampleStationStatisticsByDay(top int) ([]models.CitibikeStationStatisticsByDay, error) {
	var cs []models.CitibikeStationStatisticsByDay
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
		AND citibike_station.latitude != 0
		AND citibike_station.longitude != 0
		AND ROWNUM <= :1`,
		top,
	)
	if err != nil {
		return cs, err
	}

	defer rows.Close()
	for rows.Next() {
		var c models.CitibikeStationStatisticsByDay
		err := rows.Scan(
			&c.Latitude,
			&c.Longitude,
			&c.NumberArrivalsPrevious90Days,
			&c.NumberDeparturesPrevious90Days,
			&c.Day,
			&c.Month,
			&c.NumberArrivals,
		)
		if err != nil {
			return cs, err
		}
		cs = append(cs, c)
	}
	err = rows.Err()
	if err != nil {
		return cs, err
	}
	return cs, nil
}
