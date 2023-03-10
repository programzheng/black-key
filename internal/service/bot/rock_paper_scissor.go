package bot

import "github.com/programzheng/black-key/internal/helper"

func CalculateAmount(groupID string, amount float64) (float64, int) {
	//預設平均計算基數
	amountAvgBase := 3.0
	groupMemberCount := GetGroupMemberCount(groupID)
	amountAvgBase = helper.ConvertToFloat64(groupMemberCount)
	amountAvg := amount / amountAvgBase
	return amountAvg, groupMemberCount
}
