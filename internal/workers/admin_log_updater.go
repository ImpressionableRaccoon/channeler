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
	lastEventID = minID

	resp, err := w.sm.GetChannelAdminLog(ctx, channelID, channelAccessHash, minID)
	if err != nil {
		return lastEventID, err
	}

	if events := resp.GetEvents(); len(events) > 0 {
		// TODO: replace with storager
		fmt.Println("New events:")
		for i, event := range events {
			fmt.Printf("%d: %+v\n", i, event)

			if event.ID > lastEventID {
				lastEventID = event.ID
			}
		}
	}

	if users := resp.GetUsers(); len(users) > 0 {
		// TODO: replace with storager
		fmt.Println("New users:")
		for i, user := range users {
			fmt.Printf("%d: %+v\n", i, user)
		}
	}

	return lastEventID, nil
}
