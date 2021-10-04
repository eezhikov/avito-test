package service

type UserBalanceRequest struct {
	Id int `json:"id"`
}
type UserBalanceResponse struct {
	Balance float64 `json:"balance"`
}
type UserBalanceAddRequest struct {
	Id      int     `json:"id"`
	Quality float64 `json:"quality"`
}
type UserBalanceAddResponse struct {
	Credited bool `json:"credited"`
}
type UserBalanceTransactionRequest struct {
	UserIdFrom int     `json:"user_id_from"`
	UserIdTo   int     `json:"user_id_to"`
	Quality    float64 `json:"quality"`
}
type UserBalanceTransactionResponse struct {
	Transaction bool `json:"transaction"`
}

type Valute struct {
	ValuteMap map[string]CurrencyObject `json:"Valute"`
}

type CurrencyObject struct {
	Value float64 `json:"Value"`
}
