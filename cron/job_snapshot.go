package cron

import (
	"log"

	"github.com/baasapi/baasapi/api"
)

// SnapshotJobRunner is used to run a SnapshotJob
type SnapshotJobRunner struct {
	schedule *baasapi.Schedule
	context  *SnapshotJobContext
}

// SnapshotJobContext represents the context of execution of a SnapshotJob
type SnapshotJobContext struct {
	baask8sService baasapi.Baask8sService
	snapshotter     baasapi.Snapshotter
}

// NewSnapshotJobContext returns a new context that can be used to execute a SnapshotJob
func NewSnapshotJobContext(baask8sService baasapi.Baask8sService, snapshotter baasapi.Snapshotter) *SnapshotJobContext {
	return &SnapshotJobContext{
		baask8sService: baask8sService,
		snapshotter:     snapshotter,
	}
}

// NewSnapshotJobRunner returns a new runner that can be scheduled
func NewSnapshotJobRunner(schedule *baasapi.Schedule, context *SnapshotJobContext) *SnapshotJobRunner {
	return &SnapshotJobRunner{
		schedule: schedule,
		context:  context,
	}
}

// GetSchedule returns the schedule associated to the runner
func (runner *SnapshotJobRunner) GetSchedule() *baasapi.Schedule {
	return runner.schedule
}

// Run triggers the execution of the schedule.
// It will iterate through all the baask8ss available in the database to
// create a snapshot of each one of them.
// As a snapshot can be a long process, to avoid any concurrency issue we
// retrieve the latest version of the baask8s right after a snapshot.
func (runner *SnapshotJobRunner) Run() {
	go func() {

		log.Printf("background schedule running.  (=%s)\n", "test...")

		baask8ss, err := runner.context.baask8sService.Baask8ss()
		if err != nil {
			log.Printf("background schedule error (baask8s snapshot). Unable to retrieve baask8s list (err=%s)\n", err)
			return
		}

		for _, baask8s := range baask8ss {

			snapshot, snapshotError := runner.context.snapshotter.CreateSnapshot(&baask8s)

			latestBaask8sReference, err := runner.context.baask8sService.Baask8s(baask8s.ID)
			if latestBaask8sReference == nil {
				log.Printf("background schedule error (baask8s snapshot). Baask8s not found inside the database anymore (baask8s=%s, NetworkID=%s) (err=%s)\n", baask8s.NetworkName, baask8s.NetworkID, err)
				continue
			}

			latestBaask8sReference.Status = baasapi.Baask8sStatusUp
			if snapshotError != nil {
				log.Printf("background schedule error (baask8s snapshot). Unable to create snapshot (baask8s=%s, NetworkID=%s) (err=%s)\n", baask8s.NetworkName, baask8s.NetworkID, snapshotError)
				latestBaask8sReference.Status = baasapi.Baask8sStatusDown
			}

			if snapshot != nil {
				latestBaask8sReference.Snapshots = []baasapi.Snapshot{*snapshot}
			}

			err = runner.context.baask8sService.UpdateBaask8s(latestBaask8sReference.ID, latestBaask8sReference)
			if err != nil {
				log.Printf("background schedule error (baask8s snapshot). Unable to update baask8s (baask8s=%s, NetworkID=%s) (err=%s)\n", baask8s.NetworkName, baask8s.NetworkID, err)
				return
			}
		}
	}()
}
