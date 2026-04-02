package bilibili

import "strings"

// VideoZone 视频分区信息。
type VideoZone struct {
	ID           int         `json:"id"`
	ParentID     int         `json:"parent_id,omitempty"`
	Name         string      `json:"name"`
	DisplayName  string      `json:"display_name,omitempty"`
	Slug         string      `json:"slug"`
	Description  string      `json:"description,omitempty"`
	URI          string      `json:"uri"`
	Note         string      `json:"note,omitempty"`
	IsMain       bool        `json:"is_main,omitempty"`
	IsDeprecated bool        `json:"is_deprecated,omitempty"`
	IsRedirect   bool        `json:"is_redirect,omitempty"`
	Children     []VideoZone `json:"children,omitempty"`
}

var videoZones = []VideoZone{
	{
		ID:          1,
		Name:        "动画",
		DisplayName: "动画(主分区)",
		Slug:        "douga",
		URI:         "/v/douga",
		IsMain:      true,
		Children: []VideoZone{
			{
				ID:          24,
				ParentID:    1,
				Name:        "MAD·AMV",
				DisplayName: "MAD·AMV",
				Slug:        "mad",
				Description: "具有一定制作程度的动画或静画的二次创作视频",
				URI:         "/v/douga/mad",
			},
			{
				ID:          25,
				ParentID:    1,
				Name:        "MMD·3D",
				DisplayName: "MMD·3D",
				Slug:        "mmd",
				Description: "使用MMD（MikuMikuDance）和其他3D建模类软件制作的视频",
				URI:         "/v/douga/mmd",
			},
			{
				ID:          47,
				ParentID:    1,
				Name:        "同人·手书",
				DisplayName: "同人·手书 (原短片·手书)",
				Slug:        "handdrawn",
				Description: "追求个人特色和创意表达的手书（绘）、以及同人作品展示、宣传为主的内容",
				URI:         "/v/douga/handdrawn",
				Note:        "原短片·手书",
			},
			{
				ID:          257,
				ParentID:    1,
				Name:        "配音",
				DisplayName: "配音",
				Slug:        "voice",
				Description: "使用ACGN相关画面或台本素材进行人工配音创作的内容",
				URI:         "/v/douga/voice",
			},
			{
				ID:          210,
				ParentID:    1,
				Name:        "手办·模玩",
				DisplayName: "手办·模玩",
				Slug:        "garage_kit",
				Description: "手办模玩的测评、改造或其他衍生内容",
				URI:         "/v/douga/garage_kit",
			},
			{
				ID:          86,
				ParentID:    1,
				Name:        "特摄",
				DisplayName: "特摄",
				Slug:        "tokusatsu",
				Description: "特摄相关衍生视频",
				URI:         "/v/douga/tokusatsu",
			},
			{
				ID:          253,
				ParentID:    1,
				Name:        "动漫杂谈",
				DisplayName: "动漫杂谈",
				Slug:        "acgntalks",
				Description: "以谈话形式对ACGN文化圈进行的鉴赏、吐槽、评点、解说、推荐、科普等内容",
				URI:         "/v/douga/acgntalks",
			},
			{
				ID:          27,
				ParentID:    1,
				Name:        "综合",
				DisplayName: "综合",
				Slug:        "other",
				Description: "以动画及动画相关内容为素材，包括但不仅限于音频替换、恶搞改编、排行榜等内容",
				URI:         "/v/douga/other",
			},
		},
	},
	{
		ID:          13,
		Name:        "番剧",
		DisplayName: "番剧(主分区)",
		Slug:        "anime",
		URI:         "/anime",
		IsMain:      true,
		Children: []VideoZone{
			{
				ID:          51,
				ParentID:    13,
				Name:        "资讯",
				DisplayName: "资讯",
				Slug:        "information",
				Description: "动画番剧相关资讯视频",
				URI:         "/v/anime/information",
			},
			{
				ID:          152,
				ParentID:    13,
				Name:        "官方延伸",
				DisplayName: "官方延伸",
				Slug:        "offical",
				Description: "动画番剧为主题的宣传节目、采访视频，及声优相关视频",
				URI:         "/v/anime/offical",
			},
			{
				ID:          32,
				ParentID:    13,
				Name:        "完结动画",
				DisplayName: "完结动画",
				Slug:        "finish",
				Description: "已完结的动画番剧合集",
				URI:         "/v/anime/finish",
			},
			{
				ID:          33,
				ParentID:    13,
				Name:        "连载动画",
				DisplayName: "连载动画",
				Slug:        "serial",
				Description: "当季连载的动画番剧",
				URI:         "/v/anime/serial",
			},
		},
	},
	{
		ID:          167,
		Name:        "国创",
		DisplayName: "国创(主分区)",
		Slug:        "guochuang",
		URI:         "/guochuang",
		IsMain:      true,
		Children: []VideoZone{
			{
				ID:          153,
				ParentID:    167,
				Name:        "国产动画",
				DisplayName: "国产动画",
				Slug:        "chinese",
				Description: "我国出品的PGC动画",
				URI:         "/v/guochuang/chinese",
			},
			{
				ID:          168,
				ParentID:    167,
				Name:        "国产原创相关",
				DisplayName: "国产原创相关",
				Slug:        "original",
				URI:         "/v/guochuang/original",
			},
			{
				ID:          169,
				ParentID:    167,
				Name:        "布袋戏",
				DisplayName: "布袋戏",
				Slug:        "puppetry",
				URI:         "/v/guochuang/puppetry",
			},
			{
				ID:          170,
				ParentID:    167,
				Name:        "资讯",
				DisplayName: "资讯",
				Slug:        "information",
				URI:         "/v/guochuang/information",
			},
			{
				ID:          195,
				ParentID:    167,
				Name:        "动态漫·广播剧",
				DisplayName: "动态漫·广播剧",
				Slug:        "motioncomic",
				URI:         "/v/guochuang/motioncomic",
			},
		},
	},
	{
		ID:          3,
		Name:        "音乐",
		DisplayName: "音乐(主分区)",
		Slug:        "music",
		URI:         "/v/music",
		IsMain:      true,
		Children: []VideoZone{
			{
				ID:          28,
				ParentID:    3,
				Name:        "原创音乐",
				DisplayName: "原创音乐",
				Slug:        "original",
				Description: "原创歌曲及纯音乐，包括改编、重编曲及remix",
				URI:         "/v/music/original",
			},
			{
				ID:          29,
				ParentID:    3,
				Name:        "音乐现场",
				DisplayName: "音乐现场",
				Slug:        "live",
				Description: "音乐表演的实况视频，包括官方的综艺节目、音乐剧、音乐节、演唱会、打歌舞台现场等，以及个人演出/街头表演现场等",
				URI:         "/v/music/live",
			},
			{
				ID:          31,
				ParentID:    3,
				Name:        "翻唱",
				DisplayName: "翻唱",
				Slug:        "cover",
				Description: "对曲目的人声再演绎视频",
				URI:         "/v/music/cover",
			},
			{
				ID:          59,
				ParentID:    3,
				Name:        "演奏",
				DisplayName: "演奏",
				Slug:        "perform",
				Description: "乐器和非传统乐器器材的演奏作品",
				URI:         "/v/music/perform",
			},
			{
				ID:          243,
				ParentID:    3,
				Name:        "乐评盘点",
				DisplayName: "乐评盘点",
				Slug:        "commentary",
				Description: "音乐类新闻、盘点、点评、reaction、榜单、采访、幕后故事、唱片开箱等",
				URI:         "/v/music/commentary",
			},
			{
				ID:          30,
				ParentID:    3,
				Name:        "VOCALOID·UTAU",
				DisplayName: "VOCALOID·UTAU",
				Slug:        "vocaloid",
				Description: "以VOCALOID等歌声合成引擎为基础，运用各类音源进行的创作",
				URI:         "/v/music/vocaloid",
			},
			{
				ID:          193,
				ParentID:    3,
				Name:        "MV",
				DisplayName: "MV",
				Slug:        "mv",
				Description: "为音乐作品配合拍摄或制作的音乐录影带（Music Video），以及自制拍摄、剪辑、翻拍MV",
				URI:         "/v/music/mv",
			},
			{
				ID:          266,
				ParentID:    3,
				Name:        "音乐粉丝饭拍",
				DisplayName: "音乐粉丝饭拍",
				Slug:        "fan_videos",
				Description: "在音乐演出现场由粉丝团体或个人拍摄的非官方记录视频，包括但不限于粉丝自制饭拍、直拍、Vlog以及衍生的内容混剪等",
				URI:         "/v/music/fan_videos",
			},
			{
				ID:          265,
				ParentID:    3,
				Name:        "AI音乐",
				DisplayName: "AI音乐",
				Slug:        "ai_music",
				Description: "以AI合成技术为基础，运用各类工具进行的AI作编曲、AI作词、AI语音、AI变声、AI翻唱、AI MV等创作",
				URI:         "/v/music/ai_music",
			},
			{
				ID:          267,
				ParentID:    3,
				Name:        "电台",
				DisplayName: "电台",
				Slug:        "radio",
				Description: "音乐分享、播单、白噪音、有声读物等以听为主的播放内容",
				URI:         "/v/music/radio",
			},
			{
				ID:          244,
				ParentID:    3,
				Name:        "音乐教学",
				DisplayName: "音乐教学",
				Slug:        "tutorial",
				Description: "以音乐教学为目的的内容",
				URI:         "/v/music/tutorial",
			},
			{
				ID:          130,
				ParentID:    3,
				Name:        "音乐综合",
				DisplayName: "音乐综合",
				Slug:        "other",
				Description: "所有无法被收纳到其他音乐二级分区的音乐类视频",
				URI:         "/v/music/other",
			},
			{
				ID:           194,
				ParentID:     3,
				Name:         "电音",
				DisplayName:  "电音(已下线)",
				Slug:         "electronic",
				Description:  "以电子合成器、音乐软体等产生的电子声响制作的音乐",
				URI:          "/v/music/electronic",
				IsDeprecated: true,
			},
		},
	},
	{
		ID:          129,
		Name:        "舞蹈",
		DisplayName: "舞蹈(主分区)",
		Slug:        "dance",
		URI:         "/v/dance",
		IsMain:      true,
		Children: []VideoZone{
			{
				ID:          20,
				ParentID:    129,
				Name:        "宅舞",
				DisplayName: "宅舞",
				Slug:        "otaku",
				Description: "与ACG相关的翻跳、原创舞蹈",
				URI:         "/v/dance/otaku",
			},
			{
				ID:          198,
				ParentID:    129,
				Name:        "街舞",
				DisplayName: "街舞",
				Slug:        "hiphop",
				Description: "收录街舞相关内容，包括赛事现场、舞室作品、个人翻跳、FREESTYLE等",
				URI:         "/v/dance/hiphop",
			},
			{
				ID:          199,
				ParentID:    129,
				Name:        "明星舞蹈",
				DisplayName: "明星舞蹈",
				Slug:        "star",
				Description: "国内外明星发布的官方舞蹈及其翻跳内容",
				URI:         "/v/dance/star",
			},
			{
				ID:          200,
				ParentID:    129,
				Name:        "国风舞蹈",
				DisplayName: "国风舞蹈",
				Slug:        "china",
				Description: "收录国风向舞蹈内容，包括中国舞、民族民间舞、汉唐舞、国风爵士等",
				URI:         "/v/dance/china",
			},
			{
				ID:          255,
				ParentID:    129,
				Name:        "颜值·网红舞",
				DisplayName: "颜值·网红舞 (原手势·网红舞)",
				Slug:        "gestures",
				Description: "手势舞及网红流行舞蹈、短视频舞蹈等相关视频",
				URI:         "/v/dance/gestures",
				Note:        "原手势·网红舞",
			},
			{
				ID:          154,
				ParentID:    129,
				Name:        "舞蹈综合",
				DisplayName: "舞蹈综合",
				Slug:        "three_d",
				Description: "收录无法定义到其他舞蹈子分区的舞蹈视频",
				URI:         "/v/dance/three_d",
			},
			{
				ID:          156,
				ParentID:    129,
				Name:        "舞蹈教程",
				DisplayName: "舞蹈教程",
				Slug:        "demo",
				Description: "镜面慢速，动作分解，基础教程等具有教学意义的舞蹈视频",
				URI:         "/v/dance/demo",
			},
		},
	},
	{
		ID:          4,
		Name:        "游戏",
		DisplayName: "游戏(主分区)",
		Slug:        "game",
		URI:         "/v/game",
		IsMain:      true,
		Children: []VideoZone{
			{
				ID:          17,
				ParentID:    4,
				Name:        "单机游戏",
				DisplayName: "单机游戏",
				Slug:        "stand_alone",
				Description: "以所有平台（PC、主机、移动端）的单机或联机游戏为主的视频内容，包括游戏预告、CG、实况解说及相关的评测、杂谈与视频剪辑等",
				URI:         "/v/game/stand_alone",
			},
			{
				ID:          171,
				ParentID:    4,
				Name:        "电子竞技",
				DisplayName: "电子竞技",
				Slug:        "esports",
				Description: "具有高对抗性的电子竞技游戏项目，其相关的赛事、实况、攻略、解说、短剧等视频。",
				URI:         "/v/game/esports",
			},
			{
				ID:          172,
				ParentID:    4,
				Name:        "手机游戏",
				DisplayName: "手机游戏",
				Slug:        "mobile",
				Description: "以手机及平板设备为主要平台的游戏，其相关的实况、攻略、解说、短剧、演示等视频。",
				URI:         "/v/game/mobile",
			},
			{
				ID:          65,
				ParentID:    4,
				Name:        "网络游戏",
				DisplayName: "网络游戏",
				Slug:        "online",
				Description: "由网络运营商运营的多人在线游戏，以及电子竞技的相关游戏内容。包括赛事、攻略、实况、解说等相关视频",
				URI:         "/v/game/online",
			},
			{
				ID:          173,
				ParentID:    4,
				Name:        "桌游棋牌",
				DisplayName: "桌游棋牌",
				Slug:        "board",
				Description: "桌游、棋牌、卡牌对战等及其相关电子版游戏的实况、攻略、解说、演示等视频。",
				URI:         "/v/game/board",
			},
			{
				ID:          121,
				ParentID:    4,
				Name:        "GMV",
				DisplayName: "GMV",
				Slug:        "gmv",
				Description: "由游戏素材制作的MV视频。以游戏内容或CG为主制作的，具有一定创作程度的MV类型的视频",
				URI:         "/v/game/gmv",
			},
			{
				ID:          136,
				ParentID:    4,
				Name:        "音游",
				DisplayName: "音游",
				Slug:        "music",
				Description: "各个平台上，通过配合音乐与节奏而进行的音乐类游戏视频",
				URI:         "/v/game/music",
			},
			{
				ID:          19,
				ParentID:    4,
				Name:        "Mugen",
				DisplayName: "Mugen",
				Slug:        "mugen",
				Description: "以Mugen引擎为平台制作、或与Mugen相关的游戏视频",
				URI:         "/v/game/mugen",
			},
		},
	},
	{
		ID:          36,
		Name:        "知识",
		DisplayName: "知识(主分区)",
		Slug:        "knowledge",
		URI:         "/v/knowledge",
		IsMain:      true,
		Children: []VideoZone{
			{
				ID:          201,
				ParentID:    36,
				Name:        "科学科普",
				DisplayName: "科学科普",
				Slug:        "science",
				Description: "回答你的十万个为什么",
				URI:         "/v/knowledge/science",
			},
			{
				ID:          124,
				ParentID:    36,
				Name:        "社科·法律·心理",
				DisplayName: "社科·法律·心理(原社科人文、原趣味科普人文)",
				Slug:        "social_science",
				Description: "基于社会科学、法学、心理学展开或个人观点输出的知识视频",
				URI:         "/v/knowledge/social_science",
				Note:        "原社科人文、原趣味科普人文",
			},
			{
				ID:          228,
				ParentID:    36,
				Name:        "人文历史",
				DisplayName: "人文历史",
				Slug:        "humanity_history",
				Description: "看看古今人物，聊聊历史过往，品品文学典籍",
				URI:         "/v/knowledge/humanity_history",
			},
			{
				ID:          207,
				ParentID:    36,
				Name:        "财经商业",
				DisplayName: "财经商业",
				Slug:        "business",
				Description: "说金融市场，谈宏观经济，一起畅聊商业故事",
				URI:         "/v/knowledge/finance",
			},
			{
				ID:          208,
				ParentID:    36,
				Name:        "校园学习",
				DisplayName: "校园学习",
				Slug:        "campus",
				Description: "老师很有趣，学生也有才，我们一起搞学习",
				URI:         "/v/knowledge/campus",
			},
			{
				ID:          209,
				ParentID:    36,
				Name:        "职业职场",
				DisplayName: "职业职场",
				Slug:        "career",
				Description: "职业分享、升级指南，一起成为最有料的职场人",
				URI:         "/v/knowledge/career",
			},
			{
				ID:          229,
				ParentID:    36,
				Name:        "设计·创意",
				DisplayName: "设计·创意",
				Slug:        "design",
				Description: "天马行空，创意设计，都在这里",
				URI:         "/v/knowledge/design",
			},
			{
				ID:          122,
				ParentID:    36,
				Name:        "野生技术协会",
				DisplayName: "野生技术协会",
				Slug:        "skill",
				Description: "技能党集合，是时候展示真正的技术了",
				URI:         "/v/knowledge/skill",
			},
			{
				ID:           39,
				ParentID:     36,
				Name:         "演讲·公开课",
				DisplayName:  "演讲·公开课(已下线)",
				Slug:         "speech_course",
				Description:  "涨知识的好地方，给爱学习的你",
				URI:          "/v/technology/speech_course",
				IsDeprecated: true,
			},
			{
				ID:           96,
				ParentID:     36,
				Name:         "星海",
				DisplayName:  "星海(已下线)",
				Slug:         "military",
				Description:  "军事类内容的圣地",
				URI:          "/v/technology/military",
				IsDeprecated: true,
			},
			{
				ID:           98,
				ParentID:     36,
				Name:         "机械",
				DisplayName:  "机械(已下线)",
				Slug:         "mechanical",
				Description:  "机械设备展示或制作视频",
				URI:          "/v/technology/mechanical",
				IsDeprecated: true,
			},
		},
	},
	{
		ID:          188,
		Name:        "科技",
		DisplayName: "科技(主分区)",
		Slug:        "tech",
		URI:         "/v/tech",
		IsMain:      true,
		Children: []VideoZone{
			{
				ID:          95,
				ParentID:    188,
				Name:        "数码",
				DisplayName: "数码(原手机平板)",
				Slug:        "digital",
				Description: "科技数码产品大全，一起来做发烧友",
				URI:         "/v/tech/digital",
				Note:        "原手机平板",
			},
			{
				ID:          230,
				ParentID:    188,
				Name:        "软件应用",
				DisplayName: "软件应用",
				Slug:        "application",
				Description: "超全软件应用指南",
				URI:         "/v/tech/application",
			},
			{
				ID:          231,
				ParentID:    188,
				Name:        "计算机技术",
				DisplayName: "计算机技术",
				Slug:        "computer_tech",
				Description: "研究分析、教学演示、经验分享......有关计算机技术的都在这里",
				URI:         "/v/tech/computer_tech",
			},
			{
				ID:          232,
				ParentID:    188,
				Name:        "科工机械",
				DisplayName: "科工机械 (原工业·工程·机械)",
				Slug:        "industry",
				Description: "从小芯片到大工程，一起见证科工力量",
				URI:         "/v/tech/industry",
				Note:        "原工业·工程·机械",
			},
			{
				ID:          233,
				ParentID:    188,
				Name:        "极客DIY",
				DisplayName: "极客DIY",
				Slug:        "diy",
				Description: "炫酷技能，极客文化，硬核技巧，准备好你的惊讶",
				URI:         "/v/tech/diy",
			},
			{
				ID:           189,
				ParentID:     188,
				Name:         "电脑装机",
				DisplayName:  "电脑装机(已下线)",
				Slug:         "pc",
				Description:  "电脑、笔记本、装机配件、外设和软件教程等相关视频",
				URI:          "/v/digital/pc",
				IsDeprecated: true,
			},
			{
				ID:           190,
				ParentID:     188,
				Name:         "摄影摄像",
				DisplayName:  "摄影摄像(已下线)",
				Slug:         "photography",
				Description:  "摄影摄像器材、拍摄剪辑技巧、拍摄作品分享等相关视频",
				URI:          "/v/digital/photography",
				IsDeprecated: true,
			},
			{
				ID:           191,
				ParentID:     188,
				Name:         "影音智能",
				DisplayName:  "影音智能(已下线)",
				Slug:         "intelligence_av",
				Description:  "影音设备、智能产品等相关视频",
				URI:          "/v/digital/intelligence_av",
				IsDeprecated: true,
			},
		},
	},
	{
		ID:          234,
		Name:        "运动",
		DisplayName: "运动(主分区)",
		Slug:        "sports",
		URI:         "/v/sports",
		IsMain:      true,
		Children: []VideoZone{
			{
				ID:          235,
				ParentID:    234,
				Name:        "篮球",
				DisplayName: "篮球",
				Slug:        "basketball",
				Description: "与篮球相关的视频，包括但不限于篮球赛事、教学、评述、剪辑、剧情等相关内容",
				URI:         "/v/sports/basketball",
			},
			{
				ID:          249,
				ParentID:    234,
				Name:        "足球",
				DisplayName: "足球",
				Slug:        "football",
				Description: "与足球相关的视频，包括但不限于足球赛事、教学、评述、剪辑、剧情等相关内容",
				URI:         "/v/sports/football",
			},
			{
				ID:          164,
				ParentID:    234,
				Name:        "健身",
				DisplayName: "健身",
				Slug:        "aerobics",
				Description: "与健身相关的视频，包括但不限于瑜伽、CrossFit、健美、力量举、普拉提、街健等相关内容",
				URI:         "/v/sports/aerobics",
			},
			{
				ID:          236,
				ParentID:    234,
				Name:        "竞技体育",
				DisplayName: "竞技体育",
				Slug:        "athletic",
				Description: "与竞技体育相关的视频，包括但不限于乒乓、羽毛球、排球、赛车等竞技项目的赛事、评述、剪辑、剧情等相关内容",
				URI:         "/v/sports/culture",
			},
			{
				ID:          237,
				ParentID:    234,
				Name:        "运动文化",
				DisplayName: "运动文化",
				Slug:        "culture",
				Description: "与运动文化相关的视频，包络但不限于球鞋、球衣、球星卡等运动衍生品的分享、解读，体育产业的分析、科普等相关内容",
				URI:         "/v/sports/culture",
			},
			{
				ID:          238,
				ParentID:    234,
				Name:        "运动综合",
				DisplayName: "运动综合",
				Slug:        "comprehensive",
				Description: "与运动综合相关的视频，包括但不限于钓鱼、骑行、滑板等日常运动分享、教学、Vlog等相关内容",
				URI:         "/v/sports/comprehensive",
			},
		},
	},
	{
		ID:          223,
		Name:        "汽车",
		DisplayName: "汽车(主分区)",
		Slug:        "car",
		URI:         "/v/car",
		IsMain:      true,
		Children: []VideoZone{
			{
				ID:          258,
				ParentID:    223,
				Name:        "汽车知识科普",
				DisplayName: "汽车知识科普",
				Slug:        "knowledge",
				Description: "关于汽车技术与文化的硬核科普，以及生活中学车、用车、养车的相关知识",
				URI:         "/v/car/knowledge",
			},
			{
				ID:          227,
				ParentID:    223,
				Name:        "购车攻略",
				DisplayName: "购车攻略",
				Slug:        "strategy",
				Description: "丰富详实的购车建议和新车体验",
				URI:         "/v/car/strategy",
			},
			{
				ID:          247,
				ParentID:    223,
				Name:        "新能源车",
				DisplayName: "新能源车",
				Slug:        "newenergyvehicle",
				Description: "新能源汽车相关内容，包括电动汽车、混合动力汽车等车型种类，包含不限于新车资讯、试驾体验、专业评测、技术解读、知识科普等内容",
				URI:         "/v/car/newenergyvehicle",
			},
			{
				ID:          245,
				ParentID:    223,
				Name:        "赛车",
				DisplayName: "赛车",
				Slug:        "racing",
				Description: "F1等汽车运动相关",
				URI:         "/v/car/racing",
			},
			{
				ID:          246,
				ParentID:    223,
				Name:        "改装玩车",
				DisplayName: "改装玩车",
				Slug:        "modifiedvehicle",
				Description: "汽车文化及改装车相关内容，包括改装车、老车修复介绍、汽车聚会分享等内容",
				URI:         "/v/car/modifiedvehicle",
			},
			{
				ID:          240,
				ParentID:    223,
				Name:        "摩托车",
				DisplayName: "摩托车",
				Slug:        "motorcycle",
				Description: "骑士们集合啦",
				URI:         "/v/car/motorcycle",
			},
			{
				ID:          248,
				ParentID:    223,
				Name:        "房车",
				DisplayName: "房车",
				Slug:        "touringcar",
				Description: "房车及营地相关内容，包括不限于产品介绍、驾驶体验、房车生活和房车旅行等内容",
				URI:         "/v/car/touringcar",
			},
			{
				ID:          176,
				ParentID:    223,
				Name:        "汽车生活",
				DisplayName: "汽车生活",
				Slug:        "life",
				Description: "分享汽车及出行相关的生活体验类视频",
				URI:         "/v/car/life",
			},
			{
				ID:           224,
				ParentID:     223,
				Name:         "汽车文化",
				DisplayName:  "汽车文化(已下线)",
				Slug:         "culture",
				Description:  "车迷的精神圣地，包括汽车赛事、品牌历史、汽车改装、经典车型和汽车模型等",
				URI:          "/v/car/culture",
				IsDeprecated: true,
			},
			{
				ID:           225,
				ParentID:     223,
				Name:         "汽车极客",
				DisplayName:  "汽车极客(已下线)",
				Slug:         "geek",
				Description:  "汽车硬核达人聚集地，包括DIY造车、专业评测和技术知识分享",
				URI:          "/v/car/geek",
				IsDeprecated: true,
			},
			{
				ID:           226,
				ParentID:     223,
				Name:         "智能出行",
				DisplayName:  "智能出行(已下线)",
				Slug:         "smart",
				Description:  "探索新能源汽车和未来智能出行的前沿阵地",
				URI:          "/v/car/smart",
				IsDeprecated: true,
			},
		},
	},
	{
		ID:          160,
		Name:        "生活",
		DisplayName: "生活(主分区)",
		Slug:        "life",
		URI:         "/v/life",
		IsMain:      true,
		Children: []VideoZone{
			{
				ID:          138,
				ParentID:    160,
				Name:        "搞笑",
				DisplayName: "搞笑",
				Slug:        "funny",
				Description: "各种沙雕有趣的搞笑剪辑，挑战，表演，配音等视频",
				URI:         "/v/life/funny",
			},
			{
				ID:          254,
				ParentID:    160,
				Name:        "亲子",
				DisplayName: "亲子",
				Slug:        "parenting",
				Description: "分享亲子、萌娃、母婴、育儿相关的视频",
				URI:         "/v/life/parenting",
			},
			{
				ID:          250,
				ParentID:    160,
				Name:        "出行",
				DisplayName: "出行",
				Slug:        "travel",
				Description: "为达到观光游览、休闲娱乐为目的的远途旅行、中近途户外生活、本地探店",
				URI:         "/v/life/travel",
			},
			{
				ID:          251,
				ParentID:    160,
				Name:        "三农",
				DisplayName: "三农",
				Slug:        "rurallife",
				Description: "分享美好农村生活",
				URI:         "/v/life/rurallife",
			},
			{
				ID:          239,
				ParentID:    160,
				Name:        "家居房产",
				DisplayName: "家居房产",
				Slug:        "home",
				Description: "与买房、装修、居家生活相关的分享",
				URI:         "/v/life/home",
			},
			{
				ID:          161,
				ParentID:    160,
				Name:        "手工",
				DisplayName: "手工",
				Slug:        "handmake",
				Description: "手工制品的制作过程或成品展示、教程、测评类视频",
				URI:         "/v/life/handmake",
			},
			{
				ID:          162,
				ParentID:    160,
				Name:        "绘画",
				DisplayName: "绘画",
				Slug:        "painting",
				Description: "绘画过程或绘画教程，以及绘画相关的所有视频",
				URI:         "/v/life/painting",
			},
			{
				ID:          21,
				ParentID:    160,
				Name:        "日常",
				DisplayName: "日常",
				Slug:        "daily",
				Description: "记录日常生活，分享生活故事",
				URI:         "/v/life/daily",
			},
			{
				ID:          76,
				ParentID:    160,
				Name:        "美食圈",
				DisplayName: "美食圈(重定向)",
				Slug:        "food",
				Description: "美食鉴赏&料理制作教程",
				URI:         "/v/life/food",
				IsRedirect:  true,
			},
			{
				ID:          75,
				ParentID:    160,
				Name:        "动物圈",
				DisplayName: "动物圈(重定向)",
				Slug:        "animal",
				Description: "萌萌的动物都在这里哦",
				URI:         "/v/life/animal",
				IsRedirect:  true,
			},
			{
				ID:          163,
				ParentID:    160,
				Name:        "运动",
				DisplayName: "运动(重定向)",
				Slug:        "sports",
				Description: "运动相关的记录、教程、装备评测和精彩瞬间剪辑视频",
				URI:         "/v/life/sports",
				IsRedirect:  true,
			},
			{
				ID:          176,
				ParentID:    160,
				Name:        "汽车",
				DisplayName: "汽车(重定向)",
				Slug:        "automobile",
				Description: "专业汽车资讯，分享车生活",
				URI:         "/v/life/automobile",
				IsRedirect:  true,
			},
			{
				ID:           174,
				ParentID:     160,
				Name:         "其他",
				DisplayName:  "其他(已下线)",
				Slug:         "other",
				Description:  "对于分区归属不明的视频进行归纳整合的特定分区",
				URI:          "/v/life/other",
				IsDeprecated: true,
			},
		},
	},
	{
		ID:          211,
		Name:        "美食",
		DisplayName: "美食(主分区)",
		Slug:        "food",
		URI:         "/v/food",
		IsMain:      true,
		Children: []VideoZone{
			{
				ID:          76,
				ParentID:    211,
				Name:        "美食制作",
				DisplayName: "美食制作(原[生活]->[美食圈])",
				Slug:        "make",
				Description: "学做人间美味，展示精湛厨艺",
				URI:         "/v/food/make",
				Note:        "原[生活]->[美食圈]",
			},
			{
				ID:          212,
				ParentID:    211,
				Name:        "美食侦探",
				DisplayName: "美食侦探",
				Slug:        "detective",
				Description: "寻找美味餐厅，发现街头美食",
				URI:         "/v/food/detective",
			},
			{
				ID:          213,
				ParentID:    211,
				Name:        "美食测评",
				DisplayName: "美食测评",
				Slug:        "measurement",
				Description: "吃货世界，品尝世间美味",
				URI:         "/v/food/measurement",
			},
			{
				ID:          214,
				ParentID:    211,
				Name:        "田园美食",
				DisplayName: "田园美食",
				Slug:        "rural",
				Description: "品味乡野美食，寻找山与海的味道",
				URI:         "/v/food/rural",
			},
			{
				ID:          215,
				ParentID:    211,
				Name:        "美食记录",
				DisplayName: "美食记录",
				Slug:        "record",
				Description: "记录一日三餐，给生活添一点幸福感",
				URI:         "/v/food/record",
			},
		},
	},
	{
		ID:          217,
		Name:        "动物圈",
		DisplayName: "动物圈(主分区)",
		Slug:        "animal",
		URI:         "/v/animal",
		IsMain:      true,
		Children: []VideoZone{
			{
				ID:          218,
				ParentID:    217,
				Name:        "喵星人",
				DisplayName: "喵星人",
				Slug:        "cat",
				Description: "喵喵喵喵喵",
				URI:         "/v/animal/cat",
			},
			{
				ID:          219,
				ParentID:    217,
				Name:        "汪星人",
				DisplayName: "汪星人",
				Slug:        "dog",
				Description: "汪汪汪汪汪",
				URI:         "/v/animal/dog",
			},
			{
				ID:          222,
				ParentID:    217,
				Name:        "小宠异宠",
				DisplayName: "小宠异宠",
				Slug:        "reptiles",
				Description: "奇妙宠物大赏",
				URI:         "/v/animal/reptiles",
			},
			{
				ID:          221,
				ParentID:    217,
				Name:        "野生动物",
				DisplayName: "野生动物",
				Slug:        "wild_animal",
				Description: "内有“猛兽”出没",
				URI:         "/v/animal/wild_animal",
			},
			{
				ID:          220,
				ParentID:    217,
				Name:        "动物二创",
				DisplayName: "动物二创",
				Slug:        "second_edition",
				Description: "解说、配音、剪辑、混剪",
				URI:         "/v/animal/second_edition",
			},
			{
				ID:          75,
				ParentID:    217,
				Name:        "动物综合",
				DisplayName: "动物综合",
				Slug:        "animal_composite",
				Description: "收录除上述子分区外，其余动物相关视频以及非动物主体或多个动物主体的动物相关延伸内容",
				URI:         "/v/animal/animal_composite",
			},
		},
	},
	{
		ID:          119,
		Name:        "鬼畜",
		DisplayName: "鬼畜(主分区)",
		Slug:        "kichiku",
		URI:         "/v/kichiku",
		IsMain:      true,
		Children: []VideoZone{
			{
				ID:          22,
				ParentID:    119,
				Name:        "鬼畜调教",
				DisplayName: "鬼畜调教",
				Slug:        "guide",
				Description: "使用素材在音频、画面上做一定处理，达到与BGM一定的同步感",
				URI:         "/v/kichiku/guide",
			},
			{
				ID:          26,
				ParentID:    119,
				Name:        "音MAD",
				DisplayName: "音MAD",
				Slug:        "mad",
				Description: "使用素材音频进行一定的二次创作来达到还原原曲的非商业性质稿件",
				URI:         "/v/kichiku/mad",
			},
			{
				ID:          126,
				ParentID:    119,
				Name:        "人力VOCALOID",
				DisplayName: "人力VOCALOID",
				Slug:        "manual_vocaloid",
				Description: "将人物或者角色的无伴奏素材进行人工调音，使其就像VOCALOID一样歌唱的技术",
				URI:         "/v/kichiku/manual_vocaloid",
			},
			{
				ID:          216,
				ParentID:    119,
				Name:        "鬼畜剧场",
				DisplayName: "鬼畜剧场",
				Slug:        "theatre",
				Description: "使用素材进行人工剪辑编排的有剧情的作品",
				URI:         "/v/kichiku/theatre",
			},
			{
				ID:          127,
				ParentID:    119,
				Name:        "教程演示",
				DisplayName: "教程演示",
				Slug:        "course",
				Description: "鬼畜相关的教程演示",
				URI:         "/v/kichiku/course",
			},
		},
	},
	{
		ID:          155,
		Name:        "时尚",
		DisplayName: "时尚(主分区)",
		Slug:        "fashion",
		URI:         "/v/fashion",
		IsMain:      true,
		Children: []VideoZone{
			{
				ID:          157,
				ParentID:    155,
				Name:        "美妆护肤",
				DisplayName: "美妆护肤",
				Slug:        "makeup",
				Description: "彩妆护肤、美甲美发、仿妆、医美相关内容分享或产品测评",
				URI:         "/v/fashion/makeup",
			},
			{
				ID:          252,
				ParentID:    155,
				Name:        "仿妆cos",
				DisplayName: "仿妆cos",
				Slug:        "cos",
				Description: "对二次元、三次元人物角色进行模仿、还原、展示、演绎的内容",
				URI:         "/v/fashion/cos",
			},
			{
				ID:          158,
				ParentID:    155,
				Name:        "穿搭",
				DisplayName: "穿搭",
				Slug:        "clothing",
				Description: "穿搭风格、穿搭技巧的展示分享，涵盖衣服、鞋靴、箱包配件、配饰（帽子、钟表、珠宝首饰）等",
				URI:         "/v/fashion/clothing",
			},
			{
				ID:          159,
				ParentID:    155,
				Name:        "时尚潮流",
				DisplayName: "时尚潮流",
				Slug:        "catwalk",
				Description: "时尚街拍、时装周、时尚大片，时尚品牌、潮流等行业相关记录及知识科普",
				URI:         "/v/fashion/catwalk",
			},
			{
				ID:          164,
				ParentID:    155,
				Name:        "健身",
				DisplayName: "健身(重定向)",
				Slug:        "aerobics",
				Description: "器械、有氧、拉伸运动等，以达到强身健体、减肥瘦身、形体塑造目的",
				URI:         "/v/fashion/aerobics",
				IsRedirect:  true,
			},
			{
				ID:           192,
				ParentID:     155,
				Name:         "风尚标",
				DisplayName:  "风尚标(已下线)",
				Slug:         "trends",
				Description:  "时尚明星专访、街拍、时尚购物相关知识科普",
				URI:          "/v/fashion/trends",
				IsDeprecated: true,
			},
		},
	},
	{
		ID:          202,
		Name:        "资讯",
		DisplayName: "资讯(主分区)",
		Slug:        "information",
		URI:         "/v/information",
		IsMain:      true,
		Children: []VideoZone{
			{
				ID:          203,
				ParentID:    202,
				Name:        "热点",
				DisplayName: "热点",
				Slug:        "hotspot",
				Description: "全民关注的时政热门资讯",
				URI:         "/v/information/hotspot",
			},
			{
				ID:          204,
				ParentID:    202,
				Name:        "环球",
				DisplayName: "环球",
				Slug:        "global",
				Description: "全球范围内发生的具有重大影响力的事件动态",
				URI:         "/v/information/global",
			},
			{
				ID:          205,
				ParentID:    202,
				Name:        "社会",
				DisplayName: "社会",
				Slug:        "social",
				Description: "日常生活的社会事件、社会问题、社会风貌的报道",
				URI:         "/v/information/social",
			},
			{
				ID:          206,
				ParentID:    202,
				Name:        "综合",
				DisplayName: "综合",
				Slug:        "multiple",
				Description: "除上述领域外其它垂直领域的综合资讯",
				URI:         "/v/information/multiple",
			},
		},
	},
	{
		ID:           165,
		Name:         "广告",
		DisplayName:  "广告(主分区)",
		Slug:         "ad",
		URI:          "/v/ad",
		IsMain:       true,
		IsDeprecated: true,
		Children: []VideoZone{
			{
				ID:           166,
				ParentID:     165,
				Name:         "广告",
				DisplayName:  "广告(已下线)",
				Slug:         "ad",
				URI:          "/v/ad/ad",
				IsDeprecated: true,
			},
		},
	},
	{
		ID:          5,
		Name:        "娱乐",
		DisplayName: "娱乐(主分区)",
		Slug:        "ent",
		URI:         "/v/ent",
		IsMain:      true,
		Children: []VideoZone{
			{
				ID:          241,
				ParentID:    5,
				Name:        "娱乐杂谈",
				DisplayName: "娱乐杂谈",
				Slug:        "talker",
				Description: "娱乐人物解读、娱乐热点点评、娱乐行业分析",
				URI:         "/v/ent/talker",
			},
			{
				ID:          262,
				ParentID:    5,
				Name:        "CP安利",
				DisplayName: "CP安利",
				Slug:        "cp_recommendation",
				Description: "以安利各类娱乐名人、角色CP之间默契于火花为主题的混剪、解说，观点表达类视频",
				URI:         "/v/ent/cp_recommendation",
			},
			{
				ID:          263,
				ParentID:    5,
				Name:        "颜值安利",
				DisplayName: "颜值安利",
				Slug:        "beauty",
				Description: "以各类娱乐名人、角色的颜值、气质魅力为核心的混剪视频",
				URI:         "/v/ent/beauty",
			},
			{
				ID:          242,
				ParentID:    5,
				Name:        "娱乐粉丝创作",
				DisplayName: "娱乐粉丝创作 (原粉丝创作)",
				Slug:        "fans",
				Description: "粉丝向创作视频",
				URI:         "/v/ent/fans",
				Note:        "原粉丝创作",
			},
			{
				ID:          264,
				ParentID:    5,
				Name:        "娱乐资讯",
				DisplayName: "娱乐资讯",
				Slug:        "entertainment_news",
				Description: "具备趣味价值的文化娱乐新闻与动态报道，如名人动态，作品发布，舞台演出，趣闻盘点等",
				URI:         "/v/ent/entertainment_news",
			},
			{
				ID:          137,
				ParentID:    5,
				Name:        "明星综合",
				DisplayName: "明星综合",
				Slug:        "celebrity",
				Description: "娱乐圈动态、明星资讯相关",
				URI:         "/v/ent/celebrity",
			},
			{
				ID:          71,
				ParentID:    5,
				Name:        "综艺",
				DisplayName: "综艺",
				Slug:        "variety",
				Description: "所有综艺相关，全部一手掌握！",
				URI:         "/v/ent/variety",
			},
			{
				ID:           131,
				ParentID:     5,
				Name:         "Korea相关",
				DisplayName:  "Korea相关(已下线)",
				Slug:         "korea",
				Description:  "Korea相关音乐、舞蹈、综艺等视频",
				URI:          "/v/ent/korea",
				IsDeprecated: true,
			},
		},
	},
	{
		ID:          181,
		Name:        "影视",
		DisplayName: "影视(主分区)",
		Slug:        "cinephile",
		URI:         "/v/cinephile",
		IsMain:      true,
		Children: []VideoZone{
			{
				ID:          182,
				ParentID:    181,
				Name:        "影视杂谈",
				DisplayName: "影视杂谈",
				Slug:        "cinecism",
				Description: "影视评论、解说、吐槽、科普等",
				URI:         "/v/cinephile/cinecism",
			},
			{
				ID:          183,
				ParentID:    181,
				Name:        "影视剪辑",
				DisplayName: "影视剪辑",
				Slug:        "montage",
				Description: "对影视素材进行剪辑再创作的视频",
				URI:         "/v/cinephile/montage",
			},
			{
				ID:          260,
				ParentID:    181,
				Name:        "影视整活",
				DisplayName: "影视整活",
				Slug:        "mashup",
				Description: "使用影视素材制造的有趣、有梗的创意混剪、配音、特效玩梗视频",
				URI:         "/v/cinephile/mashup",
			},
			{
				ID:          259,
				ParentID:    181,
				Name:        "AI影像",
				DisplayName: "AI影像",
				Slug:        "ai_imaging",
				Description: "分享AI制作的影像作品、创作历程、技术风向",
				URI:         "/v/cinephile/ai_imaging",
			},
			{
				ID:          184,
				ParentID:    181,
				Name:        "预告·资讯",
				DisplayName: "预告·资讯",
				Slug:        "trailer_info",
				Description: "影视类相关资讯，预告，花絮等视频",
				URI:         "/v/cinephile/trailer_info",
			},
			{
				ID:          85,
				ParentID:    181,
				Name:        "小剧场",
				DisplayName: "小剧场",
				Slug:        "shortplay",
				Description: "有场景、有剧情的演绎类内容",
				URI:         "/v/cinephile/shortplay",
			},
			{
				ID:          256,
				ParentID:    181,
				Name:        "短片",
				DisplayName: "短片",
				Slug:        "shortfilm",
				Description: "各种类型的短片",
				URI:         "/v/cinephile/shortfilm",
			},
			{
				ID:          261,
				ParentID:    181,
				Name:        "影视综合",
				DisplayName: "影视综合",
				Slug:        "comprehensive",
				Description: "一切无法被收纳其他影视二级分区的影视相关内容",
				URI:         "/v/cinephile/comprehensive",
			},
		},
	},
	{
		ID:          177,
		Name:        "纪录片",
		DisplayName: "纪录片(主分区)",
		Slug:        "documentary",
		URI:         "/documentary",
		IsMain:      true,
		Children: []VideoZone{
			{
				ID:          37,
				ParentID:    177,
				Name:        "人文·历史",
				DisplayName: "人文·历史",
				Slug:        "history",
				URI:         "/v/documentary/history",
			},
			{
				ID:          178,
				ParentID:    177,
				Name:        "科学·探索·自然",
				DisplayName: "科学·探索·自然",
				Slug:        "science",
				URI:         "/v/documentary/science",
			},
			{
				ID:          179,
				ParentID:    177,
				Name:        "军事",
				DisplayName: "军事",
				Slug:        "military",
				URI:         "/v/documentary/military",
			},
			{
				ID:          180,
				ParentID:    177,
				Name:        "社会·美食·旅行",
				DisplayName: "社会·美食·旅行",
				Slug:        "travel",
				URI:         "/v/documentary/travel",
			},
		},
	},
	{
		ID:          23,
		Name:        "电影",
		DisplayName: "电影(主分区)",
		Slug:        "movie",
		URI:         "/movie",
		IsMain:      true,
		Children: []VideoZone{
			{
				ID:          147,
				ParentID:    23,
				Name:        "华语电影",
				DisplayName: "华语电影",
				Slug:        "chinese",
				URI:         "/v/movie/chinese",
			},
			{
				ID:          145,
				ParentID:    23,
				Name:        "欧美电影",
				DisplayName: "欧美电影",
				Slug:        "west",
				URI:         "/v/movie/west",
			},
			{
				ID:          146,
				ParentID:    23,
				Name:        "日本电影",
				DisplayName: "日本电影",
				Slug:        "japan",
				URI:         "/v/movie/japan",
			},
			{
				ID:          83,
				ParentID:    23,
				Name:        "其他国家",
				DisplayName: "其他国家",
				Slug:        "movie",
				URI:         "/v/movie/movie",
			},
		},
	},
	{
		ID:          11,
		Name:        "电视剧",
		DisplayName: "电视剧(主分区)",
		Slug:        "tv",
		URI:         "/tv",
		IsMain:      true,
		Children: []VideoZone{
			{
				ID:          185,
				ParentID:    11,
				Name:        "国产剧",
				DisplayName: "国产剧",
				Slug:        "mainland",
				URI:         "/v/tv/mainland",
			},
			{
				ID:          187,
				ParentID:    11,
				Name:        "海外剧",
				DisplayName: "海外剧",
				Slug:        "overseas",
				URI:         "/v/tv/overseas",
			},
		},
	},
}

