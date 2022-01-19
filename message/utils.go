package message

// WidgetBanner
// TODO: widget_list的code不定
type WidgetBanner struct {
	Timestamp  int `json:"timestamp"`
	WidgetList struct {
		Field1 struct {
			Id             int      `json:"id"`
			Title          string   `json:"title"`
			Cover          string   `json:"cover"`
			WebCover       string   `json:"web_cover"`
			TipText        string   `json:"tip_text"`
			TipTextColor   string   `json:"tip_text_color"`
			TipBottomColor string   `json:"tip_bottom_color"`
			JumpUrl        string   `json:"jump_url"`
			Url            string   `json:"url"`
			StayTime       int      `json:"stay_time"`
			Site           int      `json:"site"`
			PlatformIn     []string `json:"platform_in"`
			Type           int      `json:"type"`
			BandId         int      `json:"band_id"`
			SubKey         string   `json:"sub_key"`
			SubData        string   `json:"sub_data"`
			IsAdd          bool     `json:"is_add"`
		} `json:"58"`
	} `json:"widget_list"`
}
