package usecase

import (
	"ara-server/internal/constants"
	"ara-server/internal/repository/db"
	"ara-server/util/log"
	"database/sql"
	"strconv"
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

	nextRun := time.Date(param.ScheduledAt.Year(), param.ScheduledAt.Month(), param.ScheduledAt.Day(), param.ScheduledAt.Hour(), param.ScheduledAt.Minute(), 0, 0, param.ScheduledAt.Location())
	schedule := db.ActionSchedule{
		Name:          param.Name,
		Description:   assignSQLNullString(param.Description),
		Actions:       string(actions),
		NextRunAt:     assignSQLNullTime(&nextRun),
		LastRunStatus: assignSQLNullInt(int(constants.ScheduleStatusPending)),
		IsActive:      true,
	}

	if param.DurationInMinutes > 0 {
		schedule.DurationInMinute.Valid = true
		schedule.DurationInMinute.Int32 = int32(param.DurationInMinutes)
	}

	if param.RecurringMode != RecurringModeNone {
		schedulePattern := buildSchedulePattern(param.ScheduledAt, param.RecurringMode)
		schedule.Schedule = assignSQLNullString(schedulePattern)
		cronPattern, err := uc.cronParser.Parse(schedulePattern)
		if err != nil {
			log.Error(ctx, map[string]interface{}{
				"schedule": schedule,
				"pattern":  schedulePattern,
			}, err, "failed to parse cron pattern")
			return err
		}
		nextRunTime := cronPattern.Next(time.Now()) // ensure the next run is after the time now
		nextRun := time.Date(nextRunTime.Year(), nextRunTime.Month(), nextRunTime.Day(), nextRunTime.Hour(), nextRunTime.Minute(), 0, 0, nextRunTime.Location())
		schedule.NextRunAt = assignSQLNullTime(&nextRun)
	}

	if err := uc.db.InsertActionSchedule(schedule); err != nil {
		log.Error(ctx, schedule, err, "failed to insert schedule")
		return err
	}

	return nil
}

func (uc *Usecase) DeleteSchedule(ctx context.Context, scheduleID int) error {
	if err := uc.db.DeleteScheduleByID(ctx, scheduleID); err != nil {
		log.Error(ctx, scheduleID, err, "failed to delete schedule")
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
		result = append(result, convertScheduleFromDB(ctx, schedule))
	}

	return result, nil
}

func (uc *Usecase) UpdateSchedule(ctx context.Context, schedule ActionSchedule) error {
	existingSchedule, err := uc.db.GetScheduleByID(ctx, schedule.ID)
	if err != nil {
		log.Error(ctx, schedule, err, "failed to get schedule by id")
		return err
	}

	if existingSchedule.ID == 0 {
		return constants.ErrorScheduleNotFound
	}

	if err := uc.db.UpdateActionSchedule(ctx, patchSchedule(ctx, existingSchedule, schedule)); err != nil {
		log.Error(ctx, schedule, err, "error patching schedule")
		return err
	}

	return nil
}

