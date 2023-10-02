package stats

import (
	"github.com/gotd/td/tg"

	"github.com/ImpressionableRaccoon/channeler/internal/datarealm"
)

func tgUserToDataRealmUser(u *tg.User) datarealm.User {
	return datarealm.User{
		ID:         u.GetID(),
		AccessHash: u.AccessHash,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		Username:   u.Username,
		Phone:      u.Phone,
	}
}
