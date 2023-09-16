package workers

import (
	"go.uber.org/zap"

	"github.com/ImpressionableRaccoon/channeler/internal/datarealm"
)

type workers struct {
	logger *zap.Logger
	sm     datarealm.StatManager
}

func New(logger *zap.Logger, sm datarealm.StatManager) (workers, error) {
	return workers{
		logger: logger,
		sm:     sm,
	}, nil
}
