package factories

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/notefan-golang/helpers/sliceh"
	"github.com/notefan-golang/models/entities"
)

func FakeToken() entities.Token {
	nullTime := sql.NullTime{Time: time.Time{}, Valid: false}

	token := entities.Token{
		Id: uuid.New(),
		// TokenableType: // will be filled in later ,
		// TokenableId: // will be filled in later ,
		// Type: // will be filled in later ,
		// Body: strconv.Itoa(numberh.Random(100000, 999999)), // will be filled in later ,
		UsedAt: sliceh.Random([]sql.NullTime{
			sql.NullTime{Time: time.Now().Add(-(time.Hour / 2)), Valid: true},
			nullTime,
		}),
		ExpiredAt: sliceh.Random([]sql.NullTime{
			sql.NullTime{Time: time.Now().Add(time.Hour / 2), Valid: true},
			nullTime,
		}),
		CreatedAt: time.Now(),
		UpdatedAt: sliceh.Random([]sql.NullTime{
			sql.NullTime{Time: time.Now().Add(time.Hour), Valid: true},
			nullTime,
		}),
	}

	return token
}
