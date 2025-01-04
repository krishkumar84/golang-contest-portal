package types

import (
	"time"
)

type ContestBasicInfo struct {
    ID          string    `bson:"_id" json:"contest_id"`
    Title       string    `bson:"title" json:"title"`
    StartTime   time.Time `bson:"start_time" json:"start_time"`
    EndTime     time.Time `bson:"end_time" json:"end_time"`
    Description string    `bson:"description" json:"description"`
}

