package system

import (
	sysModel "47.103.136.241/goprojects/curescan/server/model/system"
	"47.103.136.241/goprojects/curescan/server/service/system"
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderMenuAuthority = initOrderMenu + initOrderAuthority

type initMenuAuthority struct{}

// auto run
func init() {
	system.RegisterInit(initOrderMenuAuthority, &initMenuAuthority{})
}

func (i *initMenuAuthority) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil // do nothing
}

func (i *initMenuAuthority) TableCreated(ctx context.Context) bool {
	return false // always replace
}

func (i initMenuAuthority) InitializerName() string {
	return "sys_menu_authorities"
}

func (i *initMenuAuthority) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	authorities, ok := ctx.Value(initAuthority{}.InitializerName()).([]sysModel.SysAuthority)
	if !ok {
		return ctx, errors.Wrap(system.ErrMissingDependentContext, "创建 [菜单-权限] 关联失败, 未找到权限表初始化数据")
	}
	menus, ok := ctx.Value(initMenu{}.InitializerName()).([]sysModel.SysBaseMenu)
	if !ok {
		return next, errors.Wrap(errors.New(""), "创建 [菜单-权限] 关联失败, 未找到菜单表初始化数据")
	}
	next = ctx
	// 888
	if err = db.Model(&authorities[0]).Association("SysBaseMenus").Replace(menus); err != nil {
		return next, err
	}
	menus1 := make([]sysModel.SysBaseMenu, len(menus))
	menus2 := make([]sysModel.SysBaseMenu, len(menus))
	copy(menus1, menus)
	copy(menus2, menus)
	// 9528
	menu9528 := menus[:1]
	menu9528 = append(menu9528, menus[2], menus[6], menus[9], menus[22])
	menu9528 = append(menu9528, menus[30:36]...)
	menu9528 = append(menu9528, menus[38:48]...)
	if err = db.Model(&authorities[1]).Association("SysBaseMenus").Replace(menu9528); err != nil {
		return next, err
	}
	// 1913
	menu1913 := menus1[:1]
	menu1913 = append(menu1913, menus1[9])
	menu1913 = append(menu1913, menus1[30:36]...)
	menu1913 = append(menu1913, menus1[38:48]...)
	if err = db.Model(&authorities[2]).Association("SysBaseMenus").Replace(menu1913); err != nil {
		return next, err
	}
	// 1914
	menu1914 := menus2[:1]
	menu1914 = append(menu1914, menus2[9])
	menu1914 = append(menu1914, menus2[30:36]...)
	menu1914 = append(menu1914, menus2[38:48]...)
	if err = db.Model(&authorities[3]).Association("SysBaseMenus").Replace(menu1914); err != nil {
		return next, err
	}
	// 9528
	// if err = db.Model(&authorities[2]).Association("SysBaseMenus").Replace(menus[:11]); err != nil {
	// 	return next, err
	// }
	// if err = db.Model(&authorities[2]).Association("SysBaseMenus").Append(menus[12:17]); err != nil {
	// 	return next, err
	// }
	return next, nil
}

func (i *initMenuAuthority) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	auth := &sysModel.SysAuthority{}
	if ret := db.Model(auth).
		Where("authority_id = ?", 9528).Preload("SysBaseMenus").Find(auth); ret != nil {
		if ret.Error != nil {
			return false
		}
		return len(auth.SysBaseMenus) > 0
	}
	return false
}
