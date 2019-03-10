package example

import "testing"

/**
 *@author LanguageY++2013
 *2019/3/10 11:52 AM
 **/
func TestTableLevelMaterial_FindByPk(t *testing.T) {
	m, err := TableLevelMaterial_FindByPk(1)
	if err != nil {
		t.FailNow()
	}

	t.Log(m)
}

func TestTableLevelMaterial_ListAll(t *testing.T) {
	m := TableLevelMaterial_ListAll()
	if len(m) == 0 {
		t.FailNow()
	}

	t.Log(m)
}

func TestTableLevelMaterial_ListRange(t *testing.T) {
	TableLevelMaterial_ListRange(func(k uint32, v *TableLevelMaterial) bool {
		if v == nil {
			t.FailNow()
		}

		return true
	})
}


func TestTableLevelMaterial_ListLen(t *testing.T) {
	l := TableLevelMaterial_ListLen()
	if l == 0 {
		t.FailNow()
	}
}