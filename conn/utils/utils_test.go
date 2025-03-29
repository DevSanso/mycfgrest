package utils_test

import (
	"mycfgrest/conn/utils"
	"mycfgrest/types"
	"testing"
)
 
func TestSimpleConvertSql(t *testing.T) {
	numberMap := types.NewParsingMap()

	queryTemplate := "SELECT #{request.hello}, ###name , #{0.0.name}"
	queryNum := "SELECT $1, ##name , $2"
	queryQuestion := "SELECT ?, ##name , ?"

	numberMap.Set(0, "request.hello", "hello", types.STRING)
	numberMap.Set(0, "0.0.name", 1234, types.INT)

	fetcher, fetchErr := numberMap.Fetch()
	if fetchErr != nil {
		t.Error(fetchErr.Error())
	}
	numq, _, numErr := utils.ChangeSqlToNumBindSupportSql(queryTemplate, fetcher)
	qq, _, qErr := utils.ChangeSqlToQuestionMarkBindSupportSql(queryTemplate, fetcher)

	if numErr != nil || qErr != nil{
		t.Errorf("number query bind error => %s", numErr.Error())
		t.Errorf("question query bind error => %s", qErr.Error())
		return
	}

	if numq != queryNum {
		t.Errorf("num query => '%s' != '%s'", numq, queryNum)
	}
	
	if qq != queryQuestion {
		t.Errorf("question query => '%s' != '%s'", qq, queryQuestion)
	}
}