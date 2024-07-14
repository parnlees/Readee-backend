package common

type Config struct {
	Environment *uint8  `yaml:"environment" validate:"gte=1,lte=2"`
	Address     *string `yaml:"address" validate:"required"`
	PostgresDsn *string `yaml:"postgres_dsn" validate:"required"`
}
