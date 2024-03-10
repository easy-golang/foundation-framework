package dto

type ClientType string

const (
	Seller ClientType = "seller"
	Buyer  ClientType = "buyer"
	Admin  ClientType = "admin"
)

type Page[T any] struct {
	Data  []T   `json:"data"`
	Total int64 `json:"total"`
	Pages int64 `json:"pages"`
}

func NewPage[T any](data []T, total int64, pageSize int64) *Page[T] {
	return &Page[T]{
		Data:  data,
		Total: total,
		Pages: (total + pageSize - 1) / pageSize,
	}
}

type RequestContext struct {
	XReqeustContext XReqeustContext   `json:"x-request-context"`
	XRequestHeader  map[string]string `json:"x-request-header"`
	XRequestId      string            `json:"x-request-id"`
}

type XReqeustContext struct {
	AcceptCurrency string     `json:"acceptCurrency"` // 接受货币
	AcceptLanguage string     `json:"acceptLanguage"` // 接受语言
	ClientIp       string     `json:"clientIp"`       // 客户端IP
	IsAdmin        bool       `json:"isAdmin"`        // 是否是管理员
	LoginId        string     `json:"loginId"`        // 登录id
	RequestClient  ClientType `json:"requestClient"`  // 请求客户端
	ShopId         string     `json:"shopId"`         // 商店id
	ShopSettings   string     `json:"shopSettings"`   // 商店设置
	ShopTimezone   string     `json:"shopTimezone"`   // 商店时区
	UserAgent      string     `json:"userAgent"`      // 用户代理
	UserRegion     string     `json:"userRegion"`     // 用户区域
	UserInfo       *UserInfo  `json:"userInfo"`       // 用户信息
	ErpInfo        *ErpInfo   `json:"erpInfo"`        // Erp信息
}

type UserInfo struct {
	UserId            string `json:"userId"`
	UserName          string `json:"userName"`
	Email             string `json:"email"`
	Status            string `json:"status"`
	PreferredLanguage string `json:"preferredLanguage"`
}

type ErpInfo struct {
	ErpId   int64 `json:"erpId"`
	IsOwner bool  `json:"isOwner"`
}

type ShopSettings struct {
	CurrencyDisplayWay   int32         `json:"currencyDisplayWay"`   // 货币显示方式
	CurrencyList         []Currency    `json:"currencyList"`         // 货币列表
	DefaultDisplayWay    int32         `json:"defaultDisplayWay"`    // 默认显示方式
	DefaultLanguageMode  int32         `json:"defaultLanguageMode"`  // 默认语言模式
	Language             []string      `json:"language"`             // 语言
	MainCurrency         *MainCurrency `json:"mainCurrency"`         // 主要货币
	MainLanguage         string        `json:"mainLanguage"`         // 主语言
	Morelanguage         []string      `json:"morelanguage"`         // 更多语言
	Plugins              []int64       `json:"plugins"`              // 插件
	Status               bool          `json:"status"`               // 状态
	Timezone             string        `json:"timezone"`             // 时区
	UseFixedExchangeRate int32         `json:"useFixedExchangeRate"` // 使用固定汇率
	UseMultipleCurrency  int32         `json:"useMultipleCurrency"`  // 使用多种货币
}

type MainCurrency struct {
	CurrencyId    int64  `json:"currencyId"`    // 货币id
	Code          string `json:"code"`          // 代码
	DecimalLength int32  `json:"decimalLength"` // 小数长度
	Symbol        string `json:"symbol"`        // 符号
}

type Currency struct {
	Code          string  `json:"code"`          // 代码
	DecimalLength int32   `json:"decimalLength"` // 小数长度
	Symbol        string  `json:"symbol"`        // 符号
	Rate          float64 `json:"rate"`
}

type AmazonUser struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Account string   `json:"account"`
	AreaId  string   `json:"area_id"`
	OrgCode string   `json:"org_code"`
	Avatar  string   `json:"avatar"`
	Phone   string   `json:"phone"`
	Roles   []string `json:"roles"`
	Job     string   `json:"job"`
}
