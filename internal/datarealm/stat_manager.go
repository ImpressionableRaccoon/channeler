package datarealm

import "context"

type StatManager interface {
	GetChannelAdminLog(ctx context.Context, channelID, channelAccessHash, minID int64) (ChannelAdminLogResponser, error)
	GetUser(ctx context.Context, userID, userAccessHash int64) (User, error)
}

type ChannelAdminLogResponser interface {
	GetEvents() []Event
	GetUsers() []User
}
