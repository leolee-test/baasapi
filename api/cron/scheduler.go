package cron

import (
	"github.com/baasapi/baasapi/api"
	"github.com/robfig/cron"
)

// JobScheduler represents a service for managing crons
type JobScheduler struct {
	cron *cron.Cron
}

// NewJobScheduler initializes a new service
func NewJobScheduler() *JobScheduler {
	return &JobScheduler{
		cron: cron.New(),
	}
}

// ScheduleJob schedules the execution of a job via a runner
func (scheduler *JobScheduler) ScheduleJob(runner baasapi.JobRunner) error {
	return scheduler.cron.AddJob(runner.GetSchedule().CronExpression, runner)
}

// UpdateSystemJobSchedule updates the first occurence of the specified
// scheduled job based on the specified job type.
// It does so by re-creating a new cron
// and adding all the existing jobs. It will then re-schedule the new job
// with the update cron expression passed in parameter.
// NOTE: the cron library do not support updating schedules directly
// hence the work-around
func (scheduler *JobScheduler) UpdateSystemJobSchedule(jobType baasapi.JobType, newCronExpression string) error {
	cronEntries := scheduler.cron.Entries()
	newCron := cron.New()

	for _, entry := range cronEntries {
		if entry.Job.(baasapi.JobRunner).GetSchedule().JobType == jobType {
			err := newCron.AddJob(newCronExpression, entry.Job)
			if err != nil {
				return err
			}
			continue
		}

		newCron.Schedule(entry.Schedule, entry.Job)
	}

	scheduler.cron.Stop()
	scheduler.cron = newCron
	scheduler.cron.Start()
	return nil
}

// UpdateJobSchedule updates a specific scheduled job by re-creating a new cron
// and adding all the existing jobs. It will then re-schedule the new job
// via the specified JobRunner parameter.
// NOTE: the cron library do not support updating schedules directly
// hence the work-around
func (scheduler *JobScheduler) UpdateJobSchedule(runner baasapi.JobRunner) error {
	cronEntries := scheduler.cron.Entries()
	newCron := cron.New()

	for _, entry := range cronEntries {

		if entry.Job.(baasapi.JobRunner).GetSchedule().ID == runner.GetSchedule().ID {

			var jobRunner cron.Job = runner
			if entry.Job.(baasapi.JobRunner).GetSchedule().JobType == baasapi.SnapshotJobType {
				jobRunner = entry.Job
			}

			err := newCron.AddJob(runner.GetSchedule().CronExpression, jobRunner)
			if err != nil {
				return err
			}
			continue
		}

		newCron.Schedule(entry.Schedule, entry.Job)
	}

	scheduler.cron.Stop()
	scheduler.cron = newCron
	scheduler.cron.Start()
	return nil
}

// UnscheduleJob remove a scheduled job by re-creating a new cron
// and adding all the existing jobs except for the one specified via scheduleID.
// NOTE: the cron library do not support removing schedules directly
// hence the work-around
func (scheduler *JobScheduler) UnscheduleJob(scheduleID baasapi.ScheduleID) {
	cronEntries := scheduler.cron.Entries()
	newCron := cron.New()

	for _, entry := range cronEntries {

		if entry.Job.(baasapi.JobRunner).GetSchedule().ID == scheduleID {
			continue
		}

		newCron.Schedule(entry.Schedule, entry.Job)
	}

	scheduler.cron.Stop()
	scheduler.cron = newCron
	scheduler.cron.Start()
}

// Start starts the scheduled jobs
func (scheduler *JobScheduler) Start() {
	if len(scheduler.cron.Entries()) > 0 {
		scheduler.cron.Start()
	}
}
