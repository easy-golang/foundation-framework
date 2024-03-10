package security

import (
	"context"
	"encoding/json"
	"github.com/wangliujing/foundation-framework/common/dto"
	"github.com/wangliujing/foundation-framework/constant/symbol"
	"github.com/wangliujing/foundation-framework/err/comm"
	"github.com/wangliujing/foundation-framework/redis"
	"github.com/wangliujing/foundation-framework/util/ustring"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
	"net/url"
)

const (
	RequestContextHearder = "x-request-context"
	RequestContextKey     = RequestContextHearder
	shopSettingsKey       = "ShopSettings"
)

type giimallContext struct {
	context.Context
}

func NewGiimallContext(ctx context.Context) *giimallContext {
	return &giimallContext{
		Context: ctx,
	}
}

func (s *giimallContext) GetRequestContext() *dto.RequestContext {
	requestContext := s.Value(RequestContextKey)
	if requestContext != nil {
		return requestContext.(*dto.RequestContext)
	}
	md, ok := metadata.FromIncomingContext(s.Context)
	if ok {
		values := md[RequestContextKey]
		if values == nil || len(values) == 0 {
			return nil
		}
		requestContextStr, err := url.QueryUnescape(values[0])
		if err != nil {
			logx.Error(comm.Wrap(err))
			return nil
		}
		if ustring.IsEmptyString(requestContextStr) {
			return nil
		}
		requestContext := new(dto.RequestContext)
		err = json.Unmarshal([]byte(requestContextStr), requestContext)
		if err != nil {
			logx.Errorf("请求上下文解析失败：%+v", err)
			return nil
		}
		// 这里做了一层优化，连续调用可以不用重新解析json
		s.Context = context.WithValue(s.Context, RequestContextKey, requestContext)
		return requestContext
	}
	return nil
}

func (s *giimallContext) GetXRequestContext() *dto.XReqeustContext {
	requestContext := s.GetRequestContext()
	if requestContext != nil {
		return &requestContext.XReqeustContext
	}
	return nil
}

func (s *giimallContext) GetRequestId() string {
	requestContext := s.GetRequestContext()
	if requestContext != nil {
		return requestContext.XRequestId
	}
	return symbol.EmptyString
}

func (s *giimallContext) GetShopSettings() *dto.ShopSettings {
	value := s.Value(shopSettingsKey)
	if value != nil {
		return value.(*dto.ShopSettings)
	}
	xRequestContext := s.GetXRequestContext()
	if xRequestContext == nil || ustring.IsEmptyString(xRequestContext.ShopSettings) {
		return nil
	}
	shopSettings := new(dto.ShopSettings)
	err := json.Unmarshal([]byte(xRequestContext.ShopSettings), shopSettings)
	if err != nil {
		logx.Error(comm.Wrap(err))
		return nil
	}
	// 这里做了一层优化，连续调用可以不用重新解析json
	s.Context = context.WithValue(s.Context, shopSettingsKey, shopSettings)
	return shopSettings
}

func (s *giimallContext) GetLoginUser() *dto.UserInfo {
	xRequestContext := s.GetXRequestContext()
	if xRequestContext == nil {
		return nil
	}
	return xRequestContext.UserInfo
}

type amazonContext struct {
	context.Context
	redisClient redis.Client
}

func NewAmazonContext(ctx context.Context, client redis.Client) *amazonContext {
	return &amazonContext{
		Context:     ctx,
		redisClient: client,
	}
}

func (s *amazonContext) GetLoginUser() *dto.AmazonUser {
	amazonUser := s.Value(RequestContextKey)
	if amazonUser != nil {
		return amazonUser.(*dto.AmazonUser)
	}
	md, ok := metadata.FromIncomingContext(s.Context)
	if ok {
		values := md[RequestContextKey]
		if values == nil || len(values) == 0 {
			return nil
		}
		amazonUserJson, err := url.QueryUnescape(values[0])
		if err != nil {
			logx.Error(comm.Wrap(err))
			return nil
		}
		if len(amazonUserJson) == 0 {
			return nil
		}
		amazonUser := new(dto.AmazonUser)
		err = json.Unmarshal([]byte(amazonUserJson), amazonUser)
		if err != nil {
			logx.Error("用户信息解析失败：%+v", err)
			return nil
		}
		s.Context = context.WithValue(s.Context, RequestContextKey, amazonUser)
		return amazonUser
	}
	return nil
}
