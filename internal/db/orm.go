package db

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
)

var Orm *gorm.DB // Подключение к БД
var RedisCl *redis.Client

// Resp - cтруктура используемая для стандартизации ответов АПИ в json
type Resp struct {
	Body Response
}

type Response struct {
	// Status - Статус ответа (Примеры: [200, 400, 404, 500])
	Status int `json:"Status"`

	// Msg - cообщение понятное пользователю.
	Msg string `json:"Message,omitempty"`

	// Data - данные ответа.
	Data interface{} `json:"Data,omitempty"`
}
