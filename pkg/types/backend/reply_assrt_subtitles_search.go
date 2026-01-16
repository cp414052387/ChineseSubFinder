package backend

type ReplyAssrtSubtitle struct {
	SubSha256 string `json:"sub_sha256"`
	Title     string `json:"title"`
	Language  int    `json:"language"`
	Ext       string `json:"ext"`
	Season    int    `json:"season"`
	Episode   int    `json:"episode"`
}

type ReplyAssrtSubtitlesSearch struct {
	Subtitles []ReplyAssrtSubtitle `json:"subtitles"`
}
