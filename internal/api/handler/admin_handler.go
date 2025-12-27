package handler

import (
	"my-chat/internal/service"

	"github.com/gin-gonic/gin"
	//"my-chat/pkg/errno"
)

type AdminHandler struct {
	adminService *service.AdminService
}

func NewAdminHandler(adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

type GetListReq struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func (h *AdminHandler) GetUserList(c *gin.Context) {
	var req GetListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		req.Page = 1
		req.Limit = 10
	}
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.Page == 0 {
		req.Page = 1
	}
	data, err := h.adminService.GetUserList(req.Page, req.Limit)
	if err != nil {
		SendResponse(c, err, nil)
		return
	}
	SendResponse(c, err, data)
}

type OperateUserReq struct {
	UuidList []string `json:"uuid_list" binding:"required"`
}

func (h *AdminHandler) DisableUser(c *gin.Context) {
	var req OperateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		SendResponse(c, err, nil)
		return
	}
	for _, uuid := range req.UuidList {
		_ = h.adminService.BanUser(uuid)
	}
	SendResponse(c, nil, gin.H{"msg": "操作成功"})
}
func (h *AdminHandler) AbleUser(c *gin.Context) {
	var req OperateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		SendResponse(c, err, nil)
		return
	}
	for _, uuid := range req.UuidList {
		_ = h.adminService.UnBanUser(uuid)
	}
	SendResponse(c, nil, gin.H{"msg": "操作成功"})
}
func (h *AdminHandler) GetGroupList(c *gin.Context) {
	var req GetListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		req.Page = 1
		req.Limit = 10
	}
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.Page == 0 {
		req.Page = 1
	}
	data, err := h.adminService.GetGroupList(req.Page, req.Limit)
	if err != nil {
		SendResponse(c, err, nil)
		return
	}
	SendResponse(c, err, data)
}
func (h *AdminHandler) DisableGroup(c *gin.Context) {
	var req OperateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		SendResponse(c, err, nil)
		return
	}
	for _, uuid := range req.UuidList {
		_ = h.adminService.BanGroup(uuid)
	}
	SendResponse(c, nil, gin.H{"msg": "操作成功"})
}
func (h *AdminHandler) AbleGroup(c *gin.Context) {
	var req OperateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		SendResponse(c, err, nil)
		return
	}
	for _, uuid := range req.UuidList {
		_ = h.adminService.UnBanGroup(uuid)
	}
	SendResponse(c, nil, gin.H{"msg": "操作成功"})
}
