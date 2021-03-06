package ding

import (
	"github.com/galaxy-book/saturn/model/context"
	"github.com/galaxy-book/saturn/model/req"
	"github.com/galaxy-book/saturn/model/resp"
	"github.com/galaxy-book/saturn/util/queue"
	"github.com/polaris-team/dingtalk-sdk-golang/sdk"
	"log"
	"strconv"
)

func (d *dingProxy) GetDeptIds(ctx *context.Context, req req.GetDeptIdsReq) resp.GetDeptIdsResp {
	client := &sdk.DingTalkClient{
		AccessToken: ctx.TenantAccessToken,
		AgentId:     d.AgentId,
	}
	deptIdContains := map[string]bool{}
	q := queue.New()
	q.Push("1")
	if req.ParentId != "" {
		q.Clear()
		q.Push(req.ParentId)
	}
	for {
		obj, err := q.Pop()
		if err != nil {
			break
		}
		parentId := obj.(string)
		subIdsResp, err := client.GetSubDept(parentId)
		if err != nil {
			log.Println(err)
			continue
		}
		if subIdsResp.ErrCode != 0 {
			log.Println(subIdsResp.ErrCode, subIdsResp.ErrMsg)
			continue
		}
		for _, subId := range subIdsResp.SubDeptIdList {
			deptStrId := strconv.FormatInt(subId, 10)
			if !deptIdContains[deptStrId] {
				deptIdContains[deptStrId] = true
				if req.FetchChild {
					q.Push(deptStrId)
				}
			}
		}
	}
	deptIds := make([]string, 0)
	for k, _ := range deptIdContains {
		deptIds = append(deptIds, k)
	}
	return resp.GetDeptIdsResp{
		Resp: resp.SucResp(),
		Data: deptIds,
	}
}

func (d *dingProxy) GetDepts(ctx *context.Context, req req.GetDeptsReq) resp.GetDeptsResp {
	client := &sdk.DingTalkClient{
		AccessToken: ctx.TenantAccessToken,
		AgentId:     d.AgentId,
	}
	depts := make([]resp.Dept, 0)
	deptIdContains := map[string]bool{}
	q := queue.New()
	q.Push("1")
	if req.ParentId != "" {
		q.Clear()
		q.Push(req.ParentId)
	}
	for {
		obj, err := q.Pop()
		if err != nil {
			break
		}
		parentId := obj.(string)
		subDepts, err := client.GetDeptList(nil, false, parentId)
		if err != nil {
			log.Println(err)
			continue
		}
		if subDepts.ErrCode != 0 {
			log.Println(subDepts.ErrCode, subDepts.ErrMsg)
			continue
		}
		for _, subDept := range subDepts.Department {
			deptStrId := strconv.FormatInt(subDept.Id, 10)
			if !deptIdContains[deptStrId] {
				deptIdContains[deptStrId] = true
				dept := resp.Dept{
					Name:     subDept.Name,
					ID:       strconv.FormatInt(subDept.Id, 10),
					ParentID: strconv.FormatInt(subDept.ParentId, 10),
				}
				dept.OpenID = dept.ID
				dept.ParentOpenID = dept.ParentID
				depts = append(depts, dept)
				if req.FetchChild {
					q.Push(deptStrId)
				}
			}
		}
	}
	deptIds := make([]string, 0)
	for k, _ := range deptIdContains {
		deptIds = append(deptIds, k)
	}
	return resp.GetDeptsResp{
		Resp: resp.SucResp(),
		Data: resp.GetDeptsRespData{
			Depts: depts,
		},
	}
}

func (d *dingProxy) GetRootDept(ctx *context.Context) resp.GetRootDeptResp {
	client := &sdk.DingTalkClient{
		AccessToken: ctx.TenantAccessToken,
		AgentId:     d.AgentId,
	}
	deptDetailResp, err := client.GetDeptDetail("1", nil)
	if err != nil {
		return resp.GetRootDeptResp{Resp: resp.ErrResp(err)}
	}
	if deptDetailResp.ErrCode != 0 {
		return resp.GetRootDeptResp{Resp: resp.Resp{Code: deptDetailResp.ErrCode, Msg: deptDetailResp.ErrMsg}}
	}
	dept := resp.Dept{
		Name:     deptDetailResp.Name,
		ID:       strconv.FormatInt(deptDetailResp.Id, 10),
		ParentID: strconv.FormatInt(deptDetailResp.ParentId, 10),
	}
	dept.OpenID = dept.ID
	dept.ParentOpenID = dept.ParentID
	return resp.GetRootDeptResp{
		Resp: resp.SucResp(),
		Data: dept,
	}
}
