package models

type DashboardData struct {
	RoomingHouseName string                         `json:"roomingHouseName"`
	TransactionData  []TransactionDashboardResponse `json:"transactionData"`
}