func (uc *Usecase) dispatchAction(ctx context.Context, action db.ActionSchedule) (err error) {
	timeNow := time.Now()

	// Set status, lock and last run time
	action.Attempts++
	action.LastRunStatus = assignSQLNullInt(int(constants.ScheduleStatusRunning))
	action.LastLockAt = assignSQLNullTime(&timeNow)
	action.LastRunAt = assignSQLNullTime(&timeNow)
	if err := uc.db.UpdateActionSchedule(ctx, action); err != nil {
		log.Error(ctx, action, err, "failed to update action schedule")
		return err
	}

	// Dispatch action
	var actions []DispatcherParam
	if err = json.Unmarshal([]byte(action.Actions), &actions); err != nil {
		log.Error(ctx, action, err, "failed to unmarshal actions")
		return err
	}

	// If it is time to clean up the schedule, revert the action
	hasDuration := action.DurationInMinute.Valid
	cleanUpTime := action.CleanupTime.Time
	isCleanUpTime := action.CleanupTime.Valid && (timeNow.Equal(cleanUpTime) || timeNow.After(cleanUpTime))

	var (
		errs  []error
		wg    sync.WaitGroup
		mutex sync.Mutex
	)
	for _, action := range actions {
		wg.Add(1)
		if hasDuration && isCleanUpTime {
			switch t := action.Value.(type) {
			case bool:
				action.Value = !t
			}
		}

		go func(action DispatcherParam) {
			defer wg.Done()

			action.ActionBy = constants.ActionSourceScheduler
			if err := uc.DispatchAction(ctx, action); err != nil {
				log.Error(ctx, action, err, "failed to dispatch action")
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

	// Release lock, set status and next run time
	defer func() {
		status := constants.ScheduleStatusFailed
		var lastError string

		if err == nil {
			status = constants.ScheduleStatusSuccess

			// Calculate the next run time and the cleanup time based on the schedule pattern
			if action.Schedule.Valid {
				schedulePattern := getSQLNullString(action.Schedule)
				cronSchedule, err := uc.cronParser.Parse(schedulePattern)
				if err != nil {
					log.Error(ctx, action, err, "failed to parse cron schedule")
				}

				// Set next run time if it is clean up time
				nextRunTime := cronSchedule.Next(timeNow)
				if !isCleanUpTime {
					action.NextRunAt = assignSQLNullTime(&nextRunTime)
					action.CleanupTime = assignSQLNullTime(calculateCleanupTime(action, timeNow))
				} else {
					action.CleanupTime = assignSQLNullTime(calculateCleanupTime(action, nextRunTime))
				}
			} else {
				if isCleanUpTime {
					// Delete the schedule after clean up since the schedule is no longer used
					if err = uc.db.DeleteScheduleByID(ctx, action.ID); err != nil {
						log.Error(ctx, action, err, "failed to delete schedule")
					}
				} else {
					action.NextRunAt = assignSQLNullTime(nil)
					action.CleanupTime = assignSQLNullTime(calculateCleanupTime(action, timeNow))
				}
			}
		} else {
			// Store error message
			lastError = err.Error()
		}

		action.LastLockAt = assignSQLNullTime(nil)
		action.LastRunStatus = assignSQLNullInt(int(status))
		action.LastError = assignSQLNullString(lastError)

		if err := uc.db.UpdateActionSchedule(ctx, action); err != nil {
			log.Error(ctx, action, err, "failed to update action schedule")
		}
	}()

	return nil
}

func buildSchedulePattern(scheduledAt time.Time, recurringMode ScheduleRecurringMode) string {
	minute := "*"
	hour := "*"

	switch recurringMode {
	case RecurringModeHourly:
		minute = scheduledAt.Format("4")
	case RecurringModeDaily:
		minute = scheduledAt.Format("4")
		hour = scheduledAt.Format("15")
	}

	return minute + " " + hour + " * * *"
}

func calculateCleanupTime(action db.ActionSchedule, timeNow time.Time) *time.Time {
	if !action.DurationInMinute.Valid {
		return nil
	}

	duration := action.DurationInMinute.Int32
	cleanUpTime := timeNow.Add(time.Duration(duration) * time.Minute)
	return &cleanUpTime
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

func convertScheduleFromDB(ctx context.Context, schedule db.ActionSchedule) ActionSchedule {
	var actions []DispatcherParam
	if err := json.Unmarshal([]byte(schedule.Actions), &actions); err != nil {
		log.Error(ctx, schedule, err, "failed to unmarshal actions")
	}

	result := ActionSchedule{
		ID:            schedule.ID,
		Name:          schedule.Name,
		Description:   getSQLNullString(schedule.Description),
		Actions:       actions,
		RecurringMode: convertCronToRecurringMode(getSQLNullString(schedule.Schedule)),
		Duration:      int(schedule.DurationInMinute.Int32),
		IsActive:      schedule.IsActive,
		ScheduledAt:   generateScheduleTime(schedule),
		LastRunStatus: constants.ActionScheduleStatus(schedule.LastRunStatus.Int32),
		LastError:     schedule.LastError.String,
	}

	if !schedule.LastRunAt.Time.IsZero() {
		result.LastRunAt = &schedule.LastRunAt.Time
	}

	if schedule.NextRunAt.Valid {
		result.NextRunAt = schedule.NextRunAt.Time
	}

	if schedule.CleanupTime.Valid && schedule.CleanupTime.Time.Before(schedule.NextRunAt.Time) || !schedule.NextRunAt.Valid {
		result.NextRunAt = schedule.CleanupTime.Time
	}

	return result
}

func generateScheduleTime(schedule db.ActionSchedule) time.Time {
	if !schedule.NextRunAt.Valid {
		return time.Time{}
	}

	pattern := getSQLNullString(schedule.Schedule)
	recurringMode := convertCronToRecurringMode(pattern)
	if recurringMode == RecurringModeNone || recurringMode == RecurringModeMinutely {
		return schedule.NextRunAt.Time
	}

	segments := strings.Split(pattern, " ")
	minute, _ := strconv.Atoi(segments[0])
	t := schedule.NextRunAt.Time
	h := t.Hour()
	if recurringMode == RecurringModeDaily {
		hour, _ := strconv.Atoi(segments[1])
		h = hour
	}

	return time.Date(t.Year(), t.Month(), t.Day(), h, minute, t.Second(), t.Nanosecond(), t.Location())
}

func patchSchedule(ctx context.Context, old db.ActionSchedule, new ActionSchedule) db.ActionSchedule {
	old.Name = new.Name

	if new.Description != "" {
		old.Description = sql.NullString{
			Valid:  true,
			String: new.Description,
		}
	}

	b, err := json.Marshal(new.Actions)
	if err != nil {
		log.Error(ctx, new, err, "error marshalling actions")
	}

	if !new.ScheduledAt.IsZero() {
		if new.RecurringMode != RecurringModeNone {
			old.Schedule = sql.NullString{
				Valid:  true,
				String: buildSchedulePattern(new.ScheduledAt, new.RecurringMode),
			}
		}

		old.NextRunAt = assignSQLNullTime(&new.ScheduledAt)
	}

	if new.Duration > 0 {
		old.DurationInMinute = sql.NullInt32{
			Valid: true,
			Int32: int32(new.Duration),
		}
	}

	old.Actions = string(b)

	return old
}
