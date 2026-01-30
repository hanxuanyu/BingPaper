package model

import (
	"time"

	"gorm.io/gorm"
)

type Image struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Date          string         `gorm:"uniqueIndex:idx_date_mkt;type:varchar(10)" json:"date"` // YYYY-MM-DD
	Mkt           string         `gorm:"uniqueIndex:idx_date_mkt;type:varchar(10)" json:"mkt"`  // zh-CN, en-US etc.
	Title         string         `json:"title"`
	Copyright     string         `json:"copyright"`
	CopyrightLink string         `json:"copyrightlink"`
	URLBase       string         `json:"urlbase"`
	Quiz          string         `json:"quiz"`
	StartDate     string         `json:"startdate"`
	FullStartDate string         `json:"fullstartdate"`
	HSH           string         `json:"hsh"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Variants      []ImageVariant `gorm:"foreignKey:ImageID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"variants"`
}

type ImageVariant struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ImageID    uint      `gorm:"index;uniqueIndex:idx_image_variant_format" json:"image_id"`
	Variant    string    `gorm:"uniqueIndex:idx_image_variant_format;type:varchar(20)" json:"variant"` // UHD, 1920x1080, etc.
	Format     string    `gorm:"uniqueIndex:idx_image_variant_format;type:varchar(10)" json:"format"`  // jpg, webp
	StorageKey string    `json:"storage_key"`
	PublicURL  string    `json:"public_url"`
	Size       int64     `json:"size"`
	CreatedAt  time.Time `json:"created_at"`
}

type Token struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Token     string    `gorm:"uniqueIndex;type:varchar(64)" json:"token"`
	Name      string    `json:"name"`
	ExpiresAt time.Time `json:"expires_at"`
	Disabled  bool      `gorm:"default:false" json:"disabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
