package schema

import (
	"fmt"
	"ovaphlow/cratecyclone/utility"
	"strings"
)

type QueryBuilder struct {
	Columns    *string
	Schema     *string
	Table      *string
	Conditions []string
	Params     []string
	Order      *string
	Limit      int
	Offset     int64
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

func (qb *QueryBuilder) Select(columns *string) *QueryBuilder {
	qb.Columns = columns
	return qb
}

func (qb *QueryBuilder) From(schema *string, table *string) *QueryBuilder {
	qb.Schema = schema
	qb.Table = table
	return qb
}

func (qb *QueryBuilder) Equal(equal []string) *QueryBuilder {
	if len(equal)%2 != 0 {
		utility.Slogger.Warn("equal 参数错误")
		return qb
	}
	for i := 0; i < len(equal); i += 2 {
		qb.Conditions = append(qb.Conditions, fmt.Sprintf("%s = ?", equal[i]))
		qb.Params = append(qb.Params, equal[i+1])
	}
	return qb
}

func (qb *QueryBuilder) NotEqual(notEqual []string) *QueryBuilder {
	if len(notEqual)%2 != 0 {
		utility.Slogger.Warn("notEqual 参数错误")
		return qb
	}
	for i := 0; i < len(notEqual); i += 2 {
		qb.Conditions = append(qb.Conditions, fmt.Sprintf("%s != ?", notEqual[i]))
		qb.Params = append(qb.Params, notEqual[i+1])
	}
	return qb
}

func (qb *QueryBuilder) Like(like []string) *QueryBuilder {
	if len(like)%2 != 0 {
		utility.Slogger.Warn("like 参数错误")
		return qb
	}
	for i := 0; i < len(like); i += 2 {
		qb.Conditions = append(
			qb.Conditions,
			fmt.Sprintf("position(? in %s) > 0", like[i]),
		)
		qb.Params = append(qb.Params, like[i+1])
	}
	return qb
}

func (qb *QueryBuilder) Greater(greater []string) *QueryBuilder {
	if len(greater)%2 != 0 {
		utility.Slogger.Warn("greater 参数错误")
		return qb
	}
	for i := 0; i < len(greater); i += 2 {
		qb.Conditions = append(qb.Conditions, fmt.Sprintf("%s >= ?", greater[i]))
		qb.Params = append(qb.Params, greater[i+1])
	}
	return qb
}

func (qb *QueryBuilder) Lesser(lesser []string) *QueryBuilder {
	if len(lesser)%2 != 0 {
		utility.Slogger.Warn("lesser 参数错误")
		return qb
	}
	for i := 0; i < len(lesser); i += 2 {
		qb.Conditions = append(qb.Conditions, fmt.Sprintf("%s <= ?", lesser[i]))
		qb.Params = append(qb.Params, lesser[i+1])
	}
	return qb
}

func (qb *QueryBuilder) In(in []string) *QueryBuilder {
	if len(in) == 0 {
		return qb
	}
	if len(in) < 2 {
		utility.Slogger.Warn("in 参数错误")
		return qb
	}
	c := make([]string, len(in)-1)
	for i := range c {
		c[i] = "?"
	}
	qb.Conditions = append(
		qb.Conditions,
		fmt.Sprintf("%s in (%s)", in[0], strings.Join(c, ", ")),
	)
	qb.Params = append(qb.Params, in[1:]...)
	return qb
}

func (qb *QueryBuilder) NotIn(notIn []string) *QueryBuilder {
	if len(notIn) == 0 {
		return qb
	}
	if len(notIn) < 2 {
		utility.Slogger.Warn("notIn 参数错误")
		return qb
	}
	c := make([]string, len(notIn)-1)
	for i := range c {
		c[i] = "?"
	}
	qb.Conditions = append(
		qb.Conditions,
		fmt.Sprintf("%s not in (%s)", notIn[0], strings.Join(c, ", ")),
	)
	qb.Params = append(qb.Params, notIn[1:]...)
	return qb
}

func (qb *QueryBuilder) ObjectContain(objectContain []string) *QueryBuilder {
	if len(objectContain)%3 != 0 {
		utility.Slogger.Warn("objectContain 参数错误")
		return qb
	}
	for i := 0; i < len(objectContain); i += 3 {
		qb.Conditions = append(
			qb.Conditions,
			fmt.Sprintf(
				"%s @> jsonb_build_object('%s', '%s')",
				objectContain[i],
				objectContain[i+1],
				objectContain[i+2],
			),
		)
	}
	return qb
}

func (qb *QueryBuilder) ArrayContain(arrayContain []string) *QueryBuilder {
	if len(arrayContain)%2 != 0 {
		utility.Slogger.Warn("arrayContain 参数错误")
		return qb
	}
	for i := 0; i < len(arrayContain); i += 2 {
		qb.Conditions = append(
			qb.Conditions,
			fmt.Sprintf("%s @> jsonb_build_array('%s')", arrayContain[i], arrayContain[i+1]),
		)
	}
	return qb
}

func (qb *QueryBuilder) ObjectLike(objectLike []string) *QueryBuilder {
	if len(objectLike)%3 != 0 {
		utility.Slogger.Warn("objectLike 参数错误")
		return qb
	}
	for i := 0; i < len(objectLike); i += 3 {
		qb.Conditions = append(
			qb.Conditions,
			fmt.Sprintf("position(? in %s->>'%s') > 0", objectLike[i], objectLike[i+1]),
		)
		qb.Params = append(qb.Params, objectLike[i+2])
	}
	return qb
}

func (qb *QueryBuilder) OrderBy(orderBy *string) *QueryBuilder {
	qb.Order = orderBy
	return qb
}

func (qb *QueryBuilder) Take(limit int) *QueryBuilder {
	qb.Limit = limit
	return qb
}

func (qb *QueryBuilder) Skip(offset int64) *QueryBuilder {
	qb.Offset = offset
	return qb
}

func (qb *QueryBuilder) Build() (string, []string) {
	q := fmt.Sprintf("select %s from %s.%s", *qb.Columns, *qb.Schema, *qb.Table)
	q = fmt.Sprintf("%s where not (state ? 'deleted_at')", q)
	if len(qb.Conditions) > 0 {
		q = fmt.Sprintf("%s and", q)
		q_ := strings.Join(qb.Conditions, " and ")
		for i := 0; i < len(qb.Params); i++ {
			q_ = strings.Replace(q_, "?", fmt.Sprintf("$%d", i+1), 1)
		}
		q = fmt.Sprintf("%s %s", q, q_)
	}
	if qb.Order == nil || *qb.Order == "" {
		*qb.Order = "id desc"
	}
	q = fmt.Sprintf("%s order by %s", q, *qb.Order)
	q = fmt.Sprintf("%s limit %d offset %d", q, qb.Limit, qb.Offset)
	return q, qb.Params
}
