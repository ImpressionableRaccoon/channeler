package workers

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

func (w workers) AdminLogUpdater(ctx context.Context, interval time.Duration, channelID, channelAccessHash int64) {
	w.logger.Info("AdminLogUpdater is starting",
		zap.Int64("channel id", channelID),
		zap.Int64("channel access hash", channelAccessHash),
		zap.Duration("interval", interval),
	)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var lastEventID int64
	var err error

	for {
		lastEventID, err = w.updateAdminLog(ctx, channelID, channelAccessHash, lastEventID)
		if err != nil {
			w.logger.Error("update admin log failed", zap.Error(err))
		}

		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			return
		}
	}
}

func (w workers) updateAdminLog(ctx context.Context, channelID, channelAccessHash, minID int64) (
	lastEventID int64, err error,
) {
	if err = ctx.Err(); err != nil {
		return minID, err
	}

	resp, err := w.sm.GetChannelAdminLog(ctx, channelID, channelAccessHash, minID)
	if err != nil {
		return lastEventID, fmt.Errorf("error while getting channel admin log: %w", err)
	}

	var saveError bool
	maxEventID := minID

	if events := resp.GetEvents(); len(events) > 0 {
		for _, event := range events {
			err = w.st.AddEvent(ctx, event)
			if err != nil {
				w.logger.Error("error saving event", zap.Error(err))
				saveError = true
			}
			if event.ID > maxEventID {
				maxEventID = event.ID
			}
		}
	}

	if users := resp.GetUsers(); len(users) > 0 {
		for _, user := range users {
			err = w.st.AddUser(ctx, user)
			if err != nil {
				w.logger.Error("error saving user", zap.Error(err))
				saveError = true
			}
		}
	}

	if saveError {
		w.logger.Warn("error occurred while saving, will not update lastEventID")
		lastEventID = minID
	} else {
		lastEventID = maxEventID
	}

	return lastEventID, nil
}
