package system

import (
	"context"
	sysModel "github.com/EscapeBearSecond/curescan/server/model/system"
	"github.com/EscapeBearSecond/curescan/server/service/system"
	"github.com/EscapeBearSecond/curescan/server/utils"
	"github.com/gofrs/uuid/v5"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderUser = initOrderAuthority + 1

type initUser struct{}

// auto run
func init() {
	system.RegisterInit(initOrderUser, &initUser{})
}

func (i *initUser) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&sysModel.SysUser{})
}

func (i *initUser) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&sysModel.SysUser{})
}

func (i initUser) InitializerName() string {
	return sysModel.SysUser{}.TableName()
}

func (i *initUser) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	ap := ctx.Value("adminPassword")
	apStr, ok := ap.(string)
	if !ok {
		apStr = "123456"
	}

	// password := utils.BcryptHash(apStr)
	adminPassword := utils.BcryptHash(apStr)

	entities := []sysModel.SysUser{
		{
			UUID:        uuid.Must(uuid.NewV4()),
			Username:    "superAdmin",
			Password:    adminPassword,
			NickName:    "DYG",
			HeaderImg:   "",
			AuthorityId: 888,
			Phone:       "",
			Email:       "",
		},
		{
			UUID:        uuid.Must(uuid.NewV4()),
			Username:    "admin",
			Password:    utils.BcryptHash("@123456qwer"),
			NickName:    "管理员",
			HeaderImg:   "",
			AuthorityId: 9528,
			Phone:       "",
			Email:       "",
		},
		{
			UUID:        uuid.Must(uuid.NewV4()),
			Username:    "systemtester",
			Password:    utils.BcryptHash("@123456qwer"),
			NickName:    "测试员",
			HeaderImg:   "",
			AuthorityId: 1913,
			Phone:       "",
			Email:       "",
		},
		{
			UUID:        uuid.Must(uuid.NewV4()),
			Username:    "systemauditor",
			Password:    utils.BcryptHash("@123456qwer"),
			NickName:    "审计员",
			HeaderImg:   "",
			AuthorityId: 1914,
			Phone:       "",
			Email:       "",
		},
	}
	if err = db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, sysModel.SysUser{}.TableName()+"表数据初始化失败!")
	}
	next = context.WithValue(ctx, i.InitializerName(), entities)
	authorityEntities, ok := ctx.Value(initAuthority{}.InitializerName()).([]sysModel.SysAuthority)
	if !ok {
		return next, errors.Wrap(system.ErrMissingDependentContext, "创建 [用户-权限] 关联失败, 未找到权限表初始化数据")
	}
	if err = db.Model(&entities[0]).Association("Authorities").Replace(authorityEntities); err != nil {
		return next, err
	}
	if err = db.Model(&entities[1]).Association("Authorities").Replace(authorityEntities[1:2]); err != nil {
		return next, err
	}
	if err = db.Model(&entities[2]).Association("Authorities").Replace(authorityEntities[2:3]); err != nil {
		return next, err
	}
	if err = db.Model(&entities[3]).Association("Authorities").Replace(authorityEntities[3:4]); err != nil {
		return next, err
	}
	return next, err
}

func (i *initUser) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	var record sysModel.SysUser
	if errors.Is(db.Where("username = ?", "admin").
		Preload("Authorities").First(&record).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return len(record.Authorities) > 0 && record.Authorities[0].AuthorityId == 9528
}
