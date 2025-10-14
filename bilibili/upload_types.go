package bilibili

// Video 视频文件信息
type Video struct {
	Title    string `json:"title"`
	Filename string `json:"filename"`
	Desc     string `json:"desc"`
}

// Studio 投稿信息
type Studio struct {
	Copyright     int      `json:"copyright"`            // 是否转载, 1-自制 2-转载
	Source        string   `json:"source"`               // 转载来源
	Tid           int      `json:"tid"`                  // 投稿分区
	Cover         string   `json:"cover"`                // 视频封面
	Title         string   `json:"title"`                // 视频标题
	DescFormatId  int      `json:"desc_format_id"`       // 简介格式ID
	Desc          string   `json:"desc"`                 // 视频简介
	Dynamic       string   `json:"dynamic"`              // 空间动态
	Subtitle      Subtitle `json:"subtitle"`             // 字幕信息
	Tag           string   `json:"tag"`                  // 视频标签
	Videos        []Video  `json:"videos"`               // 视频文件列表
	Dtime         *int64   `json:"dtime,omitempty"`      // 延时发布时间
	OpenSubtitle  bool     `json:"open_subtitle"`        // 是否开启字幕
	Interactive   int      `json:"interactive"`          // 是否开启互动
	MissionId     *int     `json:"mission_id,omitempty"` // 任务ID
	Dolby         int      `json:"dolby"`                // 是否开启杜比音效
	LosslessMusic int      `json:"lossless_music"`       // 是否开启Hi-Res
	NoReprint     int      `json:"no_reprint"`           // 是否禁止转载
	OpenElec      int      `json:"open_elec"`            // 是否开启充电
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