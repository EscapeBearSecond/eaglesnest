package system

import (
	"context"

	. "47.103.136.241/goprojects/curescan/server/model/system"
	"47.103.136.241/goprojects/curescan/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderMenu = initOrderAuthority + 1

type initMenu struct{}

// auto run
func init() {
	system.RegisterInit(initOrderMenu, &initMenu{})
}

func (i initMenu) InitializerName() string {
	return SysBaseMenu{}.TableName()
}

func (i *initMenu) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&SysBaseMenu{},
		&SysBaseMenuParameter{},
		&SysBaseMenuBtn{},
	)
}

func (i *initMenu) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	m := db.Migrator()
	return m.HasTable(&SysBaseMenu{}) &&
		m.HasTable(&SysBaseMenuParameter{}) &&
		m.HasTable(&SysBaseMenuBtn{})
}

func (i *initMenu) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	entities := []SysBaseMenu{
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "dashboard", Name: "dashboard", Component: "view/dashboard/index.vue", Sort: 1, Meta: Meta{Title: "首页", Icon: "odometer"}},
		{MenuLevel: 0, Hidden: true, ParentId: 0, Path: "about", Name: "about", Component: "view/about/index.vue", Sort: 9, Meta: Meta{Title: "关于我们", Icon: "info-filled"}},
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "admin", Name: "superAdmin", Component: "view/superAdmin/index.vue", Sort: 3, Meta: Meta{Title: "管理员", Icon: "user"}},
		{MenuLevel: 0, Hidden: false, ParentId: 3, Path: "authority", Name: "authority", Component: "view/superAdmin/authority/authority.vue", Sort: 1, Meta: Meta{Title: "角色管理", Icon: "avatar"}},
		{MenuLevel: 0, Hidden: false, ParentId: 3, Path: "menu", Name: "menu", Component: "view/superAdmin/menu/menu.vue", Sort: 2, Meta: Meta{Title: "菜单管理", Icon: "tickets", KeepAlive: true}},
		{MenuLevel: 0, Hidden: false, ParentId: 3, Path: "api", Name: "api", Component: "view/superAdmin/api/api.vue", Sort: 3, Meta: Meta{Title: "api管理", Icon: "platform", KeepAlive: true}},
		{MenuLevel: 0, Hidden: false, ParentId: 3, Path: "user", Name: "user", Component: "view/superAdmin/user/user.vue", Sort: 4, Meta: Meta{Title: "用户管理", Icon: "coordinate"}},
		{MenuLevel: 0, Hidden: false, ParentId: 3, Path: "dictionary", Name: "dictionary", Component: "view/superAdmin/dictionary/sysDictionary.vue", Sort: 5, Meta: Meta{Title: "字典管理", Icon: "notebook"}},
		{MenuLevel: 0, Hidden: false, ParentId: 3, Path: "operation", Name: "operation", Component: "view/superAdmin/operation/sysOperationRecord.vue", Sort: 6, Meta: Meta{Title: "操作历史", Icon: "pie-chart"}},
		{MenuLevel: 0, Hidden: true, ParentId: 0, Path: "person", Name: "person", Component: "view/person/person.vue", Sort: 4, Meta: Meta{Title: "个人信息", Icon: "message"}},
		{MenuLevel: 0, Hidden: true, ParentId: 0, Path: "example", Name: "example", Component: "view/example/index.vue", Sort: 7, Meta: Meta{Title: "示例文件", Icon: "management"}},
		{MenuLevel: 0, Hidden: false, ParentId: 11, Path: "upload", Name: "upload", Component: "view/example/upload/upload.vue", Sort: 5, Meta: Meta{Title: "媒体库（上传下载）", Icon: "upload"}},
		{MenuLevel: 0, Hidden: false, ParentId: 11, Path: "breakpoint", Name: "breakpoint", Component: "view/example/breakpoint/breakpoint.vue", Sort: 6, Meta: Meta{Title: "断点续传", Icon: "upload-filled"}},
		{MenuLevel: 0, Hidden: false, ParentId: 11, Path: "customer", Name: "customer", Component: "view/example/customer/customer.vue", Sort: 7, Meta: Meta{Title: "客户列表（资源示例）", Icon: "avatar"}},
		{MenuLevel: 0, Hidden: true, ParentId: 0, Path: "systemTools", Name: "systemTools", Component: "view/systemTools/index.vue", Sort: 5, Meta: Meta{Title: "系统工具", Icon: "tools"}},
		{MenuLevel: 0, Hidden: false, ParentId: 15, Path: "autoCode", Name: "autoCode", Component: "view/systemTools/autoCode/index.vue", Sort: 1, Meta: Meta{Title: "代码生成器", Icon: "cpu", KeepAlive: true}},
		{MenuLevel: 0, Hidden: false, ParentId: 15, Path: "formCreate", Name: "formCreate", Component: "view/systemTools/formCreate/index.vue", Sort: 2, Meta: Meta{Title: "表单生成器", Icon: "magic-stick", KeepAlive: true}},
		{MenuLevel: 0, Hidden: false, ParentId: 15, Path: "system", Name: "system", Component: "view/systemTools/system/system.vue", Sort: 3, Meta: Meta{Title: "系统配置", Icon: "operation"}},
		{MenuLevel: 0, Hidden: false, ParentId: 15, Path: "autoCodeAdmin", Name: "autoCodeAdmin", Component: "view/systemTools/autoCodeAdmin/index.vue", Sort: 1, Meta: Meta{Title: "自动化代码管理", Icon: "magic-stick"}},
		{MenuLevel: 0, Hidden: true, ParentId: 15, Path: "autoCodeEdit/:id", Name: "autoCodeEdit", Component: "view/systemTools/autoCode/index.vue", Sort: 0, Meta: Meta{Title: "自动化代码-${id}", Icon: "magic-stick"}},
		{MenuLevel: 0, Hidden: false, ParentId: 15, Path: "autoPkg", Name: "autoPkg", Component: "view/systemTools/autoPkg/autoPkg.vue", Sort: 0, Meta: Meta{Title: "自动化package", Icon: "folder"}},
		{MenuLevel: 0, Hidden: true, ParentId: 0, Path: "https://www.curescan.com", Name: "https://www.curescan.com", Component: "/", Sort: 0, Meta: Meta{Title: "官方网站", Icon: "customer-gva"}},
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "state", Name: "state", Component: "view/system/state.vue", Sort: 8, Meta: Meta{Title: "服务器状态", Icon: "cloudy"}},
		{MenuLevel: 0, Hidden: true, ParentId: 0, Path: "plugin", Name: "plugin", Component: "view/routerHolder.vue", Sort: 6, Meta: Meta{Title: "插件系统", Icon: "cherry"}},
		{MenuLevel: 0, Hidden: false, ParentId: 24, Path: "https://plugin.curescan.com/", Name: "https://plugin.curescan.com/", Component: "https://plugin.curescan.com/", Sort: 0, Meta: Meta{Title: "插件市场", Icon: "shop"}},
		{MenuLevel: 0, Hidden: false, ParentId: 24, Path: "installPlugin", Name: "installPlugin", Component: "view/systemTools/installPlugin/index.vue", Sort: 1, Meta: Meta{Title: "插件安装", Icon: "box"}},
		{MenuLevel: 0, Hidden: false, ParentId: 24, Path: "autoPlug", Name: "autoPlug", Component: "view/systemTools/autoPlug/autoPlug.vue", Sort: 2, Meta: Meta{Title: "插件模板", Icon: "folder"}},
		{MenuLevel: 0, Hidden: false, ParentId: 24, Path: "pubPlug", Name: "pubPlug", Component: "view/systemTools/pubPlug/pubPlug.vue", Sort: 3, Meta: Meta{Title: "打包插件", Icon: "files"}},
		{MenuLevel: 0, Hidden: false, ParentId: 24, Path: "plugin-email", Name: "plugin-email", Component: "plugin/email/view/index.vue", Sort: 4, Meta: Meta{Title: "邮件插件", Icon: "message"}},
		{MenuLevel: 0, Hidden: false, ParentId: 15, Path: "exportTemplate", Name: "exportTemplate", Component: "view/systemTools/exportTemplate/exportTemplate.vue", Sort: 10, Meta: Meta{Title: "表格模板", Icon: "reading"}},
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "policy", Name: "policy", Component: "view/routerHolder.vue", Sort: 10, Meta: Meta{Title: "策略管理", Icon: "lollipop"}},
		{MenuLevel: 0, Hidden: false, ParentId: 31, Path: "list", Name: "list", Component: "view/policy/list.vue", Sort: 0, Meta: Meta{Title: "策略列表", Icon: "list"}},
		{MenuLevel: 0, Hidden: true, ParentId: 31, Path: "create", Name: "create", Component: "view/policy/create.vue", Sort: 3, Meta: Meta{Title: "创建策略", Icon: "briefcase"}},
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "task", Name: "task", Component: "view/routerHolder.vue", Sort: 11, Meta: Meta{Title: "扫描任务", Icon: "search"}},
		{MenuLevel: 0, Hidden: false, ParentId: 34, Path: "cronTask", Name: "cronTask", Component: "view/task/cronTask.vue", Sort: 2, Meta: Meta{Title: "定时任务", Icon: "alarm-clock"}},
		{MenuLevel: 0, Hidden: false, ParentId: 34, Path: "index", Name: "index", Component: "view/task/index.vue", Sort: 0, Meta: Meta{Title: "扫描任务", Icon: "bell"}},
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "template", Name: "template", Component: "view/template/index.vue", Sort: 16, Meta: Meta{Title: "模板管理", Icon: "document-copy"}},
		{MenuLevel: 0, Hidden: false, ParentId: 37, Path: "basic", Name: "basic", Component: "view/template/index.vue", Sort: 1, Meta: Meta{Title: "模板列表", Icon: "copy-document"}},
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "assets", Name: "assets", Component: "view/routerHolder.vue", Sort: 12, Meta: Meta{Title: "资产管理", Icon: "guide"}},
		{MenuLevel: 0, Hidden: false, ParentId: 39, Path: "assets/list", Name: "assets/list", Component: "view/assets/list.vue", Sort: 2, Meta: Meta{Title: "资产列表", Icon: "data-analysis"}},
		{MenuLevel: 0, Hidden: false, ParentId: 39, Path: "district", Name: "district", Component: "view/assets/index.vue", Sort: 2, Meta: Meta{Title: "区域管理", Icon: "place"}},
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "cve", Name: "cve", Component: "view/routerHolder.vue", Sort: 17, Meta: Meta{Title: "漏洞管理", Icon: "data-analysis"}},
		{MenuLevel: 0, Hidden: false, ParentId: 42, Path: "data", Name: "data", Component: "view/cve/index.vue", Sort: 1, Meta: Meta{Title: "漏洞数据库", Icon: "expand"}},
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "protocol", Name: "protocol", Component: "view/routerHolder.vue", Sort: 18, Meta: Meta{Title: "协议测试", Icon: "flag"}},
		{MenuLevel: 0, Hidden: false, ParentId: 44, Path: "project", Name: "project", Component: "view/shell/index.vue", Sort: 1, Meta: Meta{Title: "项目管理", Icon: "briefcase"}},
		{MenuLevel: 0, Hidden: false, ParentId: 44, Path: "test", Name: "test", Component: "view/shell/test.vue", Sort: 2, Meta: Meta{Title: "测试列表", Icon: "chicken"}},
		{MenuLevel: 0, Hidden: false, ParentId: 44, Path: "problemReport", Name: "problemReport", Component: "view/shell/problemReport.vue", Sort: 3, Meta: Meta{Title: "缺陷列表", Icon: "umbrella"}},
		{MenuLevel: 0, Hidden: false, ParentId: 44, Path: "datav", Name: "datav", Component: "view/shell/data.vue", Sort: 4, Meta: Meta{Title: "数据分析", Icon: "data-analysis"}},
	}
	if err = db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, SysBaseMenu{}.TableName()+"表数据初始化失败!")
	}
	next = context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}

func (i *initMenu) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("path = ?", "autoPkg").First(&SysBaseMenu{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
