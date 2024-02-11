package genericimplementation

import (
	"fmt"
	"ovaphlow/cratecyclone/utilities"
	"strings"
)

type Option struct {
	Take int
	Skip int64
}

type Filter struct {
	ArrayContain  []string
	Equal         []string
	Greater       []string
	In            []string
	Lesser        []string
	Like          []string
	ObjectContain []string
	ObjectLike    []string
}

func ArrayContain(arrayContain []string) ([]string, []string) {
	if len(arrayContain)%2 != 0 {
		utilities.Slogger.Warn("arrayContain length is not even")
		return []string{}, []string{}
	}
	var conditions []string
	var params []string
	for i := 0; i < len(arrayContain); i += 2 {
		conditions = append(
			conditions,
			fmt.Sprintf("%s @> jsonb_build_array('%s')", arrayContain[i], arrayContain[i+1]),
		)
	}
	return conditions, params
}

func Equal(equal []string) ([]string, []string) {
	if len(equal)%2 != 0 {
		utilities.Slogger.Warn("equal length is not even")
		return []string{}, []string{}
	}
	var conditions []string
	var params []string
	for i := 0; i < len(equal); i += 2 {
		conditions = append(conditions, fmt.Sprintf("%s = ?", equal[i]))
		params = append(params, equal[i+1])
	}
	return conditions, params
}

func Greater(greater []string) ([]string, []string) {
	if len(greater)%2 != 0 {
		utilities.Slogger.Warn("greater length is not even")
		return []string{}, []string{}
	}
	var conditions []string
	var params []string
	for i := 0; i < len(greater); i += 2 {
		conditions = append(
			conditions,
			fmt.Sprintf("%s >= ?", greater[i]),
		)
		params = append(params, greater[i+1])
	}
	return conditions, params
}

func In(in []string) ([]string, []string) {
	if len(in) < 2 {
		utilities.Slogger.Warn("in length is less than 2")
		return []string{}, []string{}
	}
	c := make([]string, len(in)-1)
	for i := range c {
		c[i] = "?"
	}
	return []string{fmt.Sprintf("%s in (%s)", in[0], strings.Join(c, ", "))}, in[1:]
}

func Lesser(lesser []string) ([]string, []string) {
	if len(lesser)%2 != 0 {
		utilities.Slogger.Warn("lesser length is not even")
		return []string{}, []string{}
	}
	var conditions []string
	var params []string
	for i := 0; i < len(lesser); i += 2 {
		conditions = append(
			conditions,
			fmt.Sprintf("%s <= ?", lesser[i]),
		)
		params = append(params, lesser[i+1])
	}
	return conditions, params
}

func Like(like []string) ([]string, []string) {
	if len(like)%2 != 0 {
		utilities.Slogger.Warn("like length is not even")
		return []string{}, []string{}
	}
	var conditions []string
	var params []string
	for i := 0; i < len(like); i += 2 {
		conditions = append(
			conditions,
			fmt.Sprintf("position(? in %s) > 0", like[i]),
		)
		params = append(params, like[i+1])
	}
	return conditions, params
}

func ObjectContain(objectContain []string) ([]string, []string) {
	if len(objectContain)%3 != 0 {
		utilities.Slogger.Warn("objectContain length is not multiple of 3")
		return []string{}, []string{}
	}
	var conditions []string
	var params []string
	for i := 0; i < len(objectContain); i += 3 {
		conditions = append(
			conditions,
			fmt.Sprintf(
				"%s @> jsonb_build_object('%s', '%s')",
				objectContain[i],
				objectContain[i+1],
				objectContain[i+2],
			),
		)
	}
	return conditions, params
}

func ObjectLike(objectLike []string) ([]string, []string) {
	if len(objectLike)%3 != 0 {
		utilities.Slogger.Warn("objectLike length is not multiple of 3")
		return []string{}, []string{}
	}
	var conditions []string
	var params []string
	for i := 0; i < len(objectLike); i += 3 {
		conditions = append(
			conditions,
			fmt.Sprintf(
				"position(? in %s->>'%s') > 0",
				objectLike[i],
				objectLike[i+1],
			),
		)
		params = append(params, objectLike[i+2])
	}
	return conditions, params
}
