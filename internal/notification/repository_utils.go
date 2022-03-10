package notification

import (
	"github.com/nbs-go/nsql"
	"github.com/nbs-go/nsql/option"
	"github.com/nbs-go/nsql/pq/query"
	"github.com/nbs-go/nsql/schema"
	qs "repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp/querystring"
)

func newEqualFilter(s *schema.Schema, col string) nsql.FilterParser {
	return func(qv string) (nsql.WhereWriter, []interface{}) {
		w := query.Equal(query.Column(col, option.Schema(s)))
		return w, []interface{}{qv}
	}
}

func newGreaterThanEqualFilter(s *schema.Schema, col string) nsql.FilterParser {
	return func(qv string) (nsql.WhereWriter, []interface{}) {
		// Parse time
		t, ok := qs.ParseTime(qv)
		if !ok {
			return nil, nil
		}

		// Create schema
		w := query.GreaterThanEqual(query.Column(col, option.Schema(s)))
		return w, []interface{}{t.UTC()}
	}
}

func newLessThanEqualFilter(s *schema.Schema, col string) nsql.FilterParser {
	return func(qv string) (nsql.WhereWriter, []interface{}) {
		// Parse time
		t, ok := qs.ParseTime(qv)
		if !ok {
			return nil, nil
		}

		w := query.LessThanEqual(query.Column(col, option.Schema(s)))
		return w, []interface{}{t.UTC()}
	}
}
