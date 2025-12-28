package handler

import (
	"my-chat/internal/service"
	"my-chat/pkg/errno"

	"github.com/gin-gonic/gin"
)

type ContactHandler struct {
	contactService *service.ContactService
}

func NewContactHandler(contactService *service.ContactService) *ContactHandler {
	return &ContactHandler{contactService: contactService}
}

type AddFriendReq struct {
	TargetId string `json:"target_id" binding:"required"`
	Msg      string `json:"msg"`
}

func (h *ContactHandler) AddFriend(c *gin.Context) {
	var req AddFriendReq
	if err := c.ShouldBindJSON(&req); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}
	//userId := c.Query("uid")
	userId := c.GetString("userId")
	err := h.contactService.AddFriendApply(userId, req.TargetId, req.Msg)
	if err != nil {
		SendResponse(c, err, nil)
		return
	}
	SendResponse(c, nil, gin.H{"message": "好友申请已发送"})
}

type AgreeReq struct {
	ApplyId uint `json:"apply_id" binding:"required"`
}

func (h *ContactHandler) AgreeFriend(c *gin.Context) {
	var req AgreeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}
	err := h.contactService.AgreeFriend(req.ApplyId)
	if err != nil {
		SendResponse(c, err, nil)
		return
	}
	SendResponse(c, nil, gin.H{"message": "已添加好友"})
}
func (h *ContactHandler) GetContactList(c *gin.Context) {
	userId := c.GetString("userId")
	list, err := h.contactService.GetContactList(userId)
	if err != nil {
		SendResponse(c, err, nil)
		return
	}
	SendResponse(c, nil, list)
}
func (h *ContactHandler) GetApplyList(c *gin.Context) {
	userId := c.GetString("userId")
	list, err := h.contactService.GetApplyList(userId)
	if err != nil {
		SendResponse(c, err, nil)
		return
	}
	SendResponse(c, nil, list)
}

type RefuseReq struct {
	ApplyId uint `json:"apply_id" binding:"required"`
}

func (h *ContactHandler) RefuseApply(c *gin.Context) {
	var req RefuseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}
	if err := h.contactService.RefuseApply(req.ApplyId); err != nil {
		SendResponse(c, err, nil)
		return
	}
	SendResponse(c, nil, gin.H{"msg": "已拒绝"})
}

type DeleteContact struct {
	TargetId string `json:"target_id" binding:"required"`
}

func (h *ContactHandler) DeleteContact(c *gin.Context) {
	var req DeleteContact
	if err := c.ShouldBindJSON(&req); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}
	userId := c.GetString("userId")
	if err := h.contactService.RemoveFriend(userId, req.TargetId); err != nil {
		SendResponse(c, err, nil)
		return
	}
	SendResponse(c, nil, gin.H{"message": "删除成功"})
}

type BlackReq struct {
	TargetId string `json:"target_id" binding:"required"`
}

func (h *ContactHandler) BlackContact(c *gin.Context) {
	var req BlackReq
	if err := c.ShouldBindJSON(&req); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}
	userId := c.GetString("userString")
	if err := h.contactService.BlackContact(userId, req.TargetId); err != nil {
		SendResponse(c, err, nil)
		return
	}
	SendResponse(c, nil, gin.H{"message": "已经拉黑"})
}
func (h *ContactHandler) UnBlackContact(c *gin.Context) {
	var req BlackReq
	if err := c.ShouldBindJSON(&req); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}
	userId := c.GetString("userId")
	if err := h.contactService.UnBlackContact(userId, req.TargetId); err != nil {
		SendResponse(c, err, nil)
		return
	}
	SendResponse(c, nil, gin.H{"msg": "已经移除黑名单"})
}
