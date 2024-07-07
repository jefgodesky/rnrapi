package initializers

import (
	"github.com/jefgodesky/rnrapi/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	dsn := "host=" + os.Getenv("POSTGRES_HOST") +
		" user=" + os.Getenv("POSTGRES_USER") +
		" password=" + os.Getenv("POSTGRES_PASSWORD") +
		" dbname=" + os.Getenv("POSTGRES_DB") +
		" port=" + os.Getenv("POSTGRES_PORT") +
		" sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database")
	}
}

func MigrateDB() error {
	err := DB.AutoMigrate(
		&models.User{},
		&models.World{},
		&models.Campaign{},
	)

	if err != nil {
		return err
	}

	err = DB.Exec("CREATE UNIQUE INDEX idx_campaign_world_slug ON campaigns (world_id, slug)").Error
	if err != nil {
		return err
	}

	return nil
}
