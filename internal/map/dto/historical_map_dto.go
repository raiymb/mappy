package dto

type HistoricalMapResponse struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Period   string   `json:"period"`
	ImageURL string   `json:"imageUrl"`
	Tags     []string `json:"tags"`
}
