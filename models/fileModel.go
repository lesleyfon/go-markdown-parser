package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type File struct {
	ID           primitive.ObjectID `bson:"_id"`
	User_id      string             `json:"user_id"`
	File_name    string             `json:"file_name"`
	File_content string             `json:"file_content"`
	Created_at   time.Time          `json:"created_at"`
	Updated_at   time.Time          `json:"updated_at"`
}
