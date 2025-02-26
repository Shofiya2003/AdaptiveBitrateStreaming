package data

type VideoEvent struct {
	VideoURL string `json:"video_url"`
	Bucket   string `json:"bucket"`
	Key      string `json:"key"`
}
