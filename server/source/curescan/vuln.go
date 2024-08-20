package curescan

import (
	. "47.103.136.241/goprojects/curescan/server/model/curescan"
	"47.103.136.241/goprojects/curescan/server/service/system"
	"context"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderVuln = system.InitOrderCurescan + 1

type initVuln struct {
}

func init() {
	system.RegisterInit(initOrderVuln, &initVuln{})
}
func (i initVuln) InitializerName() string {
	return Vuln{}.TableName()
}

func (i *initVuln) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&Vuln{},
	)
}

func (i *initVuln) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	m := db.Migrator()
	return m.HasTable(&Vuln{})
}

func (i *initVuln) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	entities := []Vuln{
		{
			TemplateID:  "3com-nj2000-default-login",
			Name:        "3COM NJ2000 - default login",
			Author:      "zy",
			Severity:    "high",
			Description: "3COM NJ2000存在默认登录漏洞。发现了默认管理员登录密码为“password”。攻击者可以获取用户账户访问权限，访问敏感信息，修改数据和/或执行未授权的操作。",
			Reference: pq.StringArray{
				"https://www.manualslib.com/manual/204158/3com-Intellijack-Nj2000.html?page=12",
			},
			Classification: JSONB{
				"cvss-metrics": "CVSS:3.0/AV:N/AC:L/PR:N/UI:N/S:C/C:L/I:L/A:L",
				"cvss-score":   8.3,
				"cwe-id":       "CWE-522",
			},
			Remediation: "为确保您的系统不受该漏洞的影响，请及时下载并安装厂商发布的修复补丁。具体操作步骤请参考厂商提供的指南",
		},
		{
			TemplateID:  "3Com-wireless-default-login",
			Name:        "3Com Wireless 8760 - default login",
			Author:      "zy",
			Severity:    "high",
			Description: "3COM Wireless 8760 Dual Radio存在默认登录漏洞。发现了默认管理员登录密码“password”。",
			Reference: pq.StringArray{
				"https://www.speedguide.net/routers/3com-wl-546-3com-wireless-8760-dual-radio-11abg-1256",
			},
			Remediation: "厂商已发布针对该漏洞的修复补丁，请尽早安装更新以保护您的设备。更多信息请参考厂商公告",
		},
		{
			TemplateID:  "3cx-management-console",
			Name:        "3CX管理控制台 - 本地文件包含",
			Author:      "zy",
			Severity:    "high",
			Description: "3CX管理控制台存在本地文件包含漏洞。",
			Reference: pq.StringArray{
				"https://medium.com/@frycos/pwning-3cx-phone-management-backends-from-the-internet-d0096339dd88",
			},
			Classification: JSONB{
				"cvss-metrics": "CVSS:3.0/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:N/A:N",
				"cvss-score":   7.5,
				"cwe-id":       "CWE-22",
			},
			Remediation: "目前厂商已发布升级补丁以修复漏洞，详情请关注厂商主页",
		},
		{
			TemplateID:  "3d-print-lite-xss",
			Name:        "3D Print Lite < 1.9.1.6 - 反射型跨站脚本攻击",
			Author:      "zy",
			Severity:    "medium",
			Description: "该插件在输出用户输入到属性中时未对其进行清理和转义，导致了反射型跨站脚本攻击问题。",
			Reference: pq.StringArray{
				"https://wpscan.com/vulnerability/5909e225-5756-472e-a2fc-3ac52c7fb909",
				"https://www.acunetix.com/vulnerabilities/web/wordpress-plugin-3dprint-lite-cross-site-scripting-1-9-1-5/",
			},
			Remediation: "更新到插件版本1.9.1.6或最新版本",
		},
		{
			TemplateID:  "3dprint-arbitrary-file-upload",
			Name:        "WordPress 3DPrint Lite <1.9.1.5 - Arbitrary File Upload",
			Author:      "zy",
			Severity:    "high",
			Description: "WordPress 3DPrint Lite插件在1.9.1.5版本之前存在任意文件上传漏洞。该插件的p3dlite_handle_upload AJAX操作没有任何授权，并且不检查上传的文件。攻击者可以将任意文件上传到服务器，从而可能会获取敏感信息，修改数据和/或执行未经授权的操作。",
			Reference: pq.StringArray{
				"https://wpscan.com/vulnerability/c46ecd0d-a132-4ad6-b936-8acde3a09282",
				"https://www.exploit-db.com/exploits/50321",
			},
			Classification: JSONB{
				"cvss-metrics": "CVSS:3.1/AV:N/AC:L/PR:L/UI:N/S:U/C:H/I:H/A:H",
				"cvss-score":   8.8,
				"cwe-id":       "CWE-434",
			},
			Remediation: "升级到1.9.1.5或更高版本。",
		},
	}
	if err = db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, Vuln{}.TableName()+"表初始化失败！")
	}
	next = context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}

func (i *initVuln) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("author = ?", "zy").First(&Vuln{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
