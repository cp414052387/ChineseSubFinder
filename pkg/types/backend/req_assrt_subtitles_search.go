package backend

type ReqAssrtSubtitlesSearch struct {
	VideoFPath string `json:"video_f_path"`
	IsMovie    bool   `json:"is_movie"`
}
