package model

type ViedeoDetails struct {
	VideoName string  `json:"videoName"`
	Viewes    float64 `json:"viewes"`
}

type Video struct {
	Name string `json:"name"`
}

type AddVideoStatus struct {
	Status string
}
