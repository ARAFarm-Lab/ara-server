package usecase

import (
	"ara-server/internal/constants"
	"ara-server/internal/repository/db"
	"ara-server/util/log"
	"context"
	"encoding/json"
	"time"
)

func (uc *Usecase) DispatchScheduler(ctx context.Context) error {
	actions, err := uc.db.GetScheduledAction()
	if err != nil {
		log.Error(ctx, nil, err, "failed to get scheduled actions")
		return err
	}

	// Set all queued actions to pending
	actionIDs := make([]int, len(actions))
	for i, action := range actions {
		actionIDs[i] = action.ID
	}

	if err := uc.db.BulkUpdateActionScheduleStatusPending(actionIDs); err != nil {
		log.Error(ctx, nil, err, "failed to bulk update action schedule status to pending")
		return err
	}

	var errs []error
	for _, action := range actions {
		if err := uc.dispatchAction(ctx, action); err != nil {
			log.Error(ctx, action, err, "failed to dispatch action")
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errorDispatchSchedulerFailed
	}

	return nil
}

func (uc *Usecase) dispatchAction(ctx context.Context, action db.ActionSchedule) (err error) {
	timeNow := time.Now()

	// Release lock, set status and next run time
	defer func() {
		status := constants.ScheduleStatusFailed
		var lastError string

		if err == nil {
			status = constants.ScheduleStatusSuccess

			// Calculate the next run time based on the schedule pattern
			schedulePattern := getSQLNullString(action.Schedule)
			if schedulePattern != "" {
				cronSchedule, err := uc.cronParser.Parse(schedulePattern)
				if err != nil {
					log.Error(ctx, action, err, "failed to parse cron schedule")
				}
				action.NextRunAt = cronSchedule.Next(timeNow)
			}

		} else {
			// Store error message
			lastError = err.Error()
		}

		action.LastLockAt = assignSQLNullTime(nil)
		action.LastRunStatus = assignSQLNullInt(int(status))
		action.LastError = assignSQLNullString(lastError)

		if err := uc.db.UpdateActionSchedule(action); err != nil {
			log.Error(ctx, action, err, "failed to update action schedule")
		}
	}()

	// Set status, lock and last run time
	action.Attempts++
	action.LastRunStatus = assignSQLNullInt(int(constants.ScheduleStatusRunning))
	action.LastLockAt = assignSQLNullTime(&timeNow)
	action.LastRunAt = assignSQLNullTime(&timeNow)
	if err := uc.db.UpdateActionSchedule(action); err != nil {
		log.Error(ctx, action, err, "failed to update action schedule")
		return err
	}

	// Dispatch action
	var actions []DispatcherParam
	if err = json.Unmarshal([]byte(action.Actions), &actions); err != nil {
		log.Error(ctx, action, err, "failed to unmarshal actions")
		return err
	}

	for _, action := range actions {
		if err := uc.DispatchAction(ctx, action); err != nil {
			log.Error(ctx, action, err, "failed to dispatch action")
			return err
		}
	}

	return nil
}