var videoZonesFlat = flattenVideoZones(videoZones)
var videoZoneByTID = makeVideoZoneByTID(videoZonesFlat)

// GetVideoZones 返回完整的视频分区树。
func GetVideoZones() []VideoZone {
	return cloneVideoZones(videoZones)
}

// GetMainVideoZones 返回所有一级视频分区。
func GetMainVideoZones() []VideoZone {
	return cloneVideoZones(videoZones)
}

// GetAllVideoZones 返回扁平化的视频分区列表。
func GetAllVideoZones() []VideoZone {
	return cloneVideoZones(videoZonesFlat)
}

// GetVideoZoneByTID 按 tid 查询视频分区。
func GetVideoZoneByTID(tid int) (*VideoZone, bool) {
	matches := FindVideoZonesByTID(tid)
	if len(matches) == 0 {
		return nil, false
	}
	best := matches[0]
	for _, zone := range matches {
		if !zone.IsRedirect && !zone.IsDeprecated {
			best = zone
			break
		}
	}
	return &best, true
}

// GetChildVideoZones 返回指定主分区下的子分区。
func GetChildVideoZones(parentTID int) []VideoZone {
	for _, zone := range videoZones {
		if zone.ID == parentTID {
			return cloneVideoZones(zone.Children)
		}
	}
	return nil
}

