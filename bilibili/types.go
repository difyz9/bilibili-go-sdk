package bilibili

import (
	"fmt"
	"strings"
)

// ResponseData 通用API响应结构
type ResponseData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	TTL     int         `json:"ttl,omitempty"`
}

// QRCodeData 二维码数据
type QRCodeData struct {
	URL      string `json:"url"`
	AuthCode string `json:"auth_code"`
}

// QRCodeResponse 二维码响应
type QRCodeResponse struct {
	Code int        `json:"code"`
	Data QRCodeData `json:"data"`
}

// LoginInfo 登录信息
type LoginInfo struct {
	CookieInfo map[string]interface{} `json:"cookie_info"`
	SSO        []string               `json:"sso"`
	TokenInfo  TokenInfo              `json:"token_info"`
	Platform   string                 `json:"platform,omitempty"`
}

// LoginResponse 登录响应 (兼容短信登录和密码登录)
type LoginResponse struct {
	Code       int                    `json:"code"`
	Message    string                 `json:"message"`
	Data       map[string]interface{} `json:"data"`
	CookieInfo CookieInfo             `json:"cookie_info"`
	TokenInfo  TokenInfo              `json:"token_info"`
}

// CookieInfo Cookie信息
type CookieInfo struct {
	Cookies string `json:"cookies"`
}

// GetCookieString 获取Cookie字符串
func (li *LoginInfo) GetCookieString() string {
	cookies, ok := li.CookieInfo["cookies"].([]interface{})
	if !ok {
		return ""
	}

	var cookieStrs []string
	for _, cookie := range cookies {
		if cookieMap, ok := cookie.(map[string]interface{}); ok {
			if name, nameOk := cookieMap["name"].(string); nameOk {
				if value, valueOk := cookieMap["value"].(string); valueOk {
					cookieStrs = append(cookieStrs, fmt.Sprintf("%s=%s", name, value))
				}
			}
		}
	}

	return strings.Join(cookieStrs, "; ")
}

// GetCSRFToken 获取CSRF token
func (li *LoginInfo) GetCSRFToken() (string, error) {
	cookies, ok := li.CookieInfo["cookies"].([]interface{})
	if !ok {
		return "", fmt.Errorf("no cookies found")
	}

	for _, cookie := range cookies {
		if cookieMap, ok := cookie.(map[string]interface{}); ok {
			if name, nameOk := cookieMap["name"].(string); nameOk && name == "bili_jct" {
				if value, valueOk := cookieMap["value"].(string); valueOk {
					return value, nil
				}
			}
		}
	}

	return "", fmt.Errorf("CSRF token not found in cookie")
}

// TokenInfo 令牌信息
type TokenInfo struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	Mid          int64  `json:"mid"`
	RefreshToken string `json:"refresh_token"`
	Uname        string `json:"uname,omitempty"`
	Face         string `json:"face,omitempty"`
}

// UserBasicInfo 用户基本信息
type UserBasicInfo struct {
	Mid             int64                  `json:"mid"`
	Name            string                 `json:"name"`
	Sex             string                 `json:"sex"`
	Face            string                 `json:"face"`
	FaceNFT         int                    `json:"face_nft,omitempty"`
	FaceNFTType     int                    `json:"face_nft_type,omitempty"`
	Sign            string                 `json:"sign"`
	Rank            int                    `json:"rank"`
	Level           int                    `json:"level"`
	Jointime        int64                  `json:"jointime,omitempty"`
	Moral           int                    `json:"moral,omitempty"`
	Silence         int                    `json:"silence,omitempty"`
	Coins           int                    `json:"coins,omitempty"`
	FansBadge       bool                   `json:"fans_badge,omitempty"`
	FansMedal       *UserFansMedal         `json:"fans_medal,omitempty"`
	Official        *UserOfficialInfo      `json:"official,omitempty"`
	VIP             *UserVIPInfo           `json:"vip,omitempty"`
	Pendant         *UserPendantInfo       `json:"pendant,omitempty"`
	Nameplate       *UserNameplateInfo     `json:"nameplate,omitempty"`
	UserHonourInfo  *UserHonourInfo        `json:"user_honour_info,omitempty"`
	IsFollowed      bool                   `json:"is_followed,omitempty"`
	TopPhoto        string                 `json:"top_photo,omitempty"`
	Theme           map[string]interface{} `json:"theme,omitempty"`
	SysNotice       map[string]interface{} `json:"sys_notice,omitempty"`
	LiveRoom        *UserLiveRoom          `json:"live_room,omitempty"`
	Birthday        string                 `json:"birthday"`
	School          *UserSchoolInfo        `json:"school,omitempty"`
	Profession      *UserProfessionInfo    `json:"profession,omitempty"`
	Tags            []string               `json:"tags,omitempty"`
	Series          map[string]interface{} `json:"series,omitempty"`
	IsSeniorMember  int                    `json:"is_senior_member,omitempty"`
	MCNInfo         interface{}            `json:"mcn_info,omitempty"`
	GaiaResType     int                    `json:"gaia_res_type,omitempty"`
	GaiaData        interface{}            `json:"gaia_data,omitempty"`
	IsRisk          bool                   `json:"is_risk,omitempty"`
	Elec            map[string]interface{} `json:"elec,omitempty"`
	Contract        *UserContractInfo      `json:"contract,omitempty"`
	CertificateShow bool                   `json:"certificate_show,omitempty"`
	NameRender      map[string]interface{} `json:"name_render,omitempty"`
}

