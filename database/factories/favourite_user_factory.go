package factories

import (
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"time"
)

func NewFavouriteUser() entities.FavouriteUser {
	return entities.FavouriteUser{
		//FavouriteableType: , will be filled in later
		//FavouriteableId: , will be filled in later
		//UserId: , will be filled in later
		Order:     1,
		CreatedAt: time.Now(),
		UpdatedAt: helper.DBRandNullOrTime(time.Now().AddDate(0, 0, 1)),
	}
}
