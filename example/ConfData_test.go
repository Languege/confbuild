package example

import (
	"testing"
	"time"
)

/**
 *@author LanguageY++2013
 *2019/3/10 11:52 AM
 **/
func TestTableLevelMaterial_FindByPk(t *testing.T) {
	Start(Configure{
		"./",
		"127.0.0.1",
		"6379",
		"SjhkHD3J5k6H8SjSbK3SC",
		1,
		time.Duration(12) * time.Hour,
		10,
	})
	m, err := TableLevelMaterial_FindByPk(1)
	if err != nil {
		t.FailNow()
	}

	t.Log(m)
}

func TestTableLevelMaterial_ListAll(t *testing.T) {
	Start(Configure{
		"./",
		"127.0.0.1",
		"6379",
		"SjhkHD3J5k6H8SjSbK3SC",
		1,
		time.Duration(12) * time.Hour,
		10,
	})
	m := TableLevelMaterial_ListAll()
	if len(m) == 0 {
		t.FailNow()
	}

	t.Log(m)
}

func TestTableLevelMaterial_ListRange(t *testing.T) {
	Start(Configure{
		"./",
		"127.0.0.1",
		"6379",
		"SjhkHD3J5k6H8SjSbK3SC",
		1,
		time.Duration(12) * time.Hour,
		10,
	})
	TableLevelMaterial_ListRange(func(k uint32, v *TableLevelMaterial) bool {
		if v == nil {
			t.FailNow()
		}

		return true
	})
}


func TestTableLevelMaterial_ListLen(t *testing.T) {
	Start(Configure{
		"./",
		"127.0.0.1",
		"6379",
		"SjhkHD3J5k6H8SjSbK3SC",
		1,
		time.Duration(12) * time.Hour,
		10,
	})

	l := TableLevelMaterial_ListLen()
	if l == 0 {
		t.FailNow()
	}
}