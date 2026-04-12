package scheduler

// ScheduledJob represents a job to be executed on a cron schedule.
// It contains the cron expression and the function to execute.
type ScheduledJob struct {
	// Cron is the cron schedule expression (e.g., "0 9 * * *" for 9 AM daily)
	Cron string
	// Job is the function to execute when the schedule triggers
	Job func()
}
