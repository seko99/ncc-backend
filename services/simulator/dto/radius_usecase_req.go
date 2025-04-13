package dto

type RadiusUsecaseRequest struct {
	Limit         int
	Secret        string
	NasIP         string
	NasIdentifier string
	Auth          string
	Acct          string
}

type RadiusUsecaseResponse struct {
	Leases             int
	Sent               int
	Accepted           int
	Rejected           int
	PersistentSessions int
	CacheSessions      int
}

type RadiusKillSessionsUsecaseRequest struct {
	Sessions int  `json:"sessions"`
	Random   bool `json:"random"`
}

type RadiusKillSessionsUsecaseResponse struct {
	Killed int `json:"killed"`
}
