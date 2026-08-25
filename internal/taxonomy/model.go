package taxonomy

type Location struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Specialty struct {
	ID         int64  `json:"id"`
	LocationID int64  `json:"location_id"`
	Name       string `json:"name"`
	URL        string `json:"url"`
}
