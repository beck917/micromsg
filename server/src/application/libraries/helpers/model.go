package helpers

import (
	"application/libraries/toml"
)

type Model struct {
	DB        *Database
	TableName string
}

func (this *Model) Construct(tableName string) {
	this.DB = InstanceDatabase(toml.GlobalTomlConfig.Mysql)

	this.TableName = tableName
}

func (this *Model) Update(id interface{}, updateEntity interface{}) (err error) {
	//更新表
	_, err = this.DB.XORM.Id(id).Update(updateEntity)

	return err
}
