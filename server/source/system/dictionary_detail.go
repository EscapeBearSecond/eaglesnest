package system

import (
	"context"
	"fmt"

	sysModel "47.103.136.241/goprojects/curescan/server/model/system"
	"47.103.136.241/goprojects/curescan/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderDictDetail = initOrderDict + 1

type initDictDetail struct{}

// auto run
func init() {
	system.RegisterInit(initOrderDictDetail, &initDictDetail{})
}

func (i *initDictDetail) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&sysModel.SysDictionaryDetail{})
}

func (i *initDictDetail) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&sysModel.SysDictionaryDetail{})
}

func (i initDictDetail) InitializerName() string {
	return sysModel.SysDictionaryDetail{}.TableName()
}

func (i *initDictDetail) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	dicts, ok := ctx.Value(initDict{}.InitializerName()).([]sysModel.SysDictionary)
	if !ok {
		return ctx, errors.Wrap(system.ErrMissingDependentContext,
			fmt.Sprintf("未找到 %s 表初始化数据", sysModel.SysDictionary{}.TableName()))
	}
	True := true
	dicts[0].SysDictionaryDetails = []sysModel.SysDictionaryDetail{
		{Label: "男", Value: "1", Status: &True, Sort: 1},
		{Label: "女", Value: "2", Status: &True, Sort: 2},
	}

	dicts[1].SysDictionaryDetails = []sysModel.SysDictionaryDetail{
		{Label: "smallint", Value: "1", Status: &True, Extend: "mysql", Sort: 1},
		{Label: "mediumint", Value: "2", Status: &True, Extend: "mysql", Sort: 2},
		{Label: "int", Value: "3", Status: &True, Extend: "mysql", Sort: 3},
		{Label: "bigint", Value: "4", Status: &True, Extend: "mysql", Sort: 4},
		{Label: "int2", Value: "5", Status: &True, Extend: "pgsql", Sort: 5},
		{Label: "int4", Value: "6", Status: &True, Extend: "pgsql", Sort: 6},
		{Label: "int6", Value: "7", Status: &True, Extend: "pgsql", Sort: 7},
		{Label: "int8", Value: "8", Status: &True, Extend: "pgsql", Sort: 8},
	}

	dicts[2].SysDictionaryDetails = []sysModel.SysDictionaryDetail{
		{Label: "date", Status: &True},
		{Label: "time", Value: "1", Status: &True, Extend: "mysql", Sort: 1},
		{Label: "year", Value: "2", Status: &True, Extend: "mysql", Sort: 2},
		{Label: "datetime", Value: "3", Status: &True, Extend: "mysql", Sort: 3},
		{Label: "timestamp", Value: "5", Status: &True, Extend: "mysql", Sort: 5},
		{Label: "timestamptz", Value: "6", Status: &True, Extend: "pgsql", Sort: 5},
	}
	dicts[3].SysDictionaryDetails = []sysModel.SysDictionaryDetail{
		{Label: "float", Status: &True},
		{Label: "double", Value: "1", Status: &True, Extend: "mysql", Sort: 1},
		{Label: "decimal", Value: "2", Status: &True, Extend: "mysql", Sort: 2},
		{Label: "numeric", Value: "3", Status: &True, Extend: "pgsql", Sort: 3},
		{Label: "smallserial", Value: "4", Status: &True, Extend: "pgsql", Sort: 4},
	}

	dicts[4].SysDictionaryDetails = []sysModel.SysDictionaryDetail{
		{Label: "char", Status: &True},
		{Label: "varchar", Value: "1", Status: &True, Extend: "mysql", Sort: 1},
		{Label: "tinyblob", Value: "2", Status: &True, Extend: "mysql", Sort: 2},
		{Label: "tinytext", Value: "3", Status: &True, Extend: "mysql", Sort: 3},
		{Label: "text", Value: "4", Status: &True, Extend: "mysql", Sort: 4},
		{Label: "blob", Value: "5", Status: &True, Extend: "mysql", Sort: 5},
		{Label: "mediumblob", Value: "6", Status: &True, Extend: "mysql", Sort: 6},
		{Label: "mediumtext", Value: "7", Status: &True, Extend: "mysql", Sort: 7},
		{Label: "longblob", Value: "8", Status: &True, Extend: "mysql", Sort: 8},
		{Label: "longtext", Value: "9", Status: &True, Extend: "mysql", Sort: 9},
	}

	dicts[5].SysDictionaryDetails = []sysModel.SysDictionaryDetail{
		{Label: "tinyint", Value: "1", Extend: "mysql", Status: &True},
		{Label: "bool", Value: "2", Extend: "pgsql", Status: &True},
	}
	dicts[6].SysDictionaryDetails = []sysModel.SysDictionaryDetail{
		{Label: "立即执行", Value: "1", Status: &True, Sort: 1},
		{Label: "稍后执行", Value: "2", Status: &True, Sort: 2},
		{Label: "定时执行", Value: "3", Status: &True, Sort: 3},
	}
	dicts[7].SysDictionaryDetails = []sysModel.SysDictionaryDetail{
		{Label: "资产发现", Value: "1", Status: &True, Sort: 1},
		{Label: "漏洞扫描", Value: "2", Status: &True, Sort: 2},
		{Label: "弱口令", Value: "3", Status: &True, Sort: 3},
	}
	dicts[8].SysDictionaryDetails = []sysModel.SysDictionaryDetail{
		{Label: "创建", Value: "7", Status: &True, Extend: "普通任务", Sort: 1},
		{Label: "执行中", Value: "1", Status: &True, Extend: "普通任务", Sort: 2},
		{Label: "成功", Value: "2", Status: &True, Extend: "普通任务", Sort: 3},
		{Label: "失败", Value: "3", Status: &True, Extend: "普通任务", Sort: 4},
		{Label: "中止", Value: "4", Status: &True, Extend: "普通任务", Sort: 5},
		{Label: "运行", Value: "5", Status: &True, Extend: "普通任务", Sort: 6},
		{Label: "停止", Value: "6", Status: &True, Extend: "普通任务", Sort: 7},
	}
	dicts[9].SysDictionaryDetails = []sysModel.SysDictionaryDetail{
		{Label: "严重", Value: "critical", Status: &True, Sort: 1},
		{Label: "高危", Value: "high", Status: &True, Sort: 2},
		{Label: "中危", Value: "medium", Status: &True, Sort: 3},
		{Label: "低危", Value: "low", Status: &True, Sort: 4},
	}

	for _, dict := range dicts {
		if err := db.Model(&dict).Association("SysDictionaryDetails").
			Replace(dict.SysDictionaryDetails); err != nil {
			return ctx, errors.Wrap(err, sysModel.SysDictionaryDetail{}.TableName()+"表数据初始化失败!")
		}
	}
	return ctx, nil
}

func (i *initDictDetail) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	var dict sysModel.SysDictionary
	if err := db.Preload("SysDictionaryDetails").
		First(&dict, &sysModel.SysDictionary{Name: "数据库bool类型"}).Error; err != nil {
		return false
	}
	return len(dict.SysDictionaryDetails) > 0 && dict.SysDictionaryDetails[0].Label == "tinyint"
}
