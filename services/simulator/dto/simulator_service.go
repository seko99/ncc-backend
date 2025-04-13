package dto

import "time"

type LatLng struct {
	Lat float64
	Lng float64
}

type FakerRequestDTO struct {
	InitDictionaries bool
	FakeData         bool
	Cleanup          bool
	Disconnect       bool
	UpdateMap        bool
	UpdateSessions   bool
	UpdateLeases     bool

	CustomersNumber int
	MaxMapNodes     int
	MinBuildLevel   int
	LeftUpper       LatLng
	RightBottom     LatLng
	FirstPayment    time.Time
}
