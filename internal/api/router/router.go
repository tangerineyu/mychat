package router

import (
	"my-chat/internal/api/handler"
	"my-chat/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

// Register registers all HTTP routes.
func Register(r *gin.Engine, userHandler *handler.UserHandler, wsHandler *handler.WSHandler,
	groupHandler *handler.GroupHandler, chatHandler *handler.ChatHandler, contactHandler *handler.ContactHandler,
	sessionHandler *handler.SessionHandler, adminHandler *handler.AdminHandler,
) {
	v1 := r.Group("/api/v1")
	{
		v1.POST("/register", userHandler.Register)
		v1.POST("/login", userHandler.Login)
		v1.POST("/refresh-token", userHandler.RefreshToken)
	}

	authGroup := v1.Group("/")
	authGroup.Use(middleware.Auth())
	{
		authGroup.GET("/ws", wsHandler.Connect)
		// 用户相关
		authGroup.POST("/upload/avatar", userHandler.UploadAvatar)
		authGroup.POST("/user/updateUserInfo", userHandler.UpdateUserInfo)
		// 群组相关
		authGroup.POST("/group/create", groupHandler.Create)
		authGroup.POST("/group/join", groupHandler.Join)
		authGroup.POST("/group/getGroupInfo", groupHandler.GetGroupInfo)
		authGroup.POST("/group/getGroupMemberList", groupHandler.GetGroupMemberList)
		authGroup.POST("/group/loadMyGroup", groupHandler.LoadMyJoinedGroup)
		authGroup.POST("/group/leaveGroup", groupHandler.LeaveGroup)
		authGroup.POST("/group/kickGroupMember", groupHandler.KickGroupMember)
		authGroup.POST("/group/dismissGroup", groupHandler.DismissGroup)
		// 联系人相关
		authGroup.POST("/contact/add", contactHandler.AddFriend)
		authGroup.POST("/contact/agree", contactHandler.AgreeFriend)
		authGroup.POST("/contact/list", contactHandler.GetContactList)
		authGroup.POST("/contact/applyList", contactHandler.GetApplyList)
		authGroup.POST("/contact/refuseContactApply", contactHandler.RefuseApply)
		authGroup.POST("/contact/deleteContact", contactHandler.DeleteContact)
		authGroup.POST("/contact/blackContact", contactHandler.BlackContact)
		authGroup.POST("/contact/cancelBlackContact", contactHandler.UnBlackContact)
		// 聊天历史记录
		authGroup.POST("/chat/history", chatHandler.History)
		// 会话接口
		authGroup.POST("/session/list", sessionHandler.List)
		// Admin User
		authGroup.POST("/user/getUserInfoList", adminHandler.GetUserList)
		authGroup.POST("/user/disableUsers", adminHandler.DisableUser)
		authGroup.POST("/user/ableUsers", adminHandler.AbleUser)
		// Admin Group
		authGroup.POST("/group/getGroupInfoList", adminHandler.GetGroupList)
		authGroup.POST("/group/disableGroup", adminHandler.DisableGroup)
		authGroup.POST("/group/ableGroups", adminHandler.AbleGroup)
	}
}
