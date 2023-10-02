package stats

import (
	"context"
	"errors"

	"github.com/gotd/td/tg"

	"github.com/ImpressionableRaccoon/channeler/internal/datarealm"
)

func (c client) GetUser(ctx context.Context, userID, userAccessHash int64) (datarealm.User, error) {
	fullUser, err := c.t.API().UsersGetFullUser(ctx,
		&tg.InputUser{
			UserID:     userID,
			AccessHash: userAccessHash,
		},
	)
	if err != nil {
		return datarealm.User{}, err
	}

	users := fullUser.GetUsers()
	if len(users) == 0 {
		return datarealm.User{}, errors.New("no users found")
	}

	u, ok := users[0].(*tg.User)
	if !ok {
		return datarealm.User{}, errors.New("error parsing user")
	}

	return tgUserToDataRealmUser(u), nil
}
