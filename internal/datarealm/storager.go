package datarealm

import "context"

type Storager interface {
	AddEvent(ctx context.Context, event Event) error
	AddUser(ctx context.Context, user User) error
}
