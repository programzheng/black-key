package line

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/programzheng/black-key/config"
	model "github.com/programzheng/black-key/internal/model/bot"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestConvertTimeToPushDateTime(t *testing.T) {
	testDateTime := fmt.Sprintf(
		"%s,%s,%s,%s,%s,%s,%s|13:00:00",
		time.Sunday.String(),
		time.Monday.String(),
		time.Tuesday.String(),
		time.Wednesday.String(),
		time.Thursday.String(),
		time.Friday.String(),
		time.Saturday.String(),
	)
	pushDateTime := convertTimeToPushDateTime(testDateTime)
	nowDateTime := time.Now().Format("2006-01-02") + " 13:00:00"
	if pushDateTime != nowDateTime {
		t.Errorf("pushDateTime %s, nowDateTime %s", pushDateTime, nowDateTime)
		t.Fail()
		return
	}
	t.Log("success")
}

func TestRunPushLineNotificationSchedule(t *testing.T) {
	RunPushLineNotificationSchedule()
}

func TestRunPushLineFeatureNotificationSchedule(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	DB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM database: %v", err)
	}

	var lfns []*model.LineFeatureNotification
	query := DB.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Where(nil).Not(nil).Select("feature").Find(&lfns)
	})

	request, _ := json.Marshal(map[string]string{
		"city":    "臺北市",
		"keyword": "new",
	})
	mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{
		"id",
		"created_at",
		"updated_at",
		"deleted_at",
		"feature",
		"push_cycle",
		"push_date_time",
		"limit",
		"group_id",
		"room_id",
		"user_id",
		"request",
	}).AddRow(
		1,
		time.Now(),
		time.Now(),
		nil,
		"new_rent_homes",
		"specify",
		time.Now(),
		-1,
		config.Cfg.GetString("TEST_LINE_PUSH_ID"),
		"",
		"",
		string(request),
	))

	err = DB.Raw(query).Find(&lfns).Error
	if err != nil {
		t.Errorf("DB.Exec(query).Find(&lfns).Error: %v", err)
		return
	}

	runPushLineFeatureNotification(lfns)

	t.Log("TestRunPushLineFeatureNotificationSchedule succeeded!")
}
