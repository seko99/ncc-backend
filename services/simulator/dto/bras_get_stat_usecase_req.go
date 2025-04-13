package dto

type BrasGetStatUsecaseRequest struct{}

type BrasGetStatUsecaseResponse struct {
	Sessions  int `json:"sessions"`
	Leases    int `json:"leases"`
	Customers int `json:"customers"`
	Nases     int `json:"nases"`
}
