package shenyu

import (
	"github.com/apache/shenyu-client-golang/clients/admin_client"
	"github.com/apache/shenyu-client-golang/clients/http_client"
	"github.com/apache/shenyu-client-golang/common/constants"
	"github.com/apache/shenyu-client-golang/common/shenyu_error"
	"github.com/apache/shenyu-client-golang/model"
	"github.com/wangliujing/foundation-framework/err/comm"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/netx"
	"github.com/zeromicro/go-zero/rest"
	"strconv"
)

type Conf struct {
	UserName string `json:",default=admin"`
	Password string `json:",default=123456"`
	AdminUrl string `json:",default=http://127.0.0.1:9095"`
	Ip       string `json:",optional"`
}

func GetAdminToken(conf Conf) (adminToken model.AdminToken, err error) {
	headers := map[string][]string{}
	headers[constants.DEFAULT_CONNECTION] = []string{constants.DEFAULT_CONNECTION_VALUE}
	headers[constants.DEFAULT_CONTENT_TYPE] = []string{constants.DEFAULT_CONTENT_TYPE_VALUE}

	params := map[string]string{}
	params[constants.ADMIN_USERNAME] = conf.UserName
	params[constants.ADMIN_PASSWORD] = conf.Password
	//tokenRequest := initShenYuCommonRequest(headers, params, constants.DEFAULT_SHENYU_TOKEN, "token")
	tokenRequest := &model.ShenYuCommonRequest{
		Url:       conf.AdminUrl + constants.DEFAULT_SHENYU_TOKEN,
		Header:    headers,
		Params:    params,
		TimeoutMs: constants.DEFAULT_REQUEST_TIME,
	}

	adminToken, err = admin_client.GetShenYuAdminUser(tokenRequest)
	if err == nil {
		return adminToken, nil
	} else {
		return model.AdminToken{}, err
	}
}

func RegisterMetaData(adminTokenData model.AdminTokenData, metaData *model.MetaDataRegister, conf Conf) (registerResult bool, err error) {
	headers := adapterHeaders(adminTokenData)

	params := map[string]string{}
	params["appName"] = metaData.AppName
	params["path"] = metaData.Path
	params["contextPath"] = metaData.ContextPath
	params["host"] = metaData.Host
	params["port"] = metaData.Port

	if metaData.RPCType != "" {
		params["rpcType"] = metaData.RPCType
	} else {
		params["rpcType"] = constants.RPCTYPE_HTTP
	}

	if metaData.RuleName != "" {
		params["ruleName"] = metaData.RuleName
	} else {
		params["ruleName"] = metaData.Path
	}
	//tokenRequest := initShenYuCommonRequest(headers, params, constants.REGISTER_METADATA, "")
	tokenRequest := &model.ShenYuCommonRequest{
		Url:       conf.AdminUrl + constants.DEFAULT_BASE_PATH + constants.REGISTER_METADATA,
		Header:    headers,
		Params:    params,
		TimeoutMs: constants.DEFAULT_REQUEST_TIME,
	}

	registerResult, err = http_client.RegisterMetaData(tokenRequest)
	if err == nil {
		return registerResult, nil
	} else {
		return false, err
	}
}

func UrlRegister(adminTokenData model.AdminTokenData, urlMetaData *model.URIRegister, conf Conf) (registerResult bool, err error) {
	headers := adapterHeaders(adminTokenData)

	params := map[string]string{}
	if urlMetaData.AppName == "" || urlMetaData.RPCType == "" || urlMetaData.Host == "" || urlMetaData.Port == "" {
		return false, shenyu_error.NewShenYuError(constants.MISS_PARAM_ERROR_CODE, constants.MISS_PARAM_ERROR_MSG, err)
	}
	params["protocol"] = urlMetaData.Protocol
	params["appName"] = urlMetaData.AppName
	params["contextPath"] = urlMetaData.ContextPath
	params["host"] = urlMetaData.Host
	params["port"] = urlMetaData.Port
	params["rpcType"] = urlMetaData.RPCType

	tokenRequest := &model.ShenYuCommonRequest{
		Url:       conf.AdminUrl + constants.DEFAULT_BASE_PATH + constants.REGISTER_URI,
		Header:    headers,
		Params:    params,
		TimeoutMs: constants.DEFAULT_REQUEST_TIME,
	}

	registerResult, err = http_client.DoUrlRegister(tokenRequest)
	if err == nil {
		return registerResult, nil
	} else {
		return false, err
	}
}

