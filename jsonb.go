package squirrel

import (
	"encoding/json"
	"fmt"
	"strings"
)

type jsonbContains struct {
	column string
	path   []string
	value  interface{}
}

func JsonbContains(column string, path []string, value interface{}) Sqlizer {
	return jsonbContains{column: column, path: path, value: value}
}

func (j jsonbContains) ToSql() (sql string, args []interface{}, err error) {
	var (
		buf strings.Builder
		v   []byte
	)

	buf.WriteString(j.column)

	for _, attr := range j.path {
		buf.WriteString(fmt.Sprintf("->'%s'", attr))
	}

	buf.WriteString(" @> ")

	if v, err = json.Marshal(j.value); err != nil {
		return "", args, err
	}

	buf.WriteString("?::JSONB")

	return buf.String(), []interface{}{string(v)}, nil
}
