# 基础框架

***缩短从需求到上线的距离***
<p align="center">
<img align="center" width="150px" src="static/img/go.jpeg">
</p>


![Go](https://github.com/zeromicro/go-zero/workflows/Go/badge.svg?branch=master)
![Go Report Card](https://goreportcard.com/badge/github.com/zeromicro/go-zero)
![codecov](https://codecov.io/gh/zeromicro/go-zero/branch/master/graph/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/zeromicro/go-zero.svg)](https://pkg.go.dev/github.com/zeromicro/go-zero)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

基础框架是基于go-zero微服务框架以及公司的业务特点、中间件环境构建的一套快速开发工具集，包含对接网关、rabbitMq、redis、orm、consul、nacos、登录用户、通用常量、JsonRpc服务发布、JsonRpc服务调用等扩展组件,以及定义了一些标准struct来规范开发，比如标准响应体，标准异常、标准分页等

基础框架未对go-zero做任何修改，只是在go-zero的基础上做了很多方便业务开发的组件抽象以及扩展，方便业务进行代码复用，更加专注业务层进行代码开发，所以为了更好的使用基础框架必须先了解 [go-zero](https://go-zero.dev/docs)

## 目录
* [1 为什么选用go-zero作为微服务框架](#1为什么选用go-zero作为微服务框架)
* [2 代码生成器](#2代码生成器)
* [3 接入shenyu网关](#3接入shenyu网关)
* [4 从上下文获取请求信息](#4从上下文获取请求信息)
* [5 自定义异常](#5自定义异常)
* [6 全局异常处理器](#6全局异常处理器)
* [7 redis分布式锁](#7redis分布式锁)
* [8 分布式事物解决方案](#8分布式事物解决方案)
* [9 使用RabbitMq](#9使用RabbitMq)
* [10 GRPC与JsonRpc](#10GRPC与JsonRpc)
* [11 接入XXL-JOB](#11接入XXL-JOB)
* [11 接入XXL-JOB](#11接入XXL-JOB)
* [12 接入nacos配置中心](#12接入nacos配置中心)
* [13 服务上下文ServiceContext](#13服务上下文ServiceContext)

##  1为什么选用go-zero作为微服务框架
### 1.1 web框架派系
这个派系依旧秉承`不要让框架束缚开发`，但由于从标准库到可使用的web框架，仍旧需要一定量的开发工作，所以倾向于选择成熟的go web框架，比如gin。

所以出现了以`gin+grpc`为核心，将其他组件集成进来的微服务架构。

同样可以使用k8s+istio。

最终构建的仍旧是现代化的云原生微服务架构。

### 1.2 大一统框架
当然，并不是所有的团队都有能力去执行上述两种方式，因为go开发者大多数是从别的语言转过来的，而其他语言都有成熟的大一统框架进行支撑，比如`Java的spring-boot、spring-cloud`，那自然而然，希望go也有这样的框架。

虽然我个人也并不赞成这样的方式，但也有其优点，因为能很好的减轻工作量，达到快速开发的目的，代价就是遵循框架的规则。

go主流的微服务框架不是很多，比较火的有`go-zero（26k stars）`、`Kratos（21.6k stars）`、`go-micro（20.9k stars）`

### 1.2.1 go-micro
原本的老大，但最近变动很大，作者去搞云原生去了，分布式框架的这个分支不确定后续还能分得创作者多少精力，担心后续万一go语言有啥大改动，框架作者懒得跟上，问题就大条了
* go-micro 2.0 版本依然可以使用，但已经不维护
* go-micro 3.0 版已经转到个人名下（Micro 公司的 CEO Asim），为了避免混淆，项目名也改了，叫 Nitro[3]。现在它是一个快速的分布式应用程序开发框架
* go-micro 3.0 版中文文档较少
* 各版本兼容性差
* 国外开源项目，容易受国际局势影响，有安全隐患，会对后续使用造成重大隐患

### 1.2.2 kratos
学了几天，发现这个源码应该是写得不错的，但是文档教程部分似乎没那么重视。像是顺便把自己在工作中形成的好东西分享出来回馈社会，在社区方面就没太用心，造成文档教程品质和作品本身的品质相差悬殊。有些东西甚至是源码里变动了但官方文档还是旧的，独立自学的话，真的非常容易走弯路，耗时耗力。最靠谱的途径反而是跟着他们的官方例子学

### 1.2.3 go-zero
一款由饿了么点评公司开源的微服务框架，起之秀，感觉得出来作者非常活跃和积极，官方群都是作者亲自加微信求star再迎接入群的，诚意满满。官网的文档详细齐全。官方（作者本人）视频教程很少，但是看得出来准备很认真，有PPT（讲缓存设计的那个）。虽然起步较晚但整个社群非常活跃，github上的star量已经超过kratos
加上go-zero是一款比较关注稳定性和开发效率的框架，这些正是公司所需要的

go-zero最吸引我的一点是它的代码生成器，它不仅能自动生成http路由、grpc远程调用、数据库层代码，而且还能自动生成Dockerfile、kubernetes deployment 资源文件，简直眼前一亮，这是其他框架代码生成器不具备的


`所以最终我选择了go-zero来进行微服务开发`


## 2代码生成器
go-zero的自身的数据库框架个人感觉没有gorm做的好，而且goctl model生成的代码的确不够丰富，方法太少，所以为了集成gorm，并扩展数据库的接口，修改了原生的cli工具，当然修改代码生成器也是为了更好的迎合基础框架，代码已经上传公司代码库

[代码生成器源码地址](https://github.com/smallnest/go-web-framework-benchmark)
### 2.1 安装代码生成器
* [下载源码](http://github.com/wangliujing/goctl)
* 在go.mod文件所在的目录执行命令
```shell
go install
```
### 2.2 生成Go HTTP 代码
* 编写api文件`app.api`文件名可以任意，执行后面的命令的时候用自己的文件名就可以
```
syntax = "v1"

type (
	UserLoginReq struct{}
	UserInfoReq  {
		Id int64 `path:"id"`
	}
	UserLoginResp      struct{}
	UserInfoResp       struct{}
	UserInfoUpdateReq  struct{}
	UserInfoUpdateResp struct{}
)

type (
	UserRoleReq        struct{}
	UserRoleResp       struct{}
	UserRoleUpdateReq  struct{}
	UserRoleUpdateResp struct{}
	UserRoleAddReq     struct{}
	UserRoleAddResp    struct{}
	UserRoleDeleteReq  struct{}
	UserRoleDeleteResp struct{}
)

type (
	UserClassReq        struct{}
	UserClassResp       struct{}
	UserClassUpdateReq  struct{}
	UserClassUpdateResp struct{}
	UserClassAddReq     struct{}
	UserClassAddResp    struct{}
	UserClassDeleteReq  struct{}
	UserClassDeleteResp struct{}
)

@server (
	prefix: /appName
	group:  user
	//signature: true
	middleware: UserAgentMiddleware
)
service user-api {
	@handler UserLogin
	post /user/login (UserLoginReq) returns (UserLoginResp)

	@handler getUserInfo
	get /user/info/:id (UserInfoReq) returns (UserInfoResp)

	@handler UserInfoUpdate
	post /user/info/update (UserInfoUpdateReq) returns (UserInfoUpdateResp)

	@handler UserList
	get /user/list returns ([]UserInfoResp)
}

@server (
	prefix: /appName
	group:  role
)
service user-api {
	@handler UserRoleList
	get /user/role/list returns ([]UserRoleResp)

	@handler UserRoleUpdate
	get /user/role/update (UserRoleUpdateReq) returns (UserRoleUpdateResp)

	@handler UserRoleInfo
	get /user/role/info (UserRoleReq) returns (UserRoleResp)

	@handler UserRoleAdd
	get /user/role/add (UserRoleAddReq) returns (UserRoleAddResp)

	@handler UserRoleDelete
	get /user/role/delete (UserRoleDeleteReq) returns (UserRoleDeleteResp)
}

@server (
	prefix: /appName
	group:  class
)
service user-api {
	@handler UserClassList
	get /user/class/list returns ([]UserClassResp)

	@handler UserClassUpdate
	get /user/class/update (UserClassUpdateReq) returns (UserClassUpdateResp)

	@handler UserClassInfo
	get /user/class/info (UserClassReq) returns (UserClassResp)

	@handler UserClassAdd
	get /user/class/add (UserClassAddReq) returns (UserClassAddResp)

	@handler UserClassDelete
	get /user/class/delete (UserClassDeleteReq) returns (UserClassDeleteResp)
}
```
注意api文件中prefix必须是 `/ + 工程名` 因为后续要接入shengyu网关，需要工程名来作为路由键
* 执行命令自动生成代码
```shell
goctl api go --api app.api --dir ../code --style go_zero
```
### 2.3 生成grpc代码
* 编写.proto文件`deposit.proto`文件名可以任意，执行后面的命令的时候用自己的文件名就可以
```
syntax = "proto3";

package deposit;

option go_package = "./deposit";

message DepositRequest {
  float amount = 1;
}

message DepositResponse {
  bool ok = 1;
}

service DepositService {
  rpc Deposit(DepositRequest) returns (DepositResponse);
}
```
* 执行命令自动生成代码
```shell
goctl rpc protoc deposit.proto --go_out=../internal/grpc --go-grpc_out=../internal/grpc --zrpc_out=../code --style go_zero -m
```

### 2.4 gorm代码生成
* 执行命令自动生成代码
```shell
goctl model mysql datasource --url 'jmy:Qwer1234@tcp(rm-bp1k0fj4ej13saps7.mysql.rds.aliyuncs.com:3306)/aws_message' --table aws_msg --dir ./mapper --cache --style go_zero
```
model代码生成改动比较大，首先摒弃了go-zero原生的数据库框架，集成了gorm，并且自动生成了很多常用的方法，但是没有自动生成复杂的查询，首先复杂的查询语句需要根据业务需要定制，有同事可能会说那可以提供一个通用的查询，由用户传入不同条件自动查询出数据，这种方式的确可以提高开发效率，前期确实爽了，但是后期维护简直要命，因为这种方式无法直观的看到sql来进行调优，所以代码生成器只会自动生成一些不需要进行sql调优查询

示例代码 [`aws_msg_model.go`]
```go
// Code generated by goctl. DO NOT EDIT.

package mapper

import (
	"database/sql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/wangliujing/foundation-framework/common/dto"
	"github.com/wangliujing/foundation-framework/orm"
	"gorm.io/gorm/clause"
)

type (
	awsMsgMapper interface {
		// 保存数据，保存会自动填充CreatedAt
		Save(model *AwsMsgPO) (*string, error)
		// 根据ID删除数据，如果实体对象有gorm.DeletedAt类型字段，则逻辑删除，否者物理删除
		DeleteById(id string) error
		// 根据多个ID批量删除，如果实体对象有gorm.DeletedAt类型字段，则逻辑删除，否者物理删除
		DeleteByIds(ids ...string) error
		// 根据ID查询，如果实体对象有gorm.DeletedAt类型字段，不会查询出逻辑删除的数据
		FindById(id string) (*AwsMsgPO, error)
		// 根据ID查询，可以指定属性进行查询，提升查询效率，如果实体对象有gorm.DeletedAt类型字段，不会查询出逻辑删除的数据
		FindByIdWithField(fields []string, id string) (*AwsMsgPO, error)
		// 根据ID查询，如果实体对象有gorm.DeletedAt类型字段，查询的时候会忽略逻辑删除的条件
		FindByIdUnscoped(id string) (*AwsMsgPO, error)
		// 根据ID查询，可以指定属性进行查询，提升查询效率，如果实体对象有gorm.DeletedAt类型字段，查询的时候会忽略逻辑删除的条件
		FindByIdUnscopedWithField(fields []string, id string) (*AwsMsgPO, error)
		// 根据多个ID批量查询，如果实体对象有gorm.DeletedAt类型字段，不会查询出逻辑删除的数据
		FindByIds(ids ...string) ([]AwsMsgPO, error)
		// 根据多个ID批量查询，可以指定属性进行查询，提升查询效率，如果实体对象有gorm.DeletedAt类型字段，不会查询出逻辑删除的数据
		FindByIdsWithField(fields []string, ids ...string) ([]AwsMsgPO, error)
		// 根据多个ID批量查询，如果实体对象有gorm.DeletedAt类型字段，查询的时候会忽略逻辑删除的条件
		FindByIdsUnscoped(ids ...string) ([]AwsMsgPO, error)
		// 根据多个ID批量查询，可以指定属性进行查询，提升查询效率，如果实体对象有gorm.DeletedAt类型字段，查询的时候会忽略逻辑删除的条件
		FindByIdsUnscopedWithField(fields []string, ids ...string) ([]AwsMsgPO, error)
		// 查询前size条数据，如果实体对象有gorm.DeletedAt类型字段，不会查询出逻辑删除的数据
		List(size int64) ([]AwsMsgPO, error)
		// 查询前size条数据，如果实体对象有gorm.DeletedAt类型字段，查询的时候会忽略逻辑删除的条件
		ListUnscoped(size int64) ([]AwsMsgPO, error)
		// 查询前size条数据，可以指定属性进行查询，提升查询效率，如果实体对象有gorm.DeletedAt类型字段，不会查询出逻辑删除的数据
		ListWithField(fields []string, size int64) ([]AwsMsgPO, error)
		// 查询前size条数据，可以指定属性进行查询，提升查询效率，如果实体对象有gorm.DeletedAt类型字段，查询的时候会忽略逻辑删除的条件
		ListUnscopedWithField(fields []string, size int64) ([]AwsMsgPO, error)
		// 分页查询，如果实体对象有gorm.DeletedAt类型字段，不会查询出逻辑删除的数据
		Page(page int64, size int64) (*dto.Page[AwsMsgPO], error)
		// 分页查询，如果实体对象有gorm.DeletedAt类型字段，查询的时候会忽略逻辑删除的条件
		PageUnscoped(page int64, size int64) (*dto.Page[AwsMsgPO], error)
		// 分页查询，可以指定属性进行查询，提升查询效率，如果实体对象有gorm.DeletedAt类型字段，不会查询出逻辑删除的数据
		PageWithField(fields []string, page int64, size int64) (*dto.Page[AwsMsgPO], error)
		// 分页查询，可以指定属性进行查询，提升查询效率，如果实体对象有gorm.DeletedAt类型字段，查询的时候会忽略逻辑删除的条件
		PageUnscopedWithField(fields []string, page int64, size int64) (*dto.Page[AwsMsgPO], error)
	}

	defaultAwsMsgMapper struct {
		conn orm.Connection
	}

	AwsMsgPO struct {
		Id               string         `gorm:"column:id"`                // 主键
		SellerId         sql.NullString `gorm:"column:seller_id"`         // 卖家ID
		NotificationType sql.NullString `gorm:"column:notification_type"` // 事件类型
		EventTime        sql.NullTime   `gorm:"column:event_time"`        // 事件事件
		Md5              sql.NullString `gorm:"column:md5"`               // 消息体md5值
		Message          sql.NullString `gorm:"column:message"`           // 消息
		Body             sql.NullString `gorm:"column:body"`              // 消息体
		TriggerTime      sql.NullTime   `gorm:"column:trigger_time"`      // 触发时间
		CreatedAt        sql.NullTime   `gorm:"column:created_at"`        // 创建时间
		DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at"`        // 删除时间
	}
)

func (m *AwsMsgPO) TableName() string {
	return "aws_msg"
}

func newAwsMsgMapper(conn orm.Connection) *defaultAwsMsgMapper {
	return &defaultAwsMsgMapper{
		conn: conn,
	}
}

func (m *defaultAwsMsgMapper) Save(model *AwsMsgPO) (*string, error) {
	if err := m.conn.(*orm.GormConnection).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(model).Error; err != nil {
		return nil, err
	}
	id := model.Id
	return &id, nil
}
func (m *defaultAwsMsgMapper) DeleteById(id string) error {
	return m.conn.(*orm.GormConnection).Delete(&AwsMsgPO{
		Id: id,
	}).Error
}
func (m *defaultAwsMsgMapper) DeleteByIds(ids ...string) error {
	return m.conn.(*orm.GormConnection).Delete(&AwsMsgPO{}, ids).Error
}
func (m *defaultAwsMsgMapper) FindById(id string) (*AwsMsgPO, error) {
	result := &AwsMsgPO{}
	if err := m.conn.(*orm.GormConnection).Where("id=?", id).Find(result).Error; err != nil {
		if err == logger.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}
func (m *defaultAwsMsgMapper) FindByIdWithField(fields []string, id string) (*AwsMsgPO, error) {
	result := &AwsMsgPO{}
	if err := m.conn.(*orm.GormConnection).Select(fields).Where("id=?", id).Find(result).Error; err != nil {
		if err == logger.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}
func (m *defaultAwsMsgMapper) FindByIdUnscoped(id string) (*AwsMsgPO, error) {
	result := &AwsMsgPO{}
	if err := m.conn.(*orm.GormConnection).Unscoped().Where("id=?", id).Find(result).Error; err != nil {
		if err == logger.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}
func (m *defaultAwsMsgMapper) FindByIdUnscopedWithField(fields []string, id string) (*AwsMsgPO, error) {
	result := &AwsMsgPO{}
	if err := m.conn.(*orm.GormConnection).Unscoped().Select(fields).Where("id=?", id).Find(result).Error; err != nil {
		if err == logger.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}
func (m *defaultAwsMsgMapper) FindByIds(ids ...string) ([]AwsMsgPO, error) {
	result := make([]AwsMsgPO, 0)
	if err := m.conn.(*orm.GormConnection).Find(&result, ids).Error; err != nil {
		return nil, err
	}
	return result, nil
}
func (m *defaultAwsMsgMapper) FindByIdsWithField(fields []string, ids ...string) ([]AwsMsgPO, error) {
	result := make([]AwsMsgPO, 0)
	if err := m.conn.(*orm.GormConnection).Select(fields).Find(&result, ids).Error; err != nil {
		return nil, err
	}
	return result, nil
}
func (m *defaultAwsMsgMapper) FindByIdsUnscoped(ids ...string) ([]AwsMsgPO, error) {
	result := make([]AwsMsgPO, 0)
	if err := m.conn.(*orm.GormConnection).Unscoped().Find(&result, ids).Error; err != nil {
		return nil, err
	}
	return result, nil
}
func (m *defaultAwsMsgMapper) FindByIdsUnscopedWithField(fields []string, ids ...string) ([]AwsMsgPO, error) {
	result := make([]AwsMsgPO, 0)
	if err := m.conn.(*orm.GormConnection).Unscoped().Select(fields).Find(&result, ids).Error; err != nil {
		return nil, err
	}
	return result, nil
}
func (m *defaultAwsMsgMapper) List(size int64) ([]AwsMsgPO, error) {
	result := make([]AwsMsgPO, 0)
	if err := m.conn.(*orm.GormConnection).Limit(int(size)).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
func (m *defaultAwsMsgMapper) ListUnscoped(size int64) ([]AwsMsgPO, error) {
	result := make([]AwsMsgPO, 0)
	if err := m.conn.(*orm.GormConnection).Unscoped().Limit(int(size)).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
func (m *defaultAwsMsgMapper) ListWithField(fields []string, size int64) ([]AwsMsgPO, error) {
	result := make([]AwsMsgPO, 0)
	if err := m.conn.(*orm.GormConnection).Select(fields).Limit(int(size)).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
func (m *defaultAwsMsgMapper) ListUnscopedWithField(fields []string, size int64) ([]AwsMsgPO, error) {
	result := make([]AwsMsgPO, 0)
	if err := m.conn.(*orm.GormConnection).Unscoped().Select(fields).Limit(int(size)).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
func (m *defaultAwsMsgMapper) Page(page int64, size int64) (*dto.Page[AwsMsgPO], error) {
	var total int64
	if err := m.conn.(*orm.GormConnection).Model(&AwsMsgPO{}).Count(&total).Error; err != nil {
		return nil, err
	}
	if total == 0 {
		return dto.NewPage[AwsMsgPO](nil, 0, size), nil
	}
	result := make([]AwsMsgPO, 0)
	if err := m.conn.(*orm.GormConnection).Offset(int((page - 1) * size)).Limit(int(size)).Find(&result).Error; err != nil {
		return nil, err
	}
	return dto.NewPage[AwsMsgPO](result, total, size), nil
}
func (m *defaultAwsMsgMapper) PageUnscoped(page int64, size int64) (*dto.Page[AwsMsgPO], error) {
	var total int64
	if err := m.conn.(*orm.GormConnection).Unscoped().Model(&AwsMsgPO{}).Count(&total).Error; err != nil {
		return nil, err
	}
	if total == 0 {
		return dto.NewPage[AwsMsgPO](nil, 0, size), nil
	}
	result := make([]AwsMsgPO, 0)
	if err := m.conn.(*orm.GormConnection).Unscoped().Offset(int((page - 1) * size)).Limit(int(size)).Find(&result).
		Error; err != nil {
		return nil, err
	}
	return dto.NewPage[AwsMsgPO](result, total, size), nil
}
func (m *defaultAwsMsgMapper) PageWithField(fields []string, page int64, size int64) (*dto.Page[AwsMsgPO], error) {
	var total int64
	if err := m.conn.(*orm.GormConnection).Model(&AwsMsgPO{}).Count(&total).Error; err != nil {
		return nil, err
	}
	if total == 0 {
		return dto.NewPage[AwsMsgPO](nil, 0, size), nil
	}
	result := make([]AwsMsgPO, 0)
	if err := m.conn.(*orm.GormConnection).Select(fields).Offset(int((page - 1) * size)).Limit(int(size)).Find(&result).
		Error; err != nil {
		return nil, err
	}
	return dto.NewPage[AwsMsgPO](result, total, size), nil
}
func (m *defaultAwsMsgMapper) PageUnscopedWithField(fields []string, page int64, size int64) (*dto.Page[AwsMsgPO], error) {
	var total int64
	if err := m.conn.(*orm.GormConnection).Unscoped().Model(&AwsMsgPO{}).Count(&total).Error; err != nil {
		return nil, err
	}
	if total == 0 {
		return dto.NewPage[AwsMsgPO](nil, 0, size), nil
	}
	result := make([]AwsMsgPO, 0)
	if err := m.conn.(*orm.GormConnection).Unscoped().Select(fields).Offset(int((page - 1) * size)).Limit(int(size)).
		Find(&result).Error; err != nil {
		return nil, err
	}
	return dto.NewPage[AwsMsgPO](result, total, size), nil
}

```
示例代码 [`aws_msg_model_custom.go`]
```go
package mapper

import "github.com/wangliujing/foundation-framework/orm"

var _ AwsMsgMapper = (*customAwsMsgMapper)(nil)

type (
	AwsMsgMapper interface {
		awsMsgMapper
	}

	customAwsMsgMapper struct {
		*defaultAwsMsgMapper
	}
)

func NewAwsMsgMapper(conn orm.Connection) AwsMsgMapper {
	return &customAwsMsgMapper{
		defaultAwsMsgMapper: newAwsMsgMapper(conn),
	}
}

```
由以上代码可以看出AwsMsgMapper是集成awsMsgMapper的，所以说如果业务需要扩展操作数据库的方法只需要对AwsMsgMapper接口以及customAwsMsgMapper实现类进行扩展就可以了，不需要动aws_msg_model.go，减少了aws_msg_model.go这个文件出错的风险，需要注意的是如果需要查询逻辑删除的数据可以调用unscoped方法比如FindByIdUnscoped

以上只是对gorm进行了简单的封装，如需深入了解gorm请查看[gorm官方文档](https://gorm.io/zh_CN/docs/index.html)

## 3接入shenyu网关
我们首先看下未接入网关的，api文件自动生成的启动代码
```go
package main

import (
	"client/code/internal/config"
	"client/code/internal/handler"
	"client/code/internal/svc"
	"flag"
	"fmt"
	"github.com/wangliujing/foundation-framework/system"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/user-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	system.Start(server)
}

```
### 3.1 http的方式
配置文件
```text
Shenyu:
  AdminUrl: http://172.24.116.143:9095
  UserName: admin
  Password: 123456
```
配置类
```go
type Config struct {
    rest.RestConf
    Shenyu         shenyu.Conf
}

type Conf struct {
	UserName string `json:",default=admin"`
	Password string `json:",default=123456"`
	AdminUrl string `json:",default=http://127.0.0.1:9095"`
	Ip       string `json:",optional"`
}
```
关键代码修改
```go
handler.RegisterHandlers(server, ctx)
改为
handler.RegisterHandlers(shenyu.NewHttpServer(server, c.Shenyu, c.RestConf), ctx)
```
修改后重新启动就可以自动注册网关了

### 3.2 springCloud的方式
这种方式注册不像http方式需要把自己的ip注册到网关，只需要把自己的服务名告诉网关，网关自己根据服务名拿到服务的 `ip + port`，这种方式的好处是不会用到网关的注册中心，因为在某种情况下网关的注册中心可能与微服务的注册中心服务不一致，比如cousul的客户端挂掉，就会导致consul注册中心掉节点，而shenyu不会掉节点

### 3.2.1 nacos注册中心
配置文件
```text
Name: appName
Host: 0.0.0.0
Port: 8888
Timeout: -1
Shenyu:
  AdminUrl: http://172.24.116.143:9095
Nacos:
  Hosts:
    -
      Ip: 172.18.0.145
      Port: 8848
  Namespace: local_test
```

配置类
```go
type Config struct {
    rest.RestConf
    Shenyu         shenyu.Conf
    Nacos          nacos.Conf
}
```

示例代码
```go
package main

import (
	"flag"
	"fmt"
	"github.com/wangliujing/foundation-framework/reg/nacos"
	"github.com/wangliujing/foundation-framework/shenyu"
	"github.com/wangliujing/foundation-framework/system"

	"client/internal/config"
	"client/internal/handler"
	"client/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/user-api.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()
	ctx := svc.NewServiceContext(c)
	// springcloud注册关键代码
	handler.RegisterHandlers(shenyu.NewSpringCloudServer(server, c.Shenyu, c.RestConf), ctx)
	// 注册服务到注册中心，因为网关需要从注册中心获取服务信息
	nacos.NewRestRegistry(c.Nacos, c.RestConf).RegisterService()
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	system.Start(server)
}
```
### 3.2.2 consul注册中心
配置文件
```text
Name: appName
Host: 0.0.0.0
Port: 8888
Timeout: -1
Consul:
  Host: 172.30.20.23:8500
  Key: appName
Shenyu:
  AdminUrl: http://172.24.116.143:9095    
```
配置类
```go
type Config struct {
	rest.RestConf
	Shenyu         shenyu.Conf
	Consul         consul.Conf
}
```
示例代码
```go
package main

import (
	"flag"
	"fmt"
	"github.com/wangliujing/foundation-framework/reg/consul"
	"github.com/wangliujing/foundation-framework/shenyu"

	"client/internal/config"
	"client/internal/handler"
	"client/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	_ "github.com/zeromicro/zero-contrib/zrpc/registry/consul"
)

var configFile = flag.String("f", "etc/user-api.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()
	ctx := svc.NewServiceContext(c)
	// 注册shengyu关键代码
	handler.RegisterHandlers(shenyu.NewSpringCloudServer(server, c.Shenyu, c.RestConf), ctx)
	// 注册consul注册中心，主要是因为springcloud方式注册网关会从注册中心获取服务信息
	consul.NewRestRegistry(c.Consul, c.RestConf).RegisterService()
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

```

## 4从上下文获取请求信息
获取当前请求上下文需解析网关透传过来的信息（token、用户信息、请求上下文等），而且在rpc请求时候需要把上下文信息进行远程传输，在服务端需重新解析rpc上下文，为了屏蔽网关解析与rpc解析的差异，基础框架针对这块也做了封装，开发人员无需关心差异，直接无脑调用方法获取

### 4.1 吉喵云业务
### 4.1.1 在路由上加上上下文解析中间件
```json
{
    "x-request-context":{
        "acceptCurrency":"SAR",
        "acceptLanguage":"en",
        "clientIp":"8.210.217.33",
        "isAdmin":false,
        "loginId":"",
        "requestClient":"buyer",
        "shopId":"11121914443489828",
        "shopSettings":"{\"status\":true,\"rights_status\":1,\"validTime\":\"2037-07-15 15:59:59\",\"defaultLanguageMode\":2,\"mainCurrency\":{\"currencyId\":214,\"code\":\"SAR\",\"symbol\":\"﷼\",\"decimalLength\":2},\"mainLanguage\":\"ar\",\"language\":[\"ar\",\"en\"],\"morelanguage\":[\"ar\",\"zh-CN\",\"en\"],\"useMultipleCurrency\":1,\"useFixedExchangeRate\":1,\"currencyDisplayWay\":2,\"defaultDisplayWay\":2,\"currencyList\":[{\"code\":\"USD\",\"symbol\":\"$\",\"decimalLength\":2,\"rate\":0.2663},{\"code\":\"AED\",\"symbol\":\"د.إ\",\"decimalLength\":2,\"rate\":1}],\"timezone\":\"Asia\\/Shanghai\",\"plugins\":[1084,1007,1003,1005,1082,1006,1083,1086,1094,1089,1090,1087,1088,1091,1004,1095,1072,1008,1002,1001,1067,1011,1093,1085,1092,1009]}",
        "shopTimezone":"Asia/Shanghai",
        "userAgent":"mobile",
        "userRegion":"sg",
        "userInfo" : {
           "userId" : "123",
           "userName": "脏仲表（章仲淼）",
           "email": "123@qq.com",
           "status": "1",
           "preferredLanguage": "cn"
        }
    },
    "x-request-header":{
        "User-Agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
    },
    "x-request-id":"c21edeeb745a4d96b45b583b3f835b45.167.17026241022054781"
}
```
```go
import "github.com/wangliujing/foundation-framework/handler"

server.AddRoutes(
    rest.WithMiddlewares(
		// 上下文解析中间件，自动解析网关透传的上下文信息
        []rest.Middleware{handler.NewRequestContextMiddleware().Resolve},
        []rest.Route{
            {
                Method:  http.MethodGet,
                Path:    "/user/list",
                Handler: user.UserListHandler(serverCtx),
            },
        }...,
    ),
    rest.WithPrefix("/appName"),
)
```
### 4.1.2 使用关键代码
```go
context := security.NewGiimallContext(ctx)
// 获取请求上下文
requestContext := context.GetRequestContext()
xrequestContext := context.GetXRequestContext()
// 获取登录用户
loginUser := context.GetLoginUser()
// 获取请求ID
requestId := context.GetRequestId()
// 获取店铺设置
shopSettings := context.GetShopSettings()
```

### 4.1.3 获取当前登录用户示例代码
```go
func (s *DepositServiceServer) Deposit(ctx context.Context, in *mock.DepositRequest) (*mock.DepositResponse, error) {
    context := security.NewGiimallContext(ctx)
    loginUser := context.GetLoginUser()
    return &mock.DepositResponse{}, nil
}
```

### 4.2 亚马逊业务

### 4.2.1 在路由上加上用户信息解析中间件
```go
import "github.com/wangliujing/foundation-framework/handler"

func RegisterHandlers(server shenyu.Server, serverCtx *svc.ServiceContext) {
   server.AddRoutes(
       rest.WithMiddlewares(
           // 把token信息解析到上下文，方便后续通过上下文获取登录用户信息
           []rest.Middleware{handler.NewAmazonUserMiddleware(serverCtx.RedisClient, serverCtx.Env).Resolve},
           []rest.Route{
               {
                   Method:  http.MethodGet,
                   Path:    "/user/list",
                   Handler: user.UserListHandler(serverCtx),
               },
           }...,
       ),
       rest.WithPrefix("/appName"),
   )
}
```

### 4.2.2 使用关键代码
```go
context := security.NewAmazonContext(ctx)
// 获取登录用户
amazonUser := context.GetLoginUser()
```

### 4.2.3 获取当前登录用户示例代码
```go
func (s *DepositServiceServer) Deposit(ctx context.Context, in *mock.DepositRequest) (*mock.DepositResponse, error) {
    context := security.NewAmazonContext(ctx)
    amazonUser := context.GetLoginUser()
    return &mock.DepositResponse{}, nil
}
```

## 5自定义异常
### 5.1 背景
go中对于error的处理大概是被开发者吐槽最多的之一了,我刚接触这种写法时也是不太理解,代码中充斥着大量的if err != nil { return err }这种代码,看起来很不优雅.

我认为要想完全理解这种设计,首先要去接受和相信存在即合理的道理,而不是一味的吐槽.golang作为一种现代化的语言,从google内部演化而来,必然是吸取了前人的各种经验和优势,并加以精简才凝练出来的.

在php和java等众多语言中,对于错误的处理是exception和handler,exception在返回错误的同时还会携带堆栈信息,然后由开发者手动try{}cache{}捕获或者全局handler处理.这种对于error的处理使得代码非常简洁,而且输出到日志中时也可以轻易的定位问题.

然而exception本身也存在一定的隐患, 比如当我们在调用一个可能会抛出exception的方法时, 我们并不一定知道它是否会抛出异常以及如何处理这个异常, 因为它很可能是别人写的或者是一个第三方包. 必须有人告诉或者查阅文档才可以知道. 对于开发者来说, 就有可能没有及时的处理掉exception而得出了其他的结果.

在go语言中, 它允许你在函数和方法中返回多个结果, 如果有error就要立即返回. 而开发者在发现这个error时,也需要立即处理,是返回上级还是在当前逻辑中处理. 它要求一旦遇到error, 就需要立即的处理它, 操作权完全交由开发者自己去控制, 这就是go语言的设计哲学.

不过我要吐槽一点是go语言的错误没有异常栈信息，排查问题相当难受，针对这一点，基础框架自定义了两种带异常栈信息的异常类型，`通用异常comm.Error, 业务异常biz.Error`

### 5.2 通用异常`comm.Error`
通用异常可以在任意地方使用，目的只是为了让异常信息具有栈信息，辅助更快的定位错误

我们可以看一下官方的错误
```go
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/wangliujing/foundation-framework/common/security"
	"service/internal/config"
)

var configFile = flag.String("f", "etc/deposit.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.MustSetup(c.Log)
	logx.Error(errors.New("测试错误"))
}
```
输出内容
```text
API server listening at: 127.0.0.1:64299
2023-10-20T16:42:19.596+08:00    error  测试错误        caller=server/main.go:22

Debugger finished with the exit code 0

```
可以看出输出内容非常简洁，没有错误栈信息，如果在复杂调用场景排查错误并不容易

修改代码使用基础框架自定义错误
```go
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/wangliujing/foundation-framework/common/security"
	"github.com/wangliujing/foundation-framework/err/comm"
	"service/internal/config"
)

var configFile = flag.String("f", "etc/deposit.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.MustSetup(c.Log)

	err := errors.New("测试错误")
	logx.Error(comm.Wrap(err))
}
```
输出内容
```go
API server listening at: 127.0.0.1:64509
2023-10-20T16:46:05.482+08:00    error  github.com/wangliujing/foundation-framework/err/comm.Error:测试错误
    at D:/Program Files/Go/IdeaProject/gozero/server/main.go ( main.go:25 )
    at D:/Program Files/Go/go1.20.7/src/runtime/proc.go ( proc.go:250 )
    at D:/Program Files/Go/go1.20.7/src/runtime/asm_amd64.s ( asm_amd64.s:1598 )
        caller=server/main.go:25

Debugger finished with the exit code 0
```
可以看到有错误信息了，由于本人的主语言是java所以错误栈格式是仿照java写的，反正个人看起来是非常习惯的，哈哈

其他创建异常的方式
```go
// 创建一个通用错误
comm.New("测试错误")
// 创建一个带有错误码的通用错误，适用于需要根据错误码做业务判断的场景
comm.NewCode(101, "测试错误")
// 适用于需要字符串拼接的场景，采用占位符的方式代替字符串拼接
comm.Newf("错误信息是: %s", "测试错误")

......更多错误使用方式可以自己去尝试,由于太简单这里不一一说明了
```

### 5.3 业务异常`biz.Error`
业务异常，跟通用异常代码和功能是一样的，唯一的区别就是返回业务异常的时候可以把错误信息带给前端

比如api路由接收到一个业务异常
```go
 biz.New("测试错误")
```
则基础框架会吧这个异常自动处理成标准响应体
```json
{
    "state": false,
    "message": "测试错误",
    "data": null,
    "code": 500
}
```
除业务异常外的其他异常，因为不需要业务提示，所以统一返回
```json
{
    "state": false,
    "message": "抱歉，系统繁忙，请稍后重试！",
    "data": null,
    "code": 500
}
```

注意：
```go
comm.Wrap(biz.New("测试错误"))
```
上面这个代码虽然是通用异常，但还是会返回
```json
{
    "state": false,
    "message": "测试错误",
    "data": null,
    "code": 500
}
```
也就是只要错误类型是业务类型，或者错误类型只要包含有业务类型，则统一返回业务提示给前端

## 6全局异常处理器
全局异常处理器方便处理panic不可控因素导致的异常，统一异常响应体，给出友好提示，如果不加全局异常处理，接口出现panic则不返回任何结果，Http状态码为500

关键代码
```go
server.AddRoutes(
    rest.WithMiddlewares(
        // 路由增加全局异常处理器
        []rest.Middleware{handler.GlobalErrorMiddleware},
        []rest.Route{
            {
                Method:  http.MethodPost,
                Path:    "/user/login",
                Handler: user.UserLoginHandler(serverCtx),
            }
        }...,
    ),
    rest.WithPrefix("/appName"),
)
```

加全局异常处理器返回结果示例
```text
http 状态码200

{
    "state": false,
    "message": "抱歉，系统繁忙，请稍后重试！",
    "data": null,
    "code": 500
}
```

## 7redis分布式锁
redis分布式锁使用的时候有个锁过期的问题，很多开发人员使用的时候都会根据业务执行的时长来预估锁过期时间，比如业务执行3秒，锁过期可能稍微设置长一点比如10秒，如果业务在业务繁忙的时候可能执行时长会超过10秒，这个时候锁自动过期提前释放锁很可能就会造成业务数据不一致问题，有人可能会说那就设置1分钟甚至更长，其实这时候设置多长时间都不靠谱，你根本无法预知在服务器资源紧张的时候你的业务执行会耗费多长时间，合适的做法就是使用锁续期机制，如果业务没有执行完，会提前对即将过期的锁进行重置过期时间

go-zero提供的分布式锁，没有提供锁续期的机制，为防止业务出现问题，开发人员需自己实现锁续期，比较麻烦，且单从[官方文档](https://go-zero.dev/docs/tutorials/redis/redis-lock)使用方式来看，还是比较繁琐的，针对这些问题基础框架对go-zero的使用方式做了封装，用户不用考虑业务没有执行完锁过期的问题，基础框架redis获取锁成功后会在锁过期的前10秒重置过期时间，直到业务执行完成锁被释放

配置文件
```text
Redis:
  Host: 172.24.116.143:6379
  LockPath: "appName" # 一般使用工程名
```

配置类：
```go
type Config struct {
	Redis          redis.Conf
}
```
使用示例代码:

```go
// 导入基础框架redis包
improt "github.com/wangliujing/foundation-framework/redis"

redisClient := redis.NewClient(c.Redis)
lock := redisClient.NewRedisLock("lockName")
// 这里获取不到锁会一直等待直到获取到锁
if lock.Lock() {
    defer lock.UnLock()
    // 执行业务代码
}

// 这里是上面方式的简化用法，如果获取不到锁会一直等待
redisClient := redis.NewClient(c.Redis)
lock := redisClient.NewRedisLock("lockName")
err := lock.LockFunc(func() error {
// 执行业务代码
return nil
})
```

```go
// 导入基础框架redis包
improt "github.com/wangliujing/foundation-framework/redis"

redisClient := redis.NewClient(c.Redis)
lock := redisClient.NewRedisLock("lockName")
timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
defer cancelFunc()
// 这里获取不到锁，会在指定时间内持续尝试获取，如果指定时间内还是获取不到则不会继续尝试直接返回false，未获取到锁
if lock.LockWithContext(timeout) {
    defer lock.UnLock()
    // 执行业务代码
}

// 上面方式的简化写法
redisClient := redis.NewClient(c.Redis)
lock := redisClient.NewRedisLock("lockName")
timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
defer cancelFunc()
lock.LockFuncWithContext(timeout, func() error {
    // 执行业务代码
    return nil
})
```

```go
// 导入基础框架redis包
improt "github.com/wangliujing/foundation-framework/redis"

redisClient := redis.NewClient(c.Redis)
lock := redisClient.NewRedisLock("lockName")
// 这里只会尝试获取一次，获取不到直接返回false
if lock.TryLock() {
    defer lock.UnLock()
    // 执行业务代码
}

// 上面方式的简化写法
redisClient := redis.NewClient(c.Redis)
lock := redisClient.NewRedisLock("lockName")
lock.TryLockFunc(func() error {
    // 执行业务代码
    return nil
})

```

说明：`redis锁的格式 lock:&{LockPath}:&{lockName} 其中LockPath的值就是配置文件LockPath对应的值,锁名称lockName创建锁的时候由参数传入`

## 8分布式事物解决方案
分布式事物依据业务场景分为同步及异步，由于互联网业务场景普遍要求高性能高并发，XA事物有性能瓶颈所以一般不会采用XA事物，TCC事物虽然性能比XA事物好很多但是对业务侵入性太强，开发起来相对复杂，对性能要求苛刻的场景性能方面表现也相形见绌，综上一般会采用RocketMq的事物消息或者是RabbitMq本地消息表的最终一致方案。

为了兼容PHP消息中间件采用RabbitMq，所以异步方案选用RabbitMq本地消息表的最终一致方案

### 8.1 异步方案
    优点：把比较耗时的工作异步执行，性能有保证，可以应对高并发场景
    缺点：实现相对较复杂，数据一致性有延迟，无法应对实时性要求很高的场景
![异步方案](static/img/异步方案.png)
上图中，为了保证A系统与B系统的数据一致性采用消息最终一致的这种方案，只需要保证发送端与消费端消息不丢失，就可以保证数据一致性

发送端：发送端通过本地消息表、以及定时发送消息的方式来保证消息不会丢失，但是消息有可能会重复发送，首先本地消息表的数据是与本地事物在一个事物里面，业务执行成功事物提交后本地消息一定会有一条消息数据，如果没有执行成功事物回滚则本地消息表的数据也会回滚，提交事物成功后才发送消息，避免本地事物执行失败，消息发送成功导致数据不一致的问题，如果发送消息成功则RabbitMq会回调A系统接口修改本地消息为已发送状态，如果消息发送失败，则不会产生回调来修改状态，此时消息是未发送状态，定时任务会定时扫描未发送的消息进行重新发送，直到消息发送成功后回调A系统接口修改本地消息为已发送状态，所以发送端最终一定会保证消息被发送到RabbitMq端

消费端：消费端可以通过手动确认的方式来保证消息不丢失，其实RabbitMq自动确认也能保证消息不丢，只不过由框架来进行的自动确认，只要监听方法不出错就代表消息消费成功，则如果异步去消费必须采用手动确认的方式才能保证消息不丢，因为异步线程出错主线程无法感知，由于发送端无法保证消息唯一发送，则消费端根据业务情况来确定是否需要做幂等，大部分情况是需要进行幂等判断的

有人可能会说我把发送消息放在最后一行代码，也能最大限度的保证一致性，但是事物提交失败是有可能发生的虽然概率比较小，根据墨菲定律如果坏事有可能发生，不管可能性有多小，它肯定会发生，墨菲定律也是架构设计要遵循的定律之一

### 8.2 同步方案
    优点：实现比较简单，数据一致性实时一致没有延迟
    缺点：性能较差，无法应对高并发场景，如果系统优化到极致无法达到并发要求，只能通过堆机器的方式来进行解决，否者就放弃实时性（业务妥协的方式）采用异步方案
![同步方案](static/img/同步方案.png)
如上图下单扣减库存的场景，因为扣减库存需要校验库存是否充足，如果采用异步就会导致超卖的情况发生，所以这种情况只能选用同步的方式进行解决
我们分如下几个场景来探讨数据的一致性问题：
1. 订单系统执行成功，仓储系统执行失败
   如果仓储系统扣减库存失败或异常，因为是同步调用则订单的事物提交一定会失败，所以该场景不存在

2. 订单系统执行成功，仓储系统执行成功
   如果仓储系统与订单系统都执行成功，则没有数据一致性问题

3. 订单系统执行失败，仓储系统执行失败
   订单系统与仓储系统都执行失败，则没有数据一致性问题

4. 订单系统执行失败，仓储系统执行成功
   仓储系统执行成功，订单失败，此时需要业务可以对订单进行重试，重试的过程仓储扣减库存会重复调用，为了防止同一个订单重复扣减库存则仓储系统必须保证接口的幂等性，最终要保证订单成功，如果不能成功则可以取消订单或者人工手动处理的方式对库存数据进行回滚，如果是bug导致下单不成功则修复订单侧bug重试下单即可
   说到这里有人可能会提出为什么不用Seate的AT模式来解决，Seate的AT模式也可以，考虑到跨语言以及需要引入新的中间件，暂时基本采用上图那种简单的方式进行解决，当然也可以引入Seate进行解决此类问题，只不过跨语言调试需要付出一定的成本，Java内部使用没有问题，只不过目前公司纯Java服务之间的调用，且必须要同步事物这种业务场景目前还没遇到，或者说非常少

## 9使用RabbitMq
go-zero默认使用消息队列kafka，如果要使用rabbitMq需要自己集成 `github.com/streadway/amqp`,而且amqp这个sdk在使用层面几乎没有做任何封装，使用起来相当繁琐举几个例子
* amqp在网络端开的时候没有重连机制，必须重启服务器才能解决，如果要自动重连，需要程序员自己开发处理
* 网络断线重连后监听器的管道已经关闭，这个时候需要重新启动监听器
* amqp的事物消息非常鸡肋，他只能保证发送多条消息要么同时成功，要么同时失败，适用于数据同步的场景，并不能保证数据库与mq的事物一致性
* amqp消息监听写到后面发现每个消息监听的代码大同小异，而且比较繁琐，需要程序员做抽象简化
* amqp消息监听需要程序员手动配置死信队列，而且还需要处理死信消息，基础框架会自动为监听队列配置死信队列，并把死信消息序列化到数据库
* amqp消息重试机制非常鸡肋，要么不设置重试，要么一直重试直到成功，更别谈重试次数设定，这些只能靠程序员自己控制，基础框架结合死信消息，对消息进行重试，重试次数可以在配置文件配置

针对上面这些问题基础框架对amqp的交换机模式（fanout、direct、topic）做了封装，如果不使用交换机模式，则使用rabbitmq.Client自行处理,rabbitmq.Client里面本身继承了amqp.Connection

### 9.1 RabbitMq配置
```text
RabbitMq:
  Url: amqp://giikin:giikinn++@172.23.211.160:5672/test
  SendReliableMessage: true
  MaxRetryAttempts: 2
  DeadExchangeName: DeadExchangeName
  DeadQueueName: DeadQueueName
```
```go
type Conf struct {
    Url                 string
    SendReliableMessage bool   `json:",default=true"` // 使用可靠消息
    MaxRetryAttempts    int    `json:",default=3"`    // 错误消息重试次数(消息监听才需配置)
    DeadExchangeName    string `json:",optional"`     // 死信交换机(消息监听才需配置)
    DeadQueueName       string `json:",optional"`     // 死信队列(消息监听才需配置)
}
```

### 9.2 RabbitMq发送消息
### 9.2.1 声明交换机
发送消息之前需要声明交换机
```go
	rabbitmqClient := rabbitmq.NewClient(c.RabbitMq, conn, redisClient)
	rabbitmqClient.DeclareExchange(rabbitmq.ExchangeDeclareParam{
		Name:       "交换器名",   // 交换器名
		Kind:       "direct", // exchange type：一般用fanout、direct、topic
		Durable:    true,     // 是否持久化
		AutoDelete: false,    // 是否自动删除（自动删除的前提是至少有一个队列或者交换器与这和交换器绑定，之后所有与这个交换器绑定的队列或者交换器都与此解绑）
		Internal:   false,    // 设置是否内置的。true表示是内置的交换器，客户端程序无法直接发送消息到这个交换器中，只能通过交换器路由到交换器这种方式
		NoWait:     false,    // 是否阻塞若为true，则客户端不会等待RabbitMQ server返回当前交换器声明的结果。
		Args:       nil,      // 额外属性
	})
```
### 9.2.2 RabbitMq发送普通消息
```go
import (
    "github.com/wangliujing/foundation-framework/orm"
    "github.com/wangliujing/foundation-framework/rabbitmq"
    "github.com/wangliujing/foundation-framework/redis"
)

gormConnection := orm.NewGormConnection(c.Gorm)
redisClient := redis.NewClient(c.Redis)
rabbitmqClient := rabbitmq.NewClient(c.RabbitMq, gormConnection, redisClient)
// Publish(exchange, key string, message any) error  message可以是字符串也可以是对象
err := rabbitmqClient.Publish("test", "test", "abc")
// PublishDelay(exchange, key string, message any, delay int64) error  delay单位毫秒
err := rabbitmqClient.PublishDelay("test", "test", "abc", 10000) // 发送延迟消息
// 自定义消息
msg := amqp.Publishing{Body: []byte("abc"), MessageId: "123", DeliveryMode: amqp.Persistent, ContentType: "application/json"}
rabbitmqClient.DoPublish("test", "test", msg)
```
这种发送方式只能用在可容忍消息丢失的业务场景，因为这种方式只保证数据发送到Exchange交换机，但是并不能保证消息路由到队列成功，比如此时路由队列并没有绑定交换机，发送消息还是返回成功，但是实际消息已经丢失了

,根据前面分布式方案介绍这种方案有数据不一致的风险，因为消息发送成功，业务不一定执行成功，比如提交事物的时候数据库宕机，这时候业务数据没有保持，但是消息已经发送到rabbitMq了

### 9.2.3 RabbitMq发送可靠消息
```go
import (
    "github.com/wangliujing/foundation-framework/orm"
    "github.com/wangliujing/foundation-framework/rabbitmq"
    "github.com/wangliujing/foundation-framework/redis"
)


redisClient := redis.NewClient(c.Redis)
rabbitmqClient := rabbitmq.NewClient(c.RabbitMq, gormConnection, redisClient)
// NotifyPublish(exchange, key string, message any, messageId string) error
rabbitmqClient.NotifyPublish("test", "test", "abc", "123", func(messageId string) error {
    // 消息发送成功后回调
    return nil
})
// NotifyPublishDelay(exchange, key string, message any, messageId string, delay int64) error
rabbitmqClient.NotifyPublishDelay("test", "test", "abc", "123", 10000, nil) // 发送延迟消息

// 自定义消息格式
msg := amqp.Publishing{Body: []byte("abc"), MessageId: "123", DeliveryMode: amqp.Persistent,
    ContentType: "application/json"}
rabbitmqClient.DoNotifyPublish("test", "test", msg, func(messageId string) error {
    // 消息发送成功后回调
    return nil
})
```
这种方式发送消息基础框架封装了通知机制与队列路由return机制，可以保证消息一定发送到rabbitMq队列，但是该方式保证不了数据的一致性

注意：可靠消息发送SendReliableMessage必须设置为true

### 9.2.4 RabbitMq发送事物消息（7.1 异步方案发送端实现）
事物消息的实现是基于`8.2.2 可靠消息方案`的，在SendReliableMessage设置为true的情况下，基础框架会自动创建本地事物表`transactional_msg`

| 字段名                                                               | 类型           | 备注                         |
|-------------|--------------|----------------------------|
| id          | char(36)     | 消息ID                       |
| exchange    | varchar(100) | 交换机                        |
| routing_key | varchar(100) | 路由键                        |
| message | longtext     | 消息体                        |
| trigger_time | datetime     | 触发时间（定时任务会根据这个字段过滤出待发送的消息） |
| created_at | datetime       | 创建时间                       |
| deleted_at | datetime       | 删除时间                       |

发送消息时本地消息表存储消息与业务在同一个事物执行，事物提交后，发送消息，消息发送成功后逻辑删除本地事物表`transactional_msg`数据，消息发送失败，后续定时任务会扫描未删除的消息一直重试直到消息发送成功，这些逻辑基础框架全部封装，程序员不用关心细节，按如下示例代码发送事物消息即可
```go
import (
    "github.com/wangliujing/foundation-framework/orm"
    "github.com/wangliujing/foundation-framework/rabbitmq"
    "github.com/wangliujing/foundation-framework/redis"
)
	gormConnection := orm.NewGormConnection(c.Gorm)
	redisClient := redis.NewClient(c.Redis)
	rabbitmqClient := rabbitmq.NewClient(c.RabbitMq, gormConnection, redisClient)
	rabbitmqClient.SendMessageInTransaction("test", "test", func(connection orm.Connection) error {
		// 这里写业务逻辑，操作数据库是为了保证事物与本地消息表是在同一个事物执行，所有操作数据库必须使用参数 connection orm.Connection
		err := repo.NewAwsMsgRepository(connection).DeleteById(1) // 示例写法
		return nil
	}, "消息内容")
```

### 9.3 RabbitMq监听消息
消息监听器需实现如下接口
```go
type Listener interface {
	CoroutineNum() int // 消费消息协程数量
	ConsumerParam() ConsumerParam // 消费参数
	AckParam() AckParam // 确认参数
	QueueDeclareParam() QueueDeclareParam // 队列声明参数
	OnDelivery(delivery *amqp.Delivery) error // 消息监听方法
}
```

消息监听示例代码
```go
package listener

import (
	"github.com/streadway/amqp"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/wangliujing/foundation-framework/err/comm"
	"github.com/wangliujing/foundation-framework/rabbitmq"
)

type MsgListener struct {
}

func (m MsgListener) CoroutineNum() int {
	return 1
}

func (m MsgListener) ConsumerParam() rabbitmq.ConsumerParam {
	return rabbitmq.ConsumerParam{
		Consumer:  "dead_queue_consumer",
		AutoAck:   false,
		Exclusive: false,
		NoLocal:   false,
		NoWait:    false,
		Args:      nil,
	}
}

func (m MsgListener) AckParam() rabbitmq.AckParam {
	return rabbitmq.AckParam{Multipart: false, Requeue: false}
}

func (m MsgListener) QueueDeclareParam() rabbitmq.QueueDeclareParam {
	return rabbitmq.QueueDeclareParam{
		QueueName:    "wangliujing",
		Durable:      true,
		AutoDelete:   false,
		Exclusive:    false,
		NoWait:       false,
		RoutingKey:   "test",
		ExchangeName: "test",
	}
}

func (m MsgListener) OnDelivery(delivery *amqp.Delivery) error {
	logx.Info("收到消息：", string(delivery.Body))
	return nil
}
```
注册监听器
```go
import (
    "github.com/wangliujing/foundation-framework/orm"
    "github.com/wangliujing/foundation-framework/rabbitmq"
    "github.com/wangliujing/foundation-framework/redis"
)

gormConnection := orm.NewGormConnection(c.Gorm)
redisClient := redis.NewClient(c.Redis)
rabbitmqClient := rabbitmq.NewClient(c.RabbitMq, gormConnection, redisClient)
// 注册监听器后，监听器就可以接受到消息了
rabbitmqClient.RegisterListener(&listener.MsgListener{})
```
注册监听器的时候，基础框架会自动创建死信交换机、死信队列、并把业务监听队列绑定到死信交换机，会创建一张死信表`dead_msg`，用于存放死信消息

| 字段名                                                               | 类型           | 备注   |
|-------------|--------------|------|
| id          | char(36)     | 消息ID |
| exchange    | varchar(100) | 交换机  |
| routing_key | varchar(100) | 路由键  |
| message | longtext     | 消息体  |
| retry_count | bigint     | 重试次数 |
| created_at | datetime       | 创建时间 |
| deleted_at | datetime       | 删除时间 |

## 10GRPC与JsonRpc
go语言使用的rpc是GRPC协议，我们大部分工程使用PHP开发，PHP使用的rpc是JsonRpc协议，之前同事用php使用GRPC协议各种原因没有调通，所以只能在go这边做兼容，所以go首先得发布JsonRpc服务方便PHP调用，go也得封装JsonRpc的客户端用来调用PHP，以此来实现go与php的互通，为了方便使用，基础框架做了统一封装，具体看下面各个注册中心的使用方式

### 10.1 使用consul注册中心
### 10.1.1 注册GRPC服务
配置文件
```text
Name: deposit
ListenOn: 0.0.0.0:8080
Consul:
  Host: 172.30.20.23:8500
  Key: deposit
```
配置类
```go
import (
    "github.com/wangliujing/foundation-framework/reg/consul"
)

type Config struct {
    zrpc.RpcServerConf
    Consul     consul.Conf
}
```

注册服务

```go
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/wangliujing/foundation-framework/common/security"
	"github.com/wangliujing/foundation-framework/jsonrpc"
	"github.com/wangliujing/foundation-framework/reg/consul"
	"github.com/wangliujing/foundation-framework/system"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"service/internal/config"
	"service/internal/grpc/deposit"
	depositserviceServer "service/internal/server/depositservice"
	"service/internal/svc"
)

var configFile = flag.String("f", "etc/deposit.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.MustSetup(c.Log)
	ctx := svc.NewServiceContext(c)
	// 创建注册器
	registry := consul.NewRegistry(c.Consul, c.RpcServerConf)
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		// 注册服务到consul
		registry.RegisterService()
		// 注册服务到内存，方便rpc后续拿到对象调用方法
		deposit.RegisterDepositServiceServer(grpcServer, depositserviceServer.NewDepositServiceServer(ctx))
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	system.Start(s)

}
```

### 10.1.2 调用GRPC服务
配置文件
```text
DepositService:
  Timeout: -1
  Target: consul://172.30.20.23:8500/deposit?wait=14s
```
配置类
```go
import 	"github.com/zeromicro/go-zero/zrpc"

type Config struct {
	rest.RestConf
	DepositService zrpc.RpcClientConf
}
```

服务调用
```go
// 需要包，主要作用是注册consul解析器(非常重要，不然服务启动报错)
import _ "github.com/zeromicro/zero-contrib/zrpc/registry/consul"

// NewDepositService(cli zrpc.Client) DepositService 这个方法是grpc自动生成的代码
depositServiceRpc := depositService.NewDepositService(zrpc.MustNewClient(c.DepositService))
// 直接调用接口的方法
depositServiceRpc.Deposit(l.ctx, d)
```

### 10.1.3 注册JsonRpc服务
配置文件
```text
Name: deposit
Consul:
  Host: 172.30.20.23:8500
  Key: deposit
  JsonRpcConf:
    ListenOn: 0.0.0.0:8081
```

配置类
```go
package config

import (
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/wangliujing/foundation-framework/reg/consul"
)

type Config struct {
	zrpc.RpcServerConf
	Consul     consul.Conf
}
```

注册代码
```go
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/wangliujing/foundation-framework/common/security"
	"github.com/wangliujing/foundation-framework/jsonrpc"
	"github.com/wangliujing/foundation-framework/reg/consul"
	"github.com/wangliujing/foundation-framework/system"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"service/internal/config"
	"service/internal/grpc/deposit"
	depositserviceServer "service/internal/server/depositservice"
	"service/internal/svc"
)

var configFile = flag.String("f", "etc/deposit.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.MustSetup(c.Log)
	ctx := svc.NewServiceContext(c)
    // 启动JsonRpc服务器，方便其他服务进行JsonRpc远程调用
	jsonRpcService := jsonrpc.Start(c.Consul.JsonRpcConf.ListenOn)
	registry := consul.NewRegistry(c.Consul, c.RpcServerConf)
	// 创建JsonRpc注册器
	rpcRegistry := registry.NewJsonRpcRegister(jsonRpcService)
	// 注册JsonRpc服务
	rpcRegistry.RegisterService("wangliujing", new(SubscriptionRpcService), "abc=1&bbc=3")
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	system.Start(s)
}

type SubscriptionRpcService struct {
}

func (t *SubscriptionRpcService) Subscription(ctx context.Context, subscription any, reply *any) error {
	newContext := security.NewContext(ctx)
	requestContext := newContext.GetXRequestContext()
	fmt.Print(requestContext)
	*reply = subscription
	return nil
}
```

### 10.1.4 调用JsonRpc服务
配置文件
```text
Consul:
  Host: 172.30.20.23:8500
  Key: deposit
  JsonRpcConf:
    DiscoveryServices:
      - Service: JShopService
      - Service: UuidService
```

配置类
```go
package config

import (
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/wangliujing/foundation-framework/reg/consul"
)

type Config struct {
	zrpc.RpcServerConf
	Consul     consul.Conf
}
```

服务调用
```go
import "github.com/wangliujing/foundation-framework/reg/consul"


registry := consul.NewRegistry(c.Consul, c.RpcServerConf)
jsonRpcClient := registry.NewJsonRpcClient()
args := make(map[string]int64)
args["shopId"] = shopId
args["action"] = 1
var shopSecret ShopSecret
// Call(serviceName string, serviceMethod string, args any, reply any) error
err := jsonRpcClient.Call("JShopService", "/j_shop/getShopApplyKey", args, &shopSecret)
```



### 10.2 使用nacos注册中心
go-zero原生包"github.com/zeromicro/zero-contrib/zrpc/registry/nacos"在nacos有本地缓存的时候启动报错，该问题已在基础框架进行修复，详细问题如果有兴趣可以查看https://mp.csdn.net/mp_blog/creation/editor/133209819

### 10.2.1 注册GRPC服务
配置文件
```text
Name: deposit
ListenOn: 0.0.0.0:8080
Nacos:
  Hosts:
    -
      Ip: 172.18.0.145
      Port: 8848
  Namespace: local_test
```
配置类
```go
package config

import (
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/wangliujing/foundation-framework/reg/nacos"
)

type Config struct {
	zrpc.RpcServerConf
	Nacos      nacos.Conf
}
```
注册代码
```go
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/wangliujing/foundation-framework/common/security"
	"github.com/wangliujing/foundation-framework/jsonrpc"
	"github.com/wangliujing/foundation-framework/reg/nacos"
	"github.com/wangliujing/foundation-framework/system"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"service/internal/config"
	"service/internal/grpc/deposit"
	depositserviceServer "service/internal/server/depositservice"
	"service/internal/svc"
)

var configFile = flag.String("f", "etc/deposit.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	// 创建nacos注册器
	registry := nacos.NewRegistry(c.Nacos, c.RpcServerConf)
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		// 注册服务到nacos
		registry.RegisterService()
		// 注册服务对象到内存，方便后续获取对象进行方法调用
		deposit.RegisterDepositServiceServer(grpcServer, depositserviceServer.NewDepositServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	system.Start(s)
}
```

### 10.2.2 调用GRPC服务
配置文件
```text
DepositService:
  Timeout: -1
  Target: nacos://172.18.0.145:8848/deposit?namespaceid=local_test
```

配置类
```go
package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/wangliujing/foundation-framework/reg/nacos"
	"github.com/wangliujing/foundation-framework/shenyu"
)

type Config struct {
	rest.RestConf
	DepositService zrpc.RpcClientConf
}

```

服务调用
```go
// NewDepositService(cli zrpc.Client) DepositService 这个方法是grpc自动生成的代码
depositServiceRpc := depositService.NewDepositService(zrpc.MustNewClient(c.DepositService))
// 直接调用接口的方法
depositServiceRpc.Deposit(l.ctx, d)
```

### 10.2.3 注册JsonRpc服务
配置文件
```text
Name: deposit
Nacos:
  Hosts:
    -
      Ip: 172.18.0.145
      Port: 8848
  Namespace: local_test
  JsonRpcConf:
    ListenOn: 0.0.0.0:8001
```

配置类
```go
package config

import (
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/wangliujing/foundation-framework/reg/nacos"
)

type Config struct {
	zrpc.RpcServerConf
	Nacos      nacos.Conf
}
```
注册代码
```go
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/wangliujing/foundation-framework/common/security"
	"github.com/wangliujing/foundation-framework/jsonrpc"
	"github.com/wangliujing/foundation-framework/reg/nacos"
	"github.com/wangliujing/foundation-framework/system"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"service/internal/config"
	"service/internal/grpc/deposit"
	depositserviceServer "service/internal/server/depositservice"
	"service/internal/svc"
)

var configFile = flag.String("f", "etc/deposit.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	// 创建注册器
	registry := nacos.NewRegistry(c.Nacos, c.RpcServerConf)
	// 创建JsonRpc注册器
	r := registry.NewJsonRpcRegistry(jsonrpc.Start(c.Nacos.JsonRpcConf.ListenOn))
	// 注册JsonRpc服务到nacos
	r.RegisterService("wangliujing", &SubscriptionRpcService{}, "abc=1&bbc=8")
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	system.Start(s)
}

type SubscriptionRpcService struct {
}

func (t *SubscriptionRpcService) Subscription(ctx context.Context, subscription any, reply *any) error {
	newContext := security.NewContext(ctx)
	requestContext := newContext.GetXRequestContext()
	fmt.Print(requestContext)
	*reply = subscription
	return nil
}
```

### 10.2.4 调用JsonRpc服务
配置文件：
```text
Nacos:
  Hosts:
    -
      Ip: 172.18.0.145
      Port: 8848
  Namespace: local_test
  JsonRpcConf:
    DiscoveryServices:
      - Service: JShopService
```

配置类
```go
package config

import (
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/wangliujing/foundation-framework/reg/nacos"
)

type Config struct {
	zrpc.RpcServerConf
	Nacos      nacos.Conf
}
```

服务调用
```go
improt 	"github.com/wangliujing/foundation-framework/reg/nacos"

registry := nacos.NewRegistry(c.Nacos, c.RpcServerConf)
jsonRpcClient := registry.NewJsonRpcClient()
args := make(map[string]int64)
args["shopId"] = shopId
args["action"] = 1
var shopSecret ShopSecret
// Call(serviceName string, serviceMethod string, args any, reply any) error
err := jsonRpcClient.Call("JShopService", "/j_shop/getShopApplyKey", args, &shopSecret)
```

## 11接入XXL-JOB
### 11.1 配置类
```text
type Conf struct {
	Token        string `json:",default=default_token"`
	AppName      string `json:",optional"` // 
	ClientIp     string `json:",optional"`
	ClientPort   string `json:",optional"`
	AdminAddress string
	Log          LogConf
}

type LogConf struct {
	KeepDays int `json:",default=30"`
}
```
### 11.2 配置示例
```text
Xxl:
  AppName: appName
  ClientIp: 192.168.1.29
  ClientPort: "9999"
  AdminAddress: "http://192.168.1.79:8080/xxl-job-admin"
  Log:
    KeepDays: 7 // 本地任务日志保留天数
```
### 11.3 使用示例
### 11.3.1 示例代码
```go
   gormConnection := orm.NewGormConnection(c.Gorm)
   redisClient := redis.NewClient(c.Redis)
   executor := xxl.NewExecutor(c.Xxl)
   // 设置日志处理器方便写入日志，以及admin调用日志进行查看
   executor.Log(xxl.NewMysqlLog(gormConnection, redisClient, c.Xxl.Log))
   //注册任务handler
   executor.RegTask("TestJob1", TestJob1)
   executor.RegTask("TestJob2", TestJob2)
```
```go
func TestJob1(ctx context.Context, param *xxl.RunReq, log xxl.Log) error {
	fmt.Println("test one task" + param.ExecutorHandler + " param：" + param.ExecutorParams + " log_id:" + xxl.Int64ToStr(param.LogID))
	// 注意这里必须使用log参数写入日志，该log对象就是之前xxl.NewMysqlLog(gormConnection, redisClient, c.Xxl.Log)
	log.Info(param, "测试日志1")
	log.Info(param, "测试日志2")
	log.Error(param, "错误信息：%+v", biz.New("自定义错误"))
	time.Sleep(10 * time.Second)
	return nil
}

// 任务执行时间比较长，而且需要手动终止任务可以参照该写法
func TestJob2(ctx context.Context, param *xxl.RunReq, log xxl.Log) error {
	num := 1
	for {
		select {
		case <-ctx.Done():
			log.Info(param, "task"+param.ExecutorHandler+"被手动终止")
			return nil
		default:
			num++
			time.Sleep(10 * time.Second)
			if num > 10 {
				log.Info(param, "test one task"+param.ExecutorHandler+" param："+param.ExecutorParams+"执行完毕！")
				return nil
			}
		}
	}

}
```

### 11.3.2 在admin手动注册执行器与任务
上述代码编写完成后，还需如下3步才能自动调度任务
*  需要在admin控制台注册执行器，执行器的AppName值一定要与配置文件AppName值一样
*  在任务管理界面注册任务，执行器还是选择上一步注册的执行器，JobHandler的值一定要与代码里面保持一致,比如 `executor.RegTask("TestJob1", TestJob1)`中 `TestJob1`
*  在控制台启动任务

### 11.3.3 MysqlLog详解
调用xxl.NewMysqlLog(gormConnection, redisClient, c.Xxl.Log)的时候会自动使用gormConnection对象在数据库创建日志表`xxl_log`，该表结构用来存储任务执行的过程中产生的日志，前提是使用MysqlLog对象进行日志写入，可以参照 `11.3.1 示例代码`，admin查看执行日志会主动通过接口获取任务相关执行日志

xxl_log表结构信息如下：

| 字段名          | 类型           | 备注                         |
|--------------|--------------|----------------------------|
| id           | bigint     | 主键ID                       |
| log_id       | bigint| admin日志ID                  |
| level        | varchar(100) | 日志级别 0 info 1 error        |
| text         | text     | 日志内容                       |
| created_at   | datetime       | 创建时间                       |

当然用户如果想使用mongodb作为日志存储则不使用框架提供的MysqlLog，只需要自定义实现如下日志接口来实现自己的相关逻辑即可
```go
type Log interface {
	Info(param *RunReq, format string, a ...any)
	Error(param *RunReq, format string, a ...any)
	Handler(req *LogReq) *LogRes
}
```
## 12接入nacos配置中心
### 12.1 配置类
```go
type Conf struct {
	Hosts               []Host
	Namespace           string `json:",default=public"`
	NotLoadCacheAtStart bool   `json:",default=true"`
	LogLevel            string `json:",default=info"` // the level of log, it's must be debug,info,warn,error, default value is info
	LogDir              string `json:",default=/tmp/nacos/log"`
	CacheDir            string `json:",default=/tmp/nacos/cache"`
	ConfigParam         ConfigParam
}

type ConfigParam struct {
   DataId string
   Group  string `json:",default=DEFAULT_GROUP"`
   Type   vo.ConfigType `json:",default=yaml"`
}

type Host struct {
   Ip   string
   Port uint64
}

```
### 12.2 配置示例
```text
Hosts:
  -
    Ip: 172.18.0.145
    Port: 8848
Namespace: test    
ConfigParam:
  DataId: deposit
  Type: yaml
```

### 12.3 nacos配置内容
```text
Name: deposit
ListenOn: 0.0.0.0:8080
Consul:
  Host: 172.30.20.23:8500
  Key: deposit
  JsonRpcConf:
    ListenOn: 0.0.0.0:8081
    DiscoveryServices:
      - Service: AuthorizationService
      - Service: UuidService

Log:
  Encoding: plain
  #Mode: file
GormConfig:
  Dsn: jmy:Qwer1234@tcp(rm-bp1k0fj4ej13saps7.mysql.rds.aliyuncs.com:3306)/aws_message?charset=utf8mb4&parseTime=True&loc=Local
  MaxIdleConns: 10
  MaxOpenConns: 50
  ConnMaxIdleTime: 5m
  ConnMaxLifetime: 1h
  OpenLogicDeleted: true
  LogLevel: 4
```

### 12.4 示例代码
```go
// 业务配置类,该配置类对应nacos配置内容
type Config struct {
	zrpc.RpcServerConf
	GormConfig orm.GormConfig
	Consul     consul.Conf
}

var c Config
// 把nacos配置中心内容拉下来给对象Config赋值，并且对nacos内容设置监听
nacos.MustLoad(*configFile, &c)
```


## 13服务上下文ServiceContext
ServiceContext 是go-zero的一种编码习惯，这个类型是goctl代码生成器自动生成的，主要用于在方法间传递共用对象，比如数据库连接，redis客户端，rabbitMq客户端，配置类，远程调用客户端等，通过ServiceContext对这些对象进行管理，避免了客户端的重复创建，规范代码，使用的时候只需要传递ServiceContext，增强代码的可读性与可维护性

示例代码
```go
package svc

import (
   "github.com/wangliujing/foundation-framework/orm"
   "service/internal/config"
)

type ServiceContext struct {
   Config      config.Config
   Conn        orm.Connection
   Transaction orm.Transaction
}

func NewServiceContext(c config.Config) *ServiceContext {
   conn := orm.NewGormConnection(c.GormConfig)
   return &ServiceContext{
      Config:      c,
      Conn:        conn,
      Transaction: orm.NewTransaction(conn),
   }
}

```