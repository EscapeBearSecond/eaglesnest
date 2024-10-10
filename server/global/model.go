package global

import (
	"time"

	"gorm.io/gorm"
)

type GvaModel struct {
	ID        uint           `gorm:"primaryKey" json:"ID"` // 主键ID
	CreatedAt time.Time      // 创建时间
	UpdatedAt time.Time      // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 删除时间
}

type CsModel struct {
	CreatedBy uint `gorm:"column:created_by;type:int8;comment:创建者" json:"createdBy"`
	UpdatedBy uint `gorm:"column:updated_by;type:int8;comment:更新者" json:"updatedBy"`
}
