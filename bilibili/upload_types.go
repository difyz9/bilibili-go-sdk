package bilibili

// Video 视频文件信息
type Video struct {
	Title    string `json:"title"`
	Filename string `json:"filename"`
	Desc     string `json:"desc"`
	CID      int64  `json:"cid,omitempty"`
}

// DescV2Item 简介富文本片段
type DescV2Item struct {
	BizID   string `json:"biz_id"`
	RawText string `json:"raw_text"`
	Type    int    `json:"type"`
}

// Studio 投稿信息
type Studio struct {
	Copyright        int          `json:"copyright"`              // 是否转载, 1-自制 2-转载
	Source           string       `json:"source,omitempty"`       // 转载来源
	Tid              int          `json:"tid"`                    // 投稿分区
	HumanType2       int          `json:"human_type2,omitempty"`
	Cover            string       `json:"cover,omitempty"`        // 视频封面
	Cover43          string       `json:"cover43,omitempty"`
	Title            string       `json:"title"`                  // 视频标题
	DescFormatId     int          `json:"desc_format_id"`         // 简介格式ID
	Desc             string       `json:"desc,omitempty"`         // 视频简介
	DescV2           []DescV2Item `json:"desc_v2,omitempty"`
	Dynamic          string       `json:"dynamic"`                // 空间动态
	Subtitle         Subtitle     `json:"subtitle"`               // 字幕信息
	Tag              string       `json:"tag"`                    // 视频标签
	Videos           []Video      `json:"videos"`                 // 视频文件列表
	Dtime            *int64       `json:"dtime,omitempty"`        // 延时发布时间
	OpenSubtitle     bool         `json:"open_subtitle,omitempty"`
	Interactive      int          `json:"interactive"`            // 是否开启互动
	ActReserveCreate int          `json:"act_reserve_create,omitempty"`
	MissionId        *int         `json:"mission_id,omitempty"` // 任务ID
	TopicID          *int         `json:"topic_id,omitempty"`
	Recreate         int          `json:"recreate,omitempty"`
	Dolby            int          `json:"dolby"`                  // 是否开启杜比音效
	LosslessMusic    int          `json:"lossless_music"`         // 是否开启Hi-Res
	NoDisturbance    int          `json:"no_disturbance,omitempty"`
	NoReprint        int          `json:"no_reprint"`             // 是否禁止转载
	UpSelectionReply bool         `json:"up_selection_reply"`
	UpCloseReply     bool         `json:"up_close_reply"`
	UpCloseDanmu     bool         `json:"up_close_danmu"`
	WebOS            int          `json:"web_os,omitempty"`
	IsOnlySelf       *int         `json:"is_only_self,omitempty"`
	Is360            int          `json:"is_360,omitempty"`
	NeutralMark      string       `json:"neutral_mark,omitempty"`
	OpenElec         int          `json:"open_elec"` // 是否开启充电
}

// Subtitle 字幕信息
type Subtitle struct {
	Open int    `json:"open"`
	Lan  string `json:"lan"`
}

// PreUploadInfo 预上传信息 - 直接对应B站API响应
type PreUploadInfo struct {
	OK              int         `json:"OK"`
	Auth            string      `json:"auth"`
	BizId           int64       `json:"biz_id"`
	ChunkRetry      int         `json:"chunk_retry"`
	ChunkRetryDelay int         `json:"chunk_retry_delay"`
	ChunkSize       int         `json:"chunk_size"`
	Endpoint        string      `json:"endpoint"`
	Endpoints       []string    `json:"endpoints"`
	ExposeParams    interface{} `json:"expose_params"`
	PutQuery        string      `json:"put_query"`
	Threads         int         `json:"threads"`
	Timeout         int         `json:"timeout"`
	Uip             string      `json:"uip"`
	UposUri         string      `json:"upos_uri"`
}