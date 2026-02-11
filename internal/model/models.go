package model

import (
	"time"

	"gorm.io/gorm"
)

type ImageRegion struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Date          string         `gorm:"uniqueIndex:idx_date_mkt;index:idx_mkt_date,priority:2;type:varchar(10)" json:"date"` // YYYY-MM-DD
	Mkt           string         `gorm:"uniqueIndex:idx_date_mkt;index:idx_mkt_date,priority:1;type:varchar(10)" json:"mkt"`  // zh-CN, en-US etc.
	HSH           string         `gorm:"type:varchar(64)" json:"hsh"`
	URLBase       string         `json:"urlbase"`
	ImageName     string         `gorm:"index" json:"image_name"`
	Title         string         `json:"title"`
	Copyright     string         `json:"copyright"`
	CopyrightLink string         `json:"copyrightlink"`
	Quiz          string         `json:"quiz"`
	StartDate     string         `json:"startdate"`
	FullStartDate string         `json:"fullstartdate"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Variants      []ImageVariant `gorm:"foreignKey:ImageName;references:ImageName" json:"variants"`
}

type ImageVariant struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ImageName  string    `gorm:"uniqueIndex:idx_name_variant_format;type:varchar(100)" json:"image_name"`
	Variant    string    `gorm:"uniqueIndex:idx_name_variant_format;type:varchar(20)" json:"variant"` // UHD, 1920x1080, etc.
	Format     string    `gorm:"uniqueIndex:idx_name_variant_format;type:varchar(10)" json:"format"`  // jpg, webp
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

type ApiStat struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Date      string    `gorm:"uniqueIndex:idx_date_endpoint_mkt;type:varchar(10)" json:"date"` // YYYY-MM-DD
	Endpoint  string    `gorm:"uniqueIndex:idx_date_endpoint_mkt;type:varchar(100)" json:"endpoint"`
	Mkt       string    `gorm:"uniqueIndex:idx_date_endpoint_mkt;type:varchar(20)" json:"mkt"`
	Count     int64     `gorm:"default:0" json:"count"`
	UpdatedAt time.Time `json:"updated_at"`
}
