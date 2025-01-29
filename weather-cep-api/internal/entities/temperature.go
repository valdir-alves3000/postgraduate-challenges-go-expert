package entities

type TemperatureResponse struct {
	City      string  `json:"city"`
	Country   string  `json:"country"`
	Localtime string  `json:"localtime"`
	TempC     float64 `json:"temp_c"`
	TempF     float64 `json:"temp_f"`
	TempK     float64 `json:"temp_k"`
}
