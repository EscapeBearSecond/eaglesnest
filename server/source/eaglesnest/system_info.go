package eaglesnest

import (
	"bufio"
	"context"
	"github.com/EscapeBearSecond/eaglesnest/server/model/eaglesnest"
	"github.com/EscapeBearSecond/eaglesnest/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"os"
	"strings"
	"time"
)

const initOrderSystemInfo = system.InitOrderSystemInfo + 1

type initSystemInfo struct {
}

func init() {
	system.RegisterInit(initOrderSystemInfo, &initSystemInfo{})
}
func (i initSystemInfo) InitializerName() string {
	return eaglesnest.SystemInfo{}.TableName()
}

func (i *initSystemInfo) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&eaglesnest.SystemInfo{},
	)
}

func (i *initSystemInfo) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	m := db.Migrator()
	return m.HasTable(&eaglesnest.SystemInfo{})
}
func readVersionInfo() (string, error) {
	filePath := "../version.ini"
	file, err := os.Open(filePath)
	if err != nil {
		return "0.0.0", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var version string
	for scanner.Scan() {
		line := scanner.Text()
		version = strings.TrimSpace(line) // 去除首尾空格
		break
	}
	return version, scanner.Err()
}

func (i *initSystemInfo) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	// systemVersion, err := readVersionInfo()
	// if err != nil {
	// 	return ctx, err
	// }
	systemVersion := "0.1.1"
	vulnVersion := "0.0.0"
	lastUpdateDate := time.Now().Format("2006-01-02 15:04:05")
	// watcher, err := license.Watch("./license.json")
	// if err != nil {
	// 	return ctx, err
	// }
	// defer watcher.Stop()
	//
	// licenseExpiration := license.L().ExpiresAt
	licenseExpiration := ""
	entities := []eaglesnest.SystemInfo{
		{
			SystemVersion:     systemVersion,
			VulnVersion:       vulnVersion,
			LastUpdateDate:    lastUpdateDate,
			LicenseExpiration: licenseExpiration,
		},
	}
	if err = db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, eaglesnest.SystemInfo{}.TableName()+"表初始化失败！")
	}
	next = context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}

func (i *initSystemInfo) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("id = ?", 1).First(&eaglesnest.SystemInfo{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
