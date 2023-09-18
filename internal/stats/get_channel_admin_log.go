package stats

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gotd/td/tg"

	"github.com/ImpressionableRaccoon/channeler/internal/datarealm"
)

type getChannelAdminLogResponse struct {
	Events []datarealm.Event
	Users  []datarealm.User
}

var _ datarealm.ChannelAdminLogResponser = getChannelAdminLogResponse{}

func (r getChannelAdminLogResponse) GetEvents() []datarealm.Event {
	return r.Events
}

func (r getChannelAdminLogResponse) GetUsers() []datarealm.User {
	return r.Users
}

func (c client) GetChannelAdminLog(ctx context.Context, channelID, channelAccessHash, minID int64) (
	datarealm.ChannelAdminLogResponser, error,
) {
	resp, err := c.t.API().ChannelsGetAdminLog(ctx, &tg.ChannelsGetAdminLogRequest{
		Channel: &tg.InputChannel{
			ChannelID:  channelID,
			AccessHash: channelAccessHash,
		},
		MinID: minID,
	})
	if err != nil {
		return getChannelAdminLogResponse{}, err
	}

	res := getChannelAdminLogResponse{
		Events: make([]datarealm.Event, 0, len(resp.GetEvents())),
		Users:  make([]datarealm.User, 0, len(resp.GetUsers())),
	}

	for _, event := range resp.GetEvents() {
		e := datarealm.Event{
			ID:         event.GetID(),
			Date:       time.Unix(int64(event.GetDate()), 0),
			UserID:     event.GetUserID(),
			ActionType: event.GetAction().TypeName(),
		}

		if data, err := json.Marshal(event.GetAction()); err == nil {
			e.Action = string(data)
		}

		res.Events = append(res.Events, e)
	}

	for _, user := range resp.GetUsers() {
		u, ok := user.(*tg.User)
		if !ok {
			continue
		}
		res.Users = append(res.Users, tgUserToDataRealmUser(u))
	}

	return res, nil
}
