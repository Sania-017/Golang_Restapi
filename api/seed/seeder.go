package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/victorsteven/fullstack/api/models"
)

var users = []models.User{
	models.User{
		Nickname: "Steven victor",
		Email:    "steven@gmail.com",
		Password: "password",
	},
	models.User{
		Nickname: "Martin Luther",
		Email:    "luther@gmail.com",
		Password: "password",
	},
}

var posts = []models.Post{
	models.Post{
		Title:   "Title 1",
		Content: "Hello world 1",
	},
	models.Post{
		Title:   "Title 2",
		Content: "Hello world 2",
	},
}
var address = []models.Address{
	models.Address{

		Address: "hostel city line 10",
		Phone:   "456712345",
	},
	models.Address{

		Address: "hostel city line 14",
		Phone:   "556789012",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Post{}, &models.User{}, &models.Address{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}

	err = db.Debug().AutoMigrate(&models.User{}, &models.Post{}, &models.Address{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	/*
		err = db.Debug().Model(&models.Post{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
		if err != nil {
			log.Fatalf("attaching foreign key error: %v", err)
		}
	*/

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		posts[i].AuthorID = users[i].ID

		err = db.Debug().Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
		address[i].AuthorID = users[i].ID
		err = db.Debug().Model(&models.Address{}).Create(&address[i]).Error
		if err != nil {
			log.Fatalf("cannot seed address table: %v", err)
		}
	}
}
