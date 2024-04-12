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

func (h *handler) HandleInsertDummyData(ctx context.Context) {
	log.Info(ctx, nil, nil, "start handling insert dummy data")
	if err := h.usecase.InsertDummyData(); err != nil {
		log.Error(ctx, nil, nil, err.Error())
	}
	log.Info(ctx, nil, nil, "finish handling insert dummy data")
}

func (h *handler) HandleSendHeartbeatRequest(ctx context.Context) {
	log.Info(ctx, nil, nil, "start handling send heartbeat request")
	h.usecase.SendHeartbeat(ctx, 1)
	log.Info(ctx, nil, nil, "finish handling send heartbeat request")
}
