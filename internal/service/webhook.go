package service

import (
	"github.com/WuKongIM/WuKongIM/internal/reactor"
	"github.com/WuKongIM/WuKongIM/internal/types"
	wkproto "github.com/WuKongIM/WuKongIMGoProto"
)

var Webhook webhook

type webhook interface {
	// Online 设备上线
	Online(uid string, deviceFlag wkproto.DeviceFlag, connId int64, deviceOnlineCount int, totalOnlineCount int)
	// Offline 设备下线
	Offline(uid string, deviceFlag wkproto.DeviceFlag, connId int64, deviceOnlineCount int, totalOnlineCount int)
	// NotifyOfflineMsg 离线消息通知
	NotifyOfflineMsg(msg reactor.ChannelMessage, subscribers []string)
	// TriggerEvent 触发事件
	TriggerEvent(event *types.Event)

	Start() error
	Stop()
}
