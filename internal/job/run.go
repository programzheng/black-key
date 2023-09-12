package job

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/programzheng/black-key/config"
	"github.com/programzheng/black-key/internal/job/line"
)

func Run() {
	s := gocron.NewScheduler(time.Now().Local().Location())
	s.Cron("*/1 * * * *").Do(func() {
		line.RunPushLineNotificationSchedule()
		if config.Cfg.GetBool("JOBS_DEBUG") {
			fmt.Printf("The job %s is scheduled to run at Cron: %s\n", "line.RunPushLineNotificationSchedule()", "*/1 * * * *")
		}
	}) // every minute
	s.Cron("*/1 * * * *").Do(func() {
		line.RunPushLineFeatureNotificationSchedule()
		if config.Cfg.GetBool("JOBS_DEBUG") {
			fmt.Printf("The job %s is scheduled to run at Cron: %s\n", "line.RunPushLineFeatureNotificationSchedule()", "*/1 * * * *")
		}
	}) // every minute
	s.Cron("0 0 * * *").Do(func() {
		line.RunRefreshLineNotificationSchedule()
		if config.Cfg.GetBool("JOBS_DEBUG") {
			fmt.Printf("The job %s is scheduled to run at Cron: %s\n", "line.RunRefreshLineNotificationSchedule()", "0 0 * * *")
		}
	}) // every daily
	s.StartAsync()
}
