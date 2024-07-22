package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/redis/go-redis/v9"
	"github.com/streadway/amqp"
	"golang-pet-api/config"
	"gorm.io/gorm"
)

var (
	Trans   ut.Translator
	Config  config.Config
	Db      *gorm.DB
	RedisDb *redis.Client
	Conn    *amqp.Connection
)