// FindVideoZonesByTID 按 tid 查询视频分区，tid 重复时返回多个结果。
func FindVideoZonesByTID(tid int) []VideoZone {
	if tid <= 0 {
		return nil
	}
	var matches []VideoZone
	for _, zone := range videoZonesFlat {
		if zone.ID == tid {
			matches = append(matches, cloneVideoZone(zone))
		}
	}
	return matches
}

// FindVideoZonesByName 按名称查询视频分区，名称重复时返回多个结果。
func FindVideoZonesByName(name string) []VideoZone {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil
	}
	var matches []VideoZone
	for _, zone := range videoZonesFlat {
		if zone.Name == name || zone.DisplayName == name {
			matches = append(matches, cloneVideoZone(zone))
		}
	}
	return matches
}

// FindVideoZonesBySlug 按代号查询视频分区，代号重复时返回多个结果。
func FindVideoZonesBySlug(slug string) []VideoZone {
	slug = strings.TrimSpace(strings.ToLower(slug))
	if slug == "" {
		return nil
	}
	var matches []VideoZone
	for _, zone := range videoZonesFlat {
		if strings.ToLower(zone.Slug) == slug {
			matches = append(matches, cloneVideoZone(zone))
		}
	}
	return matches
}

// FindVideoZonesByURI 按路由查询视频分区，路由重复时返回多个结果。
func FindVideoZonesByURI(uri string) []VideoZone {
	uri = strings.TrimSpace(uri)
	if uri == "" {
		return nil
	}
	var matches []VideoZone
	for _, zone := range videoZonesFlat {
		if zone.URI == uri {
			matches = append(matches, cloneVideoZone(zone))
		}
	}
	return matches
}

