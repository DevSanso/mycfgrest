package types_test

import (
	"mycfgrest/types"
	"testing"
)

func TestParsingMap(t *testing.T) {
	m := types.NewParsingMap()

	if err := m.Set(0, "name", "name", types.STRING); err != nil {
		t.Error(err.Error())
		return
	}
	m.Set(0, "value", 10, types.INT)
	m.Set(2, "hello", nil, types.NULL)

	fetch, _ := m.Fetch()
	{
		key, val, gt, err := fetch.GetData()
		if err != nil {
			t.Error(err.Error())
			return
		}
		if key[0] != "name" || val[0] != "name" || gt[0] != types.STRING {
			t.Errorf("not eq, name=%s:%s val=%s:%s type=%s:%s", key[0], "name", val[0],"name", gt[0], types.STRING)
			return
		}

		if key[1] != "name" || val[1] != "name" || gt[1] != types.STRING {
			t.Errorf("not eq, name=%s:%s val=%d:%d type=%s:%s", key[1], "name", val[1],10, gt[1], types.INT)
			return
		}
	}
	fetch.Next()
	fetch.Next()
	{
		key, val, gt, err := fetch.GetData()
		if err != nil {
			t.Error(err.Error())
			return
		}
		if key[0] != "hello" || val[0] != nil || gt[0] != types.NULL {
			t.Errorf("not eq, name=%s:%s val=%t type=%s:%s", key[0], "hello", val[0] == nil, gt[0], types.NULL)
			return
		}
	}

}
 