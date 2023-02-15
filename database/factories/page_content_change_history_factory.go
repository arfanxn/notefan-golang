package factories

import (
	"time"

	"github.com/notefan-golang/helper"
	"github.com/notefan-golang/models/entities"
)

func FakePageContentChangeHistory() entities.PageContentChangeHistory {
	return entities.PageContentChangeHistory{
		//BeforePageContentId: , // will be filled later
		//AfterPageContentId: , // will be filled later
		//UserId: , // will be filled later
		CreatedAt: time.Now(),
		UpdatedAt: helper.DBRandNullOrTime(time.Now().AddDate(0, 0, 1)),
	}
}
