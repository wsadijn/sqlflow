package tidb

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pingcap/parser"
	_ "github.com/pingcap/tidb/types/parser_driver" // As required by https://github.com/pingcap/parser/blob/master/parser_example_test.go#L19
)

// Parse calls TiDB's parser to parse a statement sql.  It returns
// <-1,nil> if TiDB parser accepts the statement, or <pos,nil> if TiDB
// doesn't accept but returns a `near "..."` in the error message, or
// <-1,err> if the error messages doens't contain near.
func Parse(sql string) (idx int, err error) {
	p := parser.New()
	_, _, err = p.Parse(sql, "", "")

	if err != nil {
		re := regexp.MustCompile(`.* near "([^"]+)".*`)
		matched := re.FindAllStringSubmatch(err.Error(), -1)

		if len(matched) != 1 || len(matched[0]) != 2 {
			return -1, fmt.Errorf("Cannot match near in %q", err)
		}
		return strings.Index(sql, matched[0][1]), nil
	}
	return -1, nil
}
