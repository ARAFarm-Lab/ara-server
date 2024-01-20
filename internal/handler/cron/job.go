package cron

import (
	"ara-server/util/log"
	"context"
)

func (h *handler) HandleJobDispatcher(ctx context.Context) {
	log.Info(ctx, nil, nil, "start handling job dispatcher")
	if err := h.usecase.DispatchScheduler(ctx); err != nil {
		log.Error(ctx, nil, nil, err.Error())
	}
	log.Info(ctx, nil, nil, "finish handling job dispatcher")
}
