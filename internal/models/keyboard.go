package models

import "gorm.io/gorm"

type KeyboardFlow struct {
	gorm.Model `json:"-"`
	Flow       []KeyboardEvent `json:"flow"`
	Phrase     string          `json:"phrase" gorm:"column:phrase;type:char(127);not null;<-:create"`
}

type KeyboardEvent struct {
	gorm.Model     `json:"-"`
	KeyboardFlowID uint         `json:"-" gorm:"index:flow_id"`
	KeyboardFlow   KeyboardFlow `json:"-"`
	Key            string       `json:"key"`
	Up             bool         `json:"up"`
	Time           uint         `json:"time"` // В миллисекундах
}
