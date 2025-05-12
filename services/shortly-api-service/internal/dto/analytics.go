package dto

type AnalyticsResponse struct {
	IPAddress string `json:"ipAddress"`
	OS        string `json:"os"`
	Device    string `json:"device"`
	Browser   string `json:"browser"`
	UserAgent string `json:"userAgent"`
	ClickedAt string `json:"clickedAt"`
	Referrer  string `json:"referrer"`
	Country   string `json:"country"`
}