func adapterHeaders(adminTokenData model.AdminTokenData) map[string][]string {
	headers := map[string][]string{}
	headers[constants.DEFAULT_CONNECTION] = []string{constants.DEFAULT_CONNECTION_VALUE}
	headers[constants.DEFAULT_CONTENT_TYPE] = []string{constants.DEFAULT_CONTENT_TYPE_VALUE}
	headers[constants.DEFAULT_TOKEN_HEADER_KEY] = []string{adminTokenData.Token}
	return headers
}

type Server interface {
	AddRoutes(rs []rest.Route, opts ...rest.RouteOption)
}

type SpringCloudServer struct {
	*rest.Server
	conf    Conf
	token   model.AdminToken
	appName string
}

func NewSpringCloudServer(server *rest.Server, conf Conf, restConf rest.RestConf) *SpringCloudServer {
	// 注册服务到注册中心
	token, err := GetAdminToken(conf)
	if err != nil {
		logx.Must(err)
	}
	return &SpringCloudServer{
		Server:  server,
		conf:    conf,
		appName: restConf.Name,
		token:   token,
	}
}

func (s *SpringCloudServer) AddRoutes(rs []rest.Route, opts ...rest.RouteOption) {
	s.Server.AddRoutes(rs, opts...)
	routes := s.Server.Routes()

	for _, route := range routes {
		index := nextSlash(route.Path)
		metaDataRegister := model.MetaDataRegister{
			AppName:     s.appName,
			ContextPath: route.Path[:index],
			Enabled:     true,
			Path:        route.Path,
			RPCType:     "springCloud",
			RuleName:    route.Path,
		}
		_, err := RegisterMetaData(s.token.AdminTokenData, &metaDataRegister, s.conf)
		if err != nil {
			logx.Must(err)
		}
	}
}

type HttpServer struct {
	*rest.Server
	conf          Conf
	token         model.AdminToken
	registService bool
	appName       string
	port          string
}

func NewHttpServer(server *rest.Server, conf Conf, restConf rest.RestConf) *HttpServer {
	token, err := GetAdminToken(conf)
	if err != nil {
		logx.Must(err)
	}
	return &HttpServer{
		Server:        server,
		registService: true,
		conf:          conf,
		token:         token,
		appName:       restConf.Name,
		port:          strconv.Itoa(restConf.Port),
	}
}

func (s *HttpServer) AddRoutes(rs []rest.Route, opts ...rest.RouteOption) {
	s.Server.AddRoutes(rs, opts...)
	routes := s.Server.Routes()
	if s.registService {
		if len(s.conf.Ip) == 0 {
			s.conf.Ip = netx.InternalIp()
		}
		uriRegister := model.URIRegister{
			Protocol:    "http://",
			AppName:     s.appName,
			ContextPath: "/" + s.appName,
			RPCType:     "http",
			Host:        s.conf.Ip,
			Port:        s.port,
		}
		_, err := UrlRegister(s.token.AdminTokenData, &uriRegister, s.conf)
		if err != nil {
			logx.Must(err)
		}
		s.registService = false
	}
	for _, route := range routes {
		index := nextSlash(route.Path)
		if "/"+s.appName != route.Path[:index] {
			logx.Must(comm.New("regist path must start with /" + s.appName))
		}
		metaDataRegister := model.MetaDataRegister{
			AppName:          s.appName,
			Enabled:          true,
			Path:             route.Path,
			RPCType:          "http",
			RuleName:         route.Path,
			RegisterMetaData: true,
		}
		_, err := RegisterMetaData(s.token.AdminTokenData, &metaDataRegister, s.conf)
		if err != nil {
			logx.Must(err)
		}
	}
}

func nextSlash(s string) int {
	for i := 1; i < len(s); i++ {
		if s[i] == '/' {
			return i
		}
	}
	logx.Must(comm.New("regist error path"))
	return -1
}
