package bot

import (
	"log"

	underscore "github.com/ahl5esoft/golang-underscore"
	"github.com/programzheng/black-key/internal/model/bot"
	"github.com/programzheng/black-key/internal/service"
	"github.com/programzheng/black-key/internal/service/billing"

	"github.com/jinzhu/copier"
)

type LineBilling struct {
	BillingID uint
	GroupID   string
	RoomID    string
	UserID    string

	Billing billing.Billing

	Page service.Page
}

func (lb *LineBilling) Add() (LineBilling, error) {
	model := bot.LineBilling{}
	copier.Copy(&model, &lb)
	result, err := model.Add()
	if err != nil {
		return LineBilling{}, err
	}
	lineBilling := LineBilling{}
	copier.Copy(&lineBilling, &result)

	return lineBilling, nil
}

func (lb *LineBilling) Get(where map[string]interface{}, not map[string]interface{}) ([]LineBilling, error) {
	results, err := bot.LineBilling{}.Get(service.GetDefaultWhere(where), not)
	if err != nil {
		return nil, err
	}
	var lineBillings []LineBilling
	copier.Copy(&lineBillings, &results)
	return lineBillings, nil
}

func BillingAction(lineId LineID, amount int, title string, note string) (billing.Billing, bot.LineBilling) {
	b := billing.Billing{
		Title:  title,
		Amount: amount,
		Note:   note,
	}
	billing, err := b.Add()
	if err != nil {
		log.Fatal("billingAction Billing add error:", err)
	}
	lb := bot.LineBilling{
		BillingID: billing.ID,
		GroupID:   lineId.GroupID,
		RoomID:    lineId.RoomID,
		UserID:    lineId.UserID,
	}
	_, err = lb.Add()
	if err != nil {
		log.Fatal("billingAction LineBilling add error:", err)
	}
	return b, lb
}

func getDistinctByLineBillings(lbs []bot.LineBilling) map[string]string {
	//user id line member display name
	dstByUserID := make(map[string]string, 0)
	underscore.Chain(lbs).DistinctBy("UserID").SelectMany(func(lb bot.LineBilling, _ int) map[string]string {
		dst := make(map[string]string)
		lineMember, err := GetGroupMemberProfile(lb.GroupID, lb.UserID)
		if err != nil {
			dst[lb.UserID] = "Unknow"
			return dst
		}
		dst[lb.UserID] = lineMember.DisplayName
		return dst
	}).Value(&dstByUserID)

	return dstByUserID
}
