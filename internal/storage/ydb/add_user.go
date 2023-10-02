package ydb

import (
	"database/sql"
	"time"

	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/retry"
	"go.uber.org/zap"
	"golang.org/x/net/context"

	"github.com/ImpressionableRaccoon/channeler/internal/datarealm"
)

const upsertUserQuery = `
UPSERT INTO users (
	id, access_hash, first_name, last_name, username, phone
) VALUES (
	$id, $access_hash, $first_name, $last_name, $username, $phone
);
`

func (s ydbStorage) AddUser(ctx context.Context, user datarealm.User) (err error) {
	if err = ctx.Err(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	s.logger.Info("adding new user", zap.Any("user", user))

	err = retry.Do(ydb.WithTxControl(ctx, writeTx), s.db,
		func(ctx context.Context, cc *sql.Conn) (err error) {
			if err = ctx.Err(); err != nil {
				return err
			}

			_, err = cc.ExecContext(ydb.WithQueryMode(ctx, ydb.DataQueryMode), upsertUserQuery,
				sql.Named("id", user.ID),
				sql.Named("access_hash", user.AccessHash),
				sql.Named("first_name", user.FirstName),
				sql.Named("last_name", user.LastName),
				sql.Named("username", user.Username),
				sql.Named("phone", user.Phone),
			)
			if err != nil {
				return err
			}

			return err
		},
	)

	return err
}
