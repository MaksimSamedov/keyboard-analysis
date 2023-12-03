package models

import "gorm.io/gorm"

type KeyboardFlow struct {
	gorm.Model `json:"-"`
	Flow       []*KeyboardEvent `json:"flow"`
	Password   Password         `json:"-"`
	PasswordID uint             `gorm:"index:password_id" json:"-"`
}

type KeyboardEvent struct {
	gorm.Model     `json:"-"`
	KeyboardFlowID uint         `json:"-" gorm:"index:flow_id"`
	KeyboardFlow   KeyboardFlow `json:"-"`
	Key            string       `json:"key"`
	Up             bool         `json:"up"`
	Time           uint         `json:"time"` // В миллисекундах
}

func (flow *KeyboardFlow) RemoveInvalidEvents() {
	var res []*KeyboardEvent
	for _, ev := range flow.Flow {
		if len(ev.Key) == 1 {
			res = append(res, ev)
		}
	}
	flow.Flow = res
}

func (flow *KeyboardFlow) TruncateTime() {
	l := len(flow.Flow)
	if l == 0 {
		return
	}

	start := flow.Flow[0].Time
	for i := 0; i < l; i++ {
		flow.Flow[i].Time -= start
	}
}
