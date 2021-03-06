package proxy

import (
	"github.com/galaxy-book/saturn/model/context"
	"github.com/galaxy-book/saturn/model/req"
	"github.com/galaxy-book/saturn/model/resp"
)

type Proxy interface {
	// GetTenantAccessToken 获取企业认证Token
	GetTenantAccessToken(tenantKey string) resp.GetTenantAccessTokenResp
	// CodeLogin 免登, tenantKey非必填
	CodeLogin(tenantKey, code string) resp.CodeLoginResp
	// GetUsers 获取用户列表，部门ID未指定时查询所有用户
	GetUsers(ctx *context.Context, req req.GetUsersReq) resp.GetUsersResp
	// GetUser 获取用户信息
	GetUser(ctx *context.Context, id string) resp.GetUserResp
	// GetDeptIds 获取部门ID列表，不包含顶级部门及父部门
	GetDeptIds(ctx *context.Context, req req.GetDeptIdsReq) resp.GetDeptIdsResp
	// GetDepts 获取部门列表
	GetDepts(ctx *context.Context, req req.GetDeptsReq) resp.GetDeptsResp
	// GetRootDept 获取跟部门
	GetRootDept(ctx *context.Context) resp.GetRootDeptResp
	// SendMsg 发送消息，数据结构各自实现
	SendMsg(ctx *context.Context, req req.SendMsgReq) resp.SendMsgResp
}

type Ticket func() (string, error)
