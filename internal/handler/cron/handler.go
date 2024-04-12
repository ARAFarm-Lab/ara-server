package cron

import (
	"ara-server/internal/constants"
	"ara-server/internal/usecase"
	"context"

	"github.com/robfig/cron/v3"
	"github.com/rs/xid"
)

type handler struct {
	usecase *usecase.Usecase

	cron *cron.Cron
}

func InitHandler(usecase *usecase.Usecase, cron *cron.Cron) {
	handler := &handler{
		cron:    cron,
		usecase: usecase,
	}

	handler.registerHandler()
}

func (h *handler) cronWrapper(handler func(context.Context)) func() {
	return func() {
		ctx := context.Background()
		ctx = context.WithValue(ctx, constants.CtxKeyCtxID, xid.New().String())
		handler(ctx)
	}
}

func (h *handler) registerHandler() {
	h.cron.AddFunc("* * * * *", h.cronWrapper(h.HandleJobDispatcher))
	h.cron.AddFunc("* * * * *", h.cronWrapper(h.HandleInsertDummyData))
	h.cron.AddFunc("*/5 * * * *", h.cronWrapper(h.HandleSendHeartbeatRequest))
}
