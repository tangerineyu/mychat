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
func (h *GroupHandler) LoadMyJoinedGroup(c *gin.Context) {
	userId := c.Query("user_id")
	groups, err := h.groupService.LoadMyJoinedGroup(userId)
	if err != nil {
		SendResponse(c, err, nil)
		return
	}
	SendResponse(c, nil, groups)
}

type GroupOperationReq struct {
	GroupId  string `json:"group_id" binding:"required"`
	TargetId string `json:"target_id"`
}

func (h *GroupHandler) LeaveGroup(c *gin.Context) {
	var req GroupOperationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}
	userId := c.Query("user_id")
	if err := h.groupService.LeaveGroup(req.GroupId, userId); err != nil {
		SendResponse(c, err, nil)
		return
	}
	SendResponse(c, nil, gin.H{"msg": "退群成功"})
}
func (h *GroupHandler) KickGroupMember(c *gin.Context) {
	var req GroupOperationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}
	userId := c.Query("user_id")
	if err := h.groupService.KickMember(userId, req.GroupId, req.TargetId); err != nil {
		SendResponse(c, err, nil)
		return
	}
	SendResponse(c, nil, gin.H{"msg": "移除成功"})
}
func (h *GroupHandler) DismissGroup(c *gin.Context) {
	var req GroupOperationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}
	userId := c.Query("user_id")
	if err := h.groupService.DismissGroup(userId, req.GroupId); err != nil {
		SendResponse(c, err, nil)
		return
	}
	SendResponse(c, nil, gin.H{"msg": "群聊已解散"})
}
