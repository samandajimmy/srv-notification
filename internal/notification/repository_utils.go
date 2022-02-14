package notification

import (
	"github.com/nbs-go/nsql"
	"github.com/nbs-go/nsql/option"
	"github.com/nbs-go/nsql/pq/query"
	"github.com/nbs-go/nsql/schema"
)

//func newStatusFilter(s *schema.Schema) nsql.FilterParser {
//	return func(qv string) (nsql.WhereWriter, []interface{}) {
//		// Get args
//		args := parse.IntArgs(qv)
//		w := query.In(query.Column("statusId", option.Schema(s)), len(args))
//		return w, args
//	}
//}

func newEqualFilter(s *schema.Schema, col string) nsql.FilterParser {
	return func(qv string) (nsql.WhereWriter, []interface{}) {
		w := query.Equal(query.Column(col, option.Schema(s)))
		return w, []interface{}{qv}
	}
}

func newGreaterThanEqualFilter(s *schema.Schema, col string) nsql.FilterParser {
	return func(qv string) (nsql.WhereWriter, []interface{}) {
		w := query.GreaterThanEqual(query.Column(col, option.Schema(s)))
		return w, []interface{}{qv}
	}
}

func newLessThanEqualFilter(s *schema.Schema, col string) nsql.FilterParser {
	return func(qv string) (nsql.WhereWriter, []interface{}) {
		w := query.GreaterThanEqual(query.Column(col, option.Schema(s)))
		return w, []interface{}{qv}
	}
}
