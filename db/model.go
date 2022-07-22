package db

import (
	"time"
)

type VTBInfo struct {
	ID               int64 `gorm:"primaryKey;autoIncrement"`
	Archive          bool
	BilibiliLiveRoom int64
	BilibiliUid      int64  `gorm:"index:,unique"`
	Face             string `gorm:"type:varchar(256)"`
	Hidden           bool
	IsLive           bool
	LastDanmu        int64
	LastLive         int64
	Name             string `gorm:"type:varchar(64)"`
	TotalClips       int
	TotalDanmu       int
}

type Comment struct {
	ID       int64     `gorm:"primaryKey;autoIncrement"`
	VTBUid   int64     `gorm:"not null"`
	ClipID   string    `gorm:"type:varchar(64)"`
	UserID   int64     `gorm:"not null;index:comment_user_time"`
	Username string    `gorm:"type:varchar(64)"`
	Time     time.Time `gorm:"not null;index:comment_user_time"`
	Text     string    `gorm:"type:varchar(256)"`
}

type SuperChat struct {
	ID             int64     `gorm:"primaryKey;autoIncrement"`
	VTBUid         int64     `gorm:"not null"`
	ClipID         string    `gorm:"type:varchar(64)"`
	UserID         int64     `gorm:"not null;index:super_chat_user_time"`
	Username       string    `gorm:"type:varchar(64);not null"`
	Time           time.Time `gorm:"not null;index:super_chat_user_time"`
	Text           string    `gorm:"type:varchar(256)"`
	SuperChatPrice int
}

type Gift struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	VTBUid    int64     `gorm:"not null"`
	ClipID    string    `gorm:"type:varchar(64)"`
	UserID    int64     `gorm:"not null;index:gift_user_time"`
	Username  string    `gorm:"type:varchar(64);not null"`
	Time      time.Time `gorm:"not null;index:gift_user_time"`
	GiftName  string    `gorm:"type:varchar(64)"`
	GiftNum   int
	GiftPrice float64
}

type Clip struct {
	ClipID     string    `gorm:"type:varchar(64);primaryKey"`
	VTBUid     int64     `gorm:"not null;index:clip_uid_time"`
	Cover      string    `gorm:"type:varchar(256)"`
	StartTime  time.Time `gorm:"not null"`
	EndTime    time.Time `gorm:"not null;index:clip_uid_time"`
	Title      string    `gorm:"type:varchar(64)"`
	TotalDanmu int
	Views      int
}
