package seeds

import (
	"moonshine/models"
	"moonshine/modules/database"
)

func Load() {
	// If there are no users in the database
	if database.Connection().First(&models.User{}).RecordNotFound() {
		SeedUsers()
	}
}
