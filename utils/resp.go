package utils

//Response 客户端应答
type Response struct {
	//Code 错误码
	Code    int         `json:"Code"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}

//Response 客户端应答
type CloudResp struct {
	Success bool   `json:"Success"`
	Message string `json:"Message"`
	Key     string `json:"Key"`
}

// Category  Get all category Message struct
type Category struct {
	Path       string
	Name       string
	Count      int
	CreateTime int64
	UpdateTime int64
}

// NoteAllList  Get all note Message struct
type NoteAllList struct {
	Path         string `json:"Path"`
	CategoryPath string `bson:"CategoryPath"` //笔记分组唯一途径
	//	CategoryName string `bson:"CategoryName"` //笔记分组唯一途径
	UpdateTime     int64    `bson:"UpdateTime"` //笔记更新时间
	Summary        string   `bson:"Summary"`    //笔记摘要
	CheckSum       string   //笔记CheckSum
	AttachTypeList []string //笔记附件分组
	//	CreateTime   int64  `bson:"CreateTime"`   //笔记创建时间
	NoteName     string //笔记名称
	IsPinned     uint8  `bson:"IsPinned"` //笔记是否置顶
	IsShare      uint8  //笔记是否被分享
	IsVote       uint8  //笔记是否被投票
	Vote         interface{}
	IsPassworded uint8 //笔记是否加密分享
}

// NoteDeleteList  Get all note Message struct
type NoteDeleteList struct {
	Path           string   `json:"Path"`
	CategoryPath   string   `bson:"CategoryPath"` //笔记分组唯一途径
	UpdateTime     int64    `bson:"UpdateTime"`   //笔记更新时间
	NoteName       string   //笔记名称
	CheckSum       string   //笔记CheckSum
	Summary        string   `bson:"Summary"` //笔记摘要
	AttachTypeList []string //笔记附件分组
	Type           uint8    `bson:"Type"` //笔记是否置顶
	IsShare        uint8    //笔记是否被分享
	IsVote         uint8    //笔记是否被投票
}

// attachment  note's attachment message struct
type Attachment struct {
	Path   string `json:"Path"`
	Url    string `json:"Url"`
	Uuid   string `json:"Uuid"`
	Type   string `json:"Type"`
	ObType uint8  `json:"ObType"`
}

//Note one note messgae
type Note struct {
	Title          string       `bson:"Title"`    //笔记标题
	IsPinned       uint8        `bson:"IsPinned"` //是否置顶
	IsVote         uint8        //笔记是否被投票
	Usn            int          `json:"Usn,omitempty"`
	IsPassworded   uint8        //笔记是否加密分享
	Author         string       `bson:"Author"`       //笔记作者
	CategoryName   string       `bson:"CategoryName"` //笔记作者
	CategoryPath   string       `bson:"CategoryPath"` //笔记作者
	Summary        string       `bson:"Summary"`      //笔记摘要
	Content        string       `bson:"Content"`      //笔记内容
	CheckSum       string       `bson:"CheckSum"`     //笔记校验和
	CreateTime     int64        `bson:"CreateTime"`   //笔记创建时间
	UpdateTime     int64        `bson:"UpdateTime"`   //笔记更新时间
	AttachmentList []Attachment `json:"AttachmentList"`
	Vote           interface{}
}

//UserInfo user massage
type UserInfo struct {
	Rtncode     string `json:"rtn_code"`
	Sex         string `json:"sex"`
	Nickname    string `json:"nickname"`
	Brithday    string `json:"brithday"`
	HighDefUrl  string `json:"highDefUrl"`
	HeadIconUrl string `json:"HeadIconUrl"`
}

//AppResp app check login
type AppResp struct {
	Rtncode string `json:"rtncode"`
	Rtnmsg  string `json:"rtnmsg"`
}

//RefreshResp user massage
type RefreshResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    string `json:"expires_in"`
}

//VoteResp vote massage
type VoteResp struct {
	VoteID       string
	VoteQuestion string
	VoteType     uint8
	VoteSelect   interface{}
	VoteCount    int
}

type AttchType struct {
	HasImg   bool
	HasVideo bool
	HasMusic bool
}

//同步返回结构
type SyncResp struct {
	Path         string
	Title        string `json:"Title,omitempty"`
	CategoryPath string `json:"CategoryPath,omitempty"`
	CheckSum     string `json:"CheckSum,omitempty"`
	Name         string `json:"Name,omitempty"`
	Usn          int    `json:"Usn,omitempty"`
	IsDeleted    bool
	IsPinned     uint8 `json:"IsPinned,omitempty"`
	UpdateTime   int64
	IsVote       uint8 `json:"IsVote,omitempty"`       //笔记是否被投票
	IsPassworded uint8 `json:"IsPassworded,omitempty"` //笔记是否加密分享
}

type pushPayload struct {
	Environment string   `json:"Environment"`
	ApiType     string   `json:"ApiType"`
	TimeType    int      `json:"TimeType"`
	Expires     int64    `json:"Expires"`
	IsCallback  bool     `json:"IsCallback"`
	Body        string   `json:"Body"`
	PushType    string   `json:"PushType"`
	Target      []string `json:"Target"`
}
