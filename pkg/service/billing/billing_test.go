package billing

import (
	"testing"

	"black-key/pkg/helper"
)

func TestAdd(t *testing.T) {
	b := Billing{
		Title:  "測試",
		Amount: 100,
		Payer:  "測試",
		Note:   "test",
	}
	helper.GetJSON(b)

	result, err := b.Add()
	if err != nil {
		t.Fatal("add error:", err)
	}
	t.Log(result)
}
