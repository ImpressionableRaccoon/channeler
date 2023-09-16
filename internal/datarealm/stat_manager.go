package datarealm

import "context"

type StatManager interface {
	GetChannelAdminLog(ctx context.Context, channelID, channelAccessHash, minID int64) (ChannelAdminLogResponser, error)
}

type ChannelAdminLogResponser interface {
	GetEvents() []Event
	GetUsers() []User
}
