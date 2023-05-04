package util

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	model "github.com/windbnb/reservation-service/model"
)

func ConnectToDatabase() *gorm.DB {
	connectionString := "host=localhost user=postgres dbname=ReservationServiceDB sslmode=disable password=root port=5432"

	dialect := "postgres"

	db, err := gorm.Open(dialect, connectionString)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connection to DB successfull.")
	}

	db.DropTable("reservation_requests")
	db.AutoMigrate(&model.ReservationRequest{})

	return db
}
