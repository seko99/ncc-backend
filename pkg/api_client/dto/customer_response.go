package dto

type CustomerResponseDTO struct {
	ID      string `json:"id"`
	Login   string `json:"login"`
	Deposit int    `json:"deposit"`
	Credit  int    `json:"credit"`
	Name    string `json:"name"`
}
