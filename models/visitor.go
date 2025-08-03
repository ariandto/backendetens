//models/visitor.go

package models

import "time"

type Visitor struct {
	ID          uint      `gorm:"primary key"`
	IP          string    `gorm:"size:45;not null"`
	VisitedDate time.Time `gorm:"type:date;not null"`
	CreatedAt   time.Time
}
