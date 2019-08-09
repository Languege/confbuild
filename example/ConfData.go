
package example

import(
	"errors"
)


func UpdateConfAll() {

	TableLevelMaterial_ListUpdate()

	ChefBasic_ListUpdate()

}

var ErrTableNotExit = errors.New("config table not define")

func UpdateConf(table string) error {
	switch table {
	case "TableLevelMaterial":
		TableLevelMaterial_ListUpdate()
	case "ChefBasic":
		ChefBasic_ListUpdate()
	
	default:
		return ErrTableNotExit
	}

	return nil
}
