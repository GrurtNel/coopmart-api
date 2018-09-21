package mongodb

import (
	"time"

	"github.com/trunglen/g/x/math"
)

type IModel interface {
	BeforeCreate(prefix string, length int)
	BeforeUpdate()
	BeforeDelete()
}
type BaseModel struct {
	ID        string `json:"id" bson:"_id"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
	UpdatedAt int64  `json:"updated_at" bson:"updated_at"`
}

func (b *BaseModel) BeforeCreate(prefix string, length int) {
	b.ID = math.RandString(prefix, length)
	b.CreatedAt = time.Now().UnixNano() / int64(time.Millisecond)
	b.UpdatedAt = time.Now().UnixNano() / int64(time.Millisecond)
}

func (b *BaseModel) BeforeUpdate() {
	b.UpdatedAt = time.Now().UnixNano() / int64(time.Millisecond)
}

func (b *BaseModel) BeforeDelete() {
	b.UpdatedAt = 0
}
