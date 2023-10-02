package ydb

import (
	"context"
	"database/sql"
	"time"

	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/retry"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
	"go.uber.org/zap"

	"github.com/ImpressionableRaccoon/channeler/internal/datarealm"
)

const upsertEventQuery = `
UPSERT INTO events (
	id, date, user_id, action_type, action
) VALUES (
	$id, $date, $user_id, $action_type, $action
);
`

func (s ydbStorage) AddEvent(ctx context.Context, event datarealm.Event) (err error) {
	if err = ctx.Err(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	s.logger.Info("adding new event", zap.Any("event", event))

	err = retry.Do(ydb.WithTxControl(ctx, writeTx), s.db,
		func(ctx context.Context, cc *sql.Conn) (err error) {
			if err = ctx.Err(); err != nil {
				return err
			}

			_, err = cc.ExecContext(ydb.WithQueryMode(ctx, ydb.DataQueryMode), upsertEventQuery,
				sql.Named("id", event.ID),
				sql.Named("date", event.Date),
				sql.Named("user_id", event.UserID),
				sql.Named("action_type", event.ActionType),
				sql.Named("action", types.JSONValue(event.Action)),
			)
			if err != nil {
				return err
			}

			return err
		},
	)

	return err
}
