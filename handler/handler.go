package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/easy-golang/foundation-framework/common/dto"
	"github.com/easy-golang/foundation-framework/common/security"
	"github.com/easy-golang/foundation-framework/constant/env"
	"github.com/easy-golang/foundation-framework/err"
	"github.com/easy-golang/foundation-framework/redis"
	"github.com/easy-golang/foundation-framework/util"
	"github.com/easy-golang/foundation-framework/util/md5"
	"github.com/easy-golang/foundation-framework/util/ustring"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/metadata"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	Authorization   = "Authorization"
	tokenPrefix     = "Bearer "
	userTokenPrefix = "user:token:"
)

func Response(ctx context.Context, w http.ResponseWriter, responseValue any, err error) {
	if err != nil {
		logx.Error(err)
		httpx.OkJsonCtx(ctx, w, ExceptionFail(err))
	} else {
		result, ok := responseValue.(standardResult)
		if ok {
			httpx.OkJsonCtx(ctx, w, result)
			return
		}
		httpx.OkJsonCtx(ctx, w, OK(responseValue))
	}
}

func ResponseOriginal(ctx context.Context, w http.ResponseWriter, responseValue any, err error) {
	if err != nil {
		logx.Error(err)
		httpx.OkJsonCtx(ctx, w, ExceptionFail(err))
	} else {
		httpx.OkJsonCtx(ctx, w, responseValue)
	}
}

func GlobalErrorMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer util.Recover(func(e error) {
			httpx.OkJsonCtx(r.Context(), w, standardResult{
				State:   false,
				Message: err.ServerError,
				Code:    500,
			})
		})
		next(w, r)
	}
}

type RequestContextMiddleware struct {
}

func NewRequestContextMiddleware() *RequestContextMiddleware {
	return &RequestContextMiddleware{}
}

func (*RequestContextMiddleware) Resolve(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestContextStr := r.Header.Get(security.RequestContextHearder)
		logx.Debugf("请求头上下文：{}")
		if ustring.IsNotEmptyString(requestContextStr) {
			requestContext := new(dto.RequestContext)
			unescapeContextStr, err := url.QueryUnescape(requestContextStr)
			if err != nil {
				logx.Errorf("请求上下文解码失败：%+v", err)
				next(w, r)
				return
			}
			err = json.Unmarshal([]byte(unescapeContextStr), requestContext)
			if err != nil {
				logx.Errorf("请求上下文Json解析失败：%+v", err)
				next(w, r)
				return
			}
			context := context.WithValue(r.Context(), security.RequestContextKey, requestContext)
			context = metadata.NewOutgoingContext(context, metadata.Pairs(security.RequestContextKey, requestContextStr))
			next(w, r.WithContext(context))
			return
		}
		next(w, r)
	}
}

type AmazonUserMiddleware struct {
	redisClient redis.Client
	env         env.Env
}

func NewAmazonUserMiddleware(redisClient redis.Client, env env.Env) *AmazonUserMiddleware {
	return &AmazonUserMiddleware{
		redisClient: redisClient,
		env:         env,
	}
}

func (a *AmazonUserMiddleware) Resolve(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 非生产环境下用户写死
		if a.env != env.PROD {
			amazonUser := &dto.AmazonUser{Id: "632840", Name: "余俊杰"}
			amazonUserJson, _ := json.Marshal(amazonUser)
			context := context.WithValue(r.Context(), security.RequestContextKey, amazonUser)
			context = metadata.NewOutgoingContext(context, metadata.Pairs(security.RequestContextKey,
				url.QueryEscape(string(amazonUserJson))))
			next(w, r.WithContext(context))
			return
		}
		token := r.Header.Get(Authorization)
		if !strings.HasPrefix(token, tokenPrefix) {
			httpx.WriteJson(w, 401, AuthorizationHint{
				code:    401,
				message: "无效的Authorization",
			})
			return
		}

		token = token[len(tokenPrefix):]
		if len(token) == 0 {
			httpx.WriteJson(w, 401, AuthorizationHint{
				code:    401,
				message: "无效的Authorization",
			})
			return
		}
		// 根据token获取用户信息
		amazonUserJson, err := a.redisClient.Get(userTokenPrefix + token)
		if err != nil {
			logx.Error("redis获取用户信息失败", err)
		}

		amazonUser := new(dto.AmazonUser)
		if len(amazonUserJson) != 0 {
			err := json.Unmarshal([]byte(amazonUserJson), amazonUser)
			if err != nil {
				logx.Error(err)
			}
		} else {
			// http请求
			lock := a.redisClient.NewRedisLock(userTokenPrefix + token)
			lock.LockFunc(func() error {
				// 根据token获取用户信息
				amazonUserJson, err = a.redisClient.Get(userTokenPrefix + token)
				if err != nil {
					logx.Error("redis获取用户信息失败", err)
				}
				if len(amazonUserJson) == 0 {
					// http请求用户信息
					amazonUser, err = getUserInfo(token)
					if err != nil || amazonUser == nil {
						httpx.WriteJson(w, 401, AuthorizationHint{
							code:    401,
							message: "无效的Authorization",
						})
						return nil
					}
					amazonUserByte, err := json.Marshal(amazonUser)
					if err != nil {
						logx.Error(err)
					}
					amazonUserJson = string(amazonUserByte)
					// 用户信息存入redis
					a.redisClient.Set(userTokenPrefix+token, amazonUserJson)
				}
				return nil
			})
		}
		context := context.WithValue(r.Context(), security.RequestContextKey, amazonUser)
		context = metadata.NewOutgoingContext(context, metadata.Pairs(security.RequestContextKey,
			url.QueryEscape(amazonUserJson)))
		next(w, r.WithContext(context))
	}
}

func getUserInfo(token string) (*dto.AmazonUser, error) {
	// 准备表单数据
	formData := url.Values{
		"_TOKEN": {token},
	}
	requestBody := bytes.NewBufferString(formData.Encode())
	request, err := http.NewRequest(http.MethodPost, "https://gateway.giikin.cn/gsso/sso/getUserInfo", requestBody)
	if err != nil {
		return nil, err
	}
	request.Header.Set("token", token)
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	request.Header.Set("timestamp", timestamp)
	request.Header.Set("sign", md5.Encode(timestamp+token))
	response, err := http.DefaultClient.Do(request)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {

	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var result map[string]any
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	result = result["data"].(map[string]any)
	var roles []string
	if result != nil {
		for _, role := range result["roles"].([]any) {
			roles = append(roles, role.(string))
		}
	}

	amazonUser := &dto.AmazonUser{
		Id:      result["id"].(string),
		Name:    result["name"].(string),
		Account: result["account"].(string),
		AreaId:  result["area_id"].(string),
		OrgCode: result["org_code"].(string),
		Avatar:  result["avatar"].(string),
		Phone:   result["phone"].(string),
		Roles:   roles,
		Job:     result["job"].(string),
	}
	return amazonUser, nil
}

type AuthorizationHint struct {
	code    int
	message string
}
