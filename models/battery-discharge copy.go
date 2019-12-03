package models

type CitibikeStationStatisticsByDay struct {
	Latitude                       float64
	Longitude                      float64
	NumberArrivalsPrevious90Days   int
	NumberDeparturesPrevious90Days int
	Day                            float64
	Month                          float64
	NumberArrivals                 int
}
