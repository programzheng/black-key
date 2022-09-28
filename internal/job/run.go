package job

import (
	"time"

	"github.com/go-co-op/gocron"
	"github.com/programzheng/black-key/internal/job/line"
)

func Run() {
	s := gocron.NewScheduler(time.Now().Local().Location())
	s.Cron("*/1 * * * *").Do(func() {
		line.RunSchedule()
	}) // every minute
	s.StartAsync()
}
