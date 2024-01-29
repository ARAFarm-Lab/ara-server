package usecase

import (
	"ara-server/internal/constants"
	"ara-server/internal/repository/db"
	"ara-server/util/log"
	"strings"

	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

var (
	errorDispatchScheduledActionFailed = errors.New("failed to dispatch scheduled actions")
)

func (uc *Usecase) CreateSchedule(ctx context.Context, param CreateScheduleParam) error {
	actions, err := json.Marshal(param.Actions)
	if err != nil {
		log.Error(ctx, nil, err, "failed to marshal actions")
		return err
	}

	schedule := db.ActionSchedule{
		Name:        param.Name,
		Description: assignSQLNullString(param.Description),
		Actions:     string(actions),
		NextRunAt:   param.ScheduledAt,
		IsActive:    true,
	}

	if param.RecurringMode != RecurringModeNone {
		schedule.Schedule = assignSQLNullString(buildSchedulePattern(param.ScheduledAt, param.RecurringMode))
	}

	if err := uc.db.InsertActionSchedule(schedule); err != nil {
		log.Error(ctx, schedule, err, "failed to insert action schedule")
		return err
	}

	return nil
}

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

func (uc *Usecase) GetUpcomingSchedules(ctx context.Context) ([]ActionSchedule, error) {
	schedules, err := uc.db.GetUpcomingSchedules()
	if err != nil {
		log.Error(ctx, nil, err, "failed to get upcoming schedules")
		return nil, err
	}

	var result []ActionSchedule
	for _, schedule := range schedules {
		result = append(result, convertUpcomingSchedulesFromDB(ctx, schedule))
	}

	return result, nil
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

	var (
		errs  []error
		wg    sync.WaitGroup
		mutex sync.Mutex
	)
	for _, action := range actions {
		wg.Add(1)
		go func(action DispatcherParam) {
			defer wg.Done()
			if err := uc.DispatchAction(ctx, action); err != nil {
				log.Error(ctx, action, err, "failed to dispatch action")
				mutex.Lock()
				errs = append(errs, err)
				mutex.Unlock()
				return
			}
			if err := uc.insertActionLog(InsertActionLogParam{
				DeviceID:   action.DeviceID,
				ActionType: action.ActionType,
				Value:      action.Value,
				ActionBy:   constants.ActionSourceScheduler,
				ActionAt:   timeNow,
			}); err != nil {
				log.Error(ctx, action, err, "failed to insert action log")
				mutex.Lock()
				errs = append(errs, err)
				mutex.Unlock()
				return
			}
		}(action)
	}

	wg.Wait()
	if len(errs) > 0 {
		log.Error(ctx, errs, errorDispatchScheduledActionFailed, "failed to dispatch scheduled actions")
		return errorDispatchScheduledActionFailed
	}

	return nil
}

func buildSchedulePattern(scheduledAt time.Time, recurringMode ScheduleRecurringMode) string {
	minute := "*"
	hour := "*"

	switch recurringMode {
	case RecurringModeHourly:
		minute = scheduledAt.Format("5")
	case RecurringModeDaily:
		minute = scheduledAt.Format("5")
		hour = scheduledAt.Format("15")
	}

	return minute + " " + hour + " * * *"
}

func convertCronToRecurringMode(schedule string) ScheduleRecurringMode {
	if schedule == "" {
		return RecurringModeNone
	}

	segments := strings.Split(schedule, " ")
	if len(segments) != 5 {
		return RecurringModeNone
	}

	if schedule == "* * * * *" {
		return RecurringModeMinutely
	}

	if segments[0] != "*" && segments[1] == "*" {
		return RecurringModeHourly
	}

	return RecurringModeDaily
}

func convertUpcomingSchedulesFromDB(ctx context.Context, schedule db.ActionSchedule) ActionSchedule {
	var actions []DispatcherParam
	if err := json.Unmarshal([]byte(schedule.Actions), &actions); err != nil {
		log.Error(ctx, schedule, err, "failed to unmarshal actions")
	}

	return ActionSchedule{
		ID:            schedule.ID,
		Name:          schedule.Name,
		Description:   getSQLNullString(schedule.Description),
		Actions:       actions,
		RecurringMode: convertCronToRecurringMode(getSQLNullString(schedule.Schedule)),
		IsActive:      schedule.IsActive,
		ScheduledAt:   schedule.NextRunAt,
		LastRunAt:     schedule.LastRunAt.Time,
		LastRunStatus: constants.ActionScheduleStatus(schedule.LastRunStatus.Int32),
		LastError:     schedule.LastError.String,
	}
}
