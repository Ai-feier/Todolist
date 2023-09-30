package model

type Task struct {
	TaskID int64 `gorm:"primarykey"`  // id
	UserID int64 `gorm:"index"`  // userid
	Status int `gorm:"default:0"` // 默认状态
	Title string 
	Context string `gorm:"type:longtext"`
	StartTime int64 
	EndTime int64
}

func (*Task) TableName() string {
	return "task"
}
