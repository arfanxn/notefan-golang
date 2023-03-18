package query_rules

import (
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	query_order_bys "github.com/notefan-golang/enums/query/order_bys"
	"github.com/notefan-golang/exceptions"
)

// OrderBys validates the request order bys
func OrderBys(columnNames []string) validation.RuleFunc {
	return func(value interface{}) error {
		orderBys, ok := value.([]string)
		if !ok {
			return exceptions.ValidationInvalidArgumentValue
		}
		regex, err := regexp.Compile("^(" + strings.Join(columnNames, "|") + ")=(?i)(" +
			strings.Join(query_order_bys.All(), "|") + ")$")
		if err != nil {
			return err
		}
		for _, orderBy := range orderBys {
			isMatch := regex.MatchString(orderBy)
			if isMatch == false {
				return exceptions.ValidationUnkownOrderBysField
			}
		}
		return nil
	}
}