// UserOfficialInfo 用户认证信息
type UserOfficialInfo struct {
	Role  int    `json:"role"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
	Type  int    `json:"type"`
}

// UserVIPInfo 用户会员信息
type UserVIPInfo struct {
	Type               int                    `json:"type"`
	Status             int                    `json:"status"`
	DueDate            int64                  `json:"due_date"`
	VipPayType         int                    `json:"vip_pay_type"`
	ThemeType          int                    `json:"theme_type"`
	Label              *UserVIPLabel          `json:"label,omitempty"`
	AvatarSubscript    int                    `json:"avatar_subscript"`
	NicknameColor      string                 `json:"nickname_color"`
	Role               int                    `json:"role"`
	AvatarSubscriptURL string                 `json:"avatar_subscript_url"`
	TVVipStatus        int                    `json:"tv_vip_status"`
	TVVipPayType       int                    `json:"tv_vip_pay_type"`
	TVDueDate          int64                  `json:"tv_due_date"`
	AvatarIcon         map[string]interface{} `json:"avatar_icon,omitempty"`
}

// UserVIPLabel 用户会员标签
type UserVIPLabel struct {
	Path               string `json:"path"`
	Text               string `json:"text"`
	LabelTheme         string `json:"label_theme"`
	TextColor          string `json:"text_color"`
	BgStyle            int    `json:"bg_style"`
	BgColor            string `json:"bg_color"`
	BorderColor        string `json:"border_color"`
	UseImgLabel        bool   `json:"use_img_label"`
	ImgLabelURIHans    string `json:"img_label_uri_hans"`
	ImgLabelURIHant    string `json:"img_label_uri_hant"`
	ImgLabelHansStatic string `json:"img_label_uri_hans_static"`
	ImgLabelHantStatic string `json:"img_label_uri_hant_static"`
}

// UserPendantInfo 用户头像框信息
type UserPendantInfo struct {
	PID               int64  `json:"pid"`
	Name              string `json:"name"`
	Image             string `json:"image"`
	Expire            int64  `json:"expire"`
	ImageEnhance      string `json:"image_enhance"`
	ImageEnhanceFrame string `json:"image_enhance_frame"`
	NPID              int64  `json:"n_pid"`
}

// UserNameplateInfo 用户勋章信息
type UserNameplateInfo struct {
	NID        int64  `json:"nid"`
	Name       string `json:"name"`
	Image      string `json:"image"`
	ImageSmall string `json:"image_small"`
	Level      string `json:"level"`
	Condition  string `json:"condition"`
}

// UserHonourInfo 用户荣誉信息
type UserHonourInfo struct {
	Mid    int64         `json:"mid"`
	Colour interface{}   `json:"colour"`
	Tags   []interface{} `json:"tags"`
}

// UserFansMedal 用户粉丝勋章信息
type UserFansMedal struct {
	Show  bool           `json:"show"`
	Wear  bool           `json:"wear"`
	Medal *UserMedalInfo `json:"medal,omitempty"`
}

// UserMedalInfo 粉丝勋章详情
type UserMedalInfo struct {
	UID              int64  `json:"uid"`
	TargetID         int64  `json:"target_id"`
	MedalID          int64  `json:"medal_id"`
	Level            int    `json:"level"`
	MedalName        string `json:"medal_name"`
	MedalColor       int    `json:"medal_color"`
	Intimacy         int    `json:"intimacy"`
	NextIntimacy     int    `json:"next_intimacy"`
	DayLimit         int    `json:"day_limit"`
	TodayFeed        int    `json:"today_feed"`
	MedalColorStart  int    `json:"medal_color_start"`
	MedalColorEnd    int    `json:"medal_color_end"`
	MedalColorBorder int    `json:"medal_color_border"`
	IsLighted        int    `json:"is_lighted"`
	LightStatus      int    `json:"light_status"`
	WearingStatus    int    `json:"wearing_status"`
	Score            int    `json:"score"`
}

// UserLiveRoom 用户直播间信息
type UserLiveRoom struct {
	RoomStatus    int                    `json:"roomStatus"`
	LiveStatus    int                    `json:"liveStatus"`
	URL           string                 `json:"url"`
	Title         string                 `json:"title"`
	Cover         string                 `json:"cover"`
	WatchedShow   map[string]interface{} `json:"watched_show,omitempty"`
	RoomID        int64                  `json:"roomid"`
	RoundStatus   int                    `json:"roundStatus"`
	BroadcastType int                    `json:"broadcast_type"`
}

// UserSchoolInfo 用户学校信息
type UserSchoolInfo struct {
	Name string `json:"name"`
}

// UserProfessionInfo 用户专业资质信息
type UserProfessionInfo struct {
	Name       string `json:"name"`
	Department string `json:"department"`
	Title      string `json:"title"`
	IsShow     int    `json:"is_show"`
}

// UserContractInfo 用户老粉计划信息
type UserContractInfo struct {
	IsDisplay       bool `json:"is_display"`
	IsFollowDisplay bool `json:"is_follow_display"`
}

// MyInfoResponse 详细用户信息响应结构 (myinfo API)
type MyInfoResponse struct {
	Mid       int64       `json:"mid"`
	Name      string      `json:"name"` // API返回的是name，不是uname
	Uname     string      `json:"-"`    // 为了兼容性保留，从Name复制
	UserID    string      `json:"-"`    // API中没有这个字段
	Sign      string      `json:"sign"`
	Birthday  interface{} `json:"birthday"` // 可能是字符串或数字
	Sex       string      `json:"sex"`
	NickFree  bool        `json:"-"`    // API中没有这个字段
	Rank      interface{} `json:"rank"` // 可能是字符串或数字
	Face      string      `json:"face"`
	Level     int         `json:"level"`
	Silence   int         `json:"silence"`
	Coins     interface{} `json:"coins"`     // 可能是整数或浮点数
	Follower  int         `json:"follower"`  // 粉丝数
	Following int         `json:"following"` // 关注数
	// 为了兼容性保留旧字段名
	Fans      int `json:"-"`
	Attention int `json:"-"`
	Friend    int `json:"-"`
}

// PostProcess 后处理方法，设置兼容性字段
func (m *MyInfoResponse) PostProcess() {
	m.Uname = m.Name
	m.Fans = m.Follower
	m.Attention = m.Following
}

// GetBirthdayString 获取生日字符串
func (m *MyInfoResponse) GetBirthdayString() string {
	switch v := m.Birthday.(type) {
	case string:
		return v
	case float64:
		if v == 0 {
			return ""
		}
		return fmt.Sprintf("%.0f", v)
	case int:
		if v == 0 {
			return ""
		}
		return fmt.Sprintf("%d", v)
	case int64:
		if v == 0 {
			return ""
		}
		return fmt.Sprintf("%d", v)
	default:
		return ""
	}
}

// GetCoins 获取硬币数量
func (m *MyInfoResponse) GetCoins() int {
	switch v := m.Coins.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case int64:
		return int(v)
	default:
		return 0
	}
}

// GetRankString 获取排名字符串
func (m *MyInfoResponse) GetRankString() string {
	switch v := m.Rank.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%.0f", v)
	case int:
		return fmt.Sprintf("%d", v)
	case int64:
		return fmt.Sprintf("%d", v)
	default:
		return ""
	}
}

// PartitionType 分区类型
type PartitionType struct {
	ID       int             `json:"id"`
	Name     string          `json:"name"`
	Desc     string          `json:"desc,omitempty"`
	Children []PartitionType `json:"children,omitempty"`
}

// ArchivePreData archive_pre接口返回的数据
type ArchivePreData struct {
	TypeList []PartitionType `json:"typelist"`
}
