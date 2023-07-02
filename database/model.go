package database

import (
	. "github.com/monkey-panel/control-panel-utils/types"

	"gorm.io/gorm"
)

// set up models
func setupModel(db *DB) {
	db.AutoMigrate(&DBInstance{})
}

// database base model struct
type BaseModel struct {
	ID        ID   `gorm:"primarykey" json:"id"`
	CreatedAt Time `gorm:"<-:create;autoCreateTime" json:"create_at"`
}

// set ID
func (i *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	i.ID = GlobalIDMake.Generate()
	return
}

// database instance struct
type DBInstance struct {
	BaseModel
	Name             string `gorm:"not null"`
	Description      string
	AdminDescription string
	AutoStart        bool
	Mark             InstanceMark
	LastAt           Time
	EndAt            Time
}

func (DBInstance) TableName() string { return "instance" }
