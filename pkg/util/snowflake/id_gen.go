package snowflake

import (
	"my-chat/pkg/zlog"

	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
)

var node *snowflake.Node

func Init(machineID int64) {
	var err error
	node, err = snowflake.NewNode(machineID)
	if err != nil {
		zlog.Error("init snowflake node error", zap.Error(err))
		panic(err)
	}
	zlog.Info("init snowflake node success",
		zap.Int64("machineID", machineID))
}

// 生成int64类型id
func GenID() int64 {
	return node.Generate().Int64()
}

// 生成string类型id
func GenStringID() string {
	return node.Generate().String()
}
