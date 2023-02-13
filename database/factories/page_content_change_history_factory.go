package factories

import (
	"notefan-golang/helper"
	"notefan-golang/models/entities"
	"time"
)

func NewPageContentChangeHistory() entities.PageContentChangeHistory {
	return entities.PageContentChangeHistory{
		//BeforePageContentId: , // will be filled later
		//AfterPageContentId: , // will be filled later
		//UserId: , // will be filled later
		CreatedAt: time.Now(),
		UpdatedAt: helper.DBRandNullOrTime(time.Now().AddDate(0, 0, 1)),
	}
}
