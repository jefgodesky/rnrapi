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
		&models.Species{},
		&models.Society{},
	)

	if err != nil {
		return err
	}

	createUniqueIndex("idx_campaign_world_slug", "campaigns", "world_id, slug")
	createUniqueIndex("idx_species_world_slug", "species", "world_id, slug")
	createUniqueIndex("idx_society_world_slug", "societies", "world_id, slug")

	return nil
}

func createUniqueIndex(indexName, tableName, columns string) error {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = ?)"
	DB.Raw(query, indexName).Scan(&exists)
	if !exists {
		createIndexQuery := "CREATE UNIQUE INDEX " + indexName + " ON " + tableName + " (" + columns + ")"
		err := DB.Exec(createIndexQuery).Error
		if err != nil {
			return err
		}
	}
	return nil
}
