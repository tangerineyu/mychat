package handler

import (
	"my-chat/internal/service"
	"my-chat/pkg/errno"

	"github.com/gin-gonic/gin"
)

type GroupHandler struct {
	groupService *service.GroupService
}

func NewGroupHandler(groupService *service.GroupService) *GroupHandler {
	return &GroupHandler{groupService: groupService}
}

type CreateReq struct {
	Name    string `json:"name" binding:"required"`
	OwnerId string `json:"owner_id" binding:"required"`
}

func (h *GroupHandler) Create(c *gin.Context) {
	var req CreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}
	group, err := h.groupService.CreateGroup(req.OwnerId, req.Name)
	if err != nil {
		SendResponse(c, err, nil)
		return
	}
	SendResponse(c, nil, group)
}

type JoinReq struct {
	GroupId string `json:"group_id" binding:"required"`
	UserId  string `json:"user_id" binding:"required"`
}

func (h *GroupHandler) Join(c *gin.Context) {
	var req JoinReq
	if err := c.ShouldBindJSON(&req); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}
	err := h.groupService.JoinGroup(req.GroupId, req.UserId)
	if err != nil {
		SendResponse(c, err, nil)
		return
	}
	SendResponse(c, nil, gin.H{"message": "加入群聊成功"})
}

type GroupInfoReq struct {
	GroupId string `json:"group_id" binding:"required"`
}

func (h *GroupHandler) GetGroupInfo(c *gin.Context) {
	var req GroupInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}
	group, err := h.groupService.GetGroupInfo(req.GroupId)
	if err != nil {
		SendResponse(c, err, nil)
		return
	}
	SendResponse(c, nil, group)
}
func (h *GroupHandler) GetGroupMemberList(c *gin.Context) {
	var req GroupInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}
	members, err := h.groupService.GetGroupMembers(req.GroupId)
	if err != nil {
		SendResponse(c, err, nil)
		return
	}
	SendResponse(c, nil, members)
}
