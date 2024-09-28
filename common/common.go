package cc

import (
	"Readee-Backend/type/common"

	"github.com/gofiber/fiber"
	"gorm.io/gorm"
)

var Config *common.Config
var DB *gorm.DB
var App *fiber.App