func flattenVideoZones(zones []VideoZone) []VideoZone {
	flat := make([]VideoZone, 0, len(zones))
	for _, zone := range zones {
		flat = append(flat, cloneVideoZoneWithoutChildren(zone))
		if len(zone.Children) > 0 {
			flat = append(flat, flattenVideoZones(zone.Children)...)
		}
	}
	return flat
}

func makeVideoZoneByTID(zones []VideoZone) map[int]VideoZone {
	indexed := make(map[int]VideoZone, len(zones))
	for i := range zones {
		zone := zones[i]
		if _, exists := indexed[zone.ID]; !exists && zone.ParentID == 0 {
			indexed[zone.ID] = zone
			continue
		}
		if zone.IsRedirect || zone.IsDeprecated {
			continue
		}
		indexed[zone.ID] = zone
	}
	return indexed
}

func cloneVideoZones(zones []VideoZone) []VideoZone {
	if len(zones) == 0 {
		return nil
	}
	cloned := make([]VideoZone, len(zones))
	for i := range zones {
		cloned[i] = cloneVideoZone(zones[i])
	}
	return cloned
}

func cloneVideoZone(zone VideoZone) VideoZone {
	cloned := cloneVideoZoneWithoutChildren(zone)
	if len(zone.Children) > 0 {
		cloned.Children = cloneVideoZones(zone.Children)
	}
	return cloned
}

func cloneVideoZoneWithoutChildren(zone VideoZone) VideoZone {
	zone.Children = nil
	return zone
}
