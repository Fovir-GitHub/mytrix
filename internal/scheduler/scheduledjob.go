package scheduler

type ScheduledJob struct {
	Cron string
	Job  func()
}
