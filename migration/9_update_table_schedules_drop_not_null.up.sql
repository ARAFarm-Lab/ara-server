ALTER TABLE schedules
ALTER COLUMN next_run_at DROP NOT NULL; -- Set not null since one time schedule has no next run time after the first trigger