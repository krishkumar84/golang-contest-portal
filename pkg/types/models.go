package types

import (
    "time"
)

type User struct {
    Name      string    `bson:"name" json:"name" validate:"required"`
    Email     string    `bson:"email" json:"email" validate:"required"`
    Password  string    `bson:"password" json:"password" validate:"required"`
	StudentId string    `bson:"studentId" json:"studentId" validate:"required"`
    CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

type Contest struct {
    ID        string    `bson:"_id,omitempty" json:"contest_id" validate:"required"`
    Title     string    `bson:"title" json:"title" validate:"required"`
    StartTime time.Time `bson:"start_time" json:"start_time" validate:"required"`
    EndTime   time.Time `bson:"end_time" json:"end_time" validate:"required"`
    Questions []string  `bson:"questions" json:"questions" validate:"required"`
}
