package squirrel

import (
	"bytes"
	"errors"

	"github.com/lann/builder"
)

func init() {
	builder.Register(SetBuilder{}, setData{})
}

type SetBuilder builder.Builder

type Intersect []Sqlizer

func (i Intersect) ToSql() (sql string, args []interface{}, err error) {
	return setToSql(i, "INTERSECT")
}

type Union []Sqlizer

func (i Union) ToSql() (sql string, args []interface{}, err error) {
	return setToSql(i, "UNION")
}

func (b SetBuilder) Query(q Sqlizer) SetBuilder {
	return builder.Set(b, "Query", q).(SetBuilder)
}

func (b SetBuilder) ToSql() (sql string, args []interface{}, err error) {
	data := builder.GetStruct(b).(setData)
	return data.ToSql()
}

func (b SetBuilder) PlaceholderFormat(fmt PlaceholderFormat) SetBuilder {
	return builder.Set(b, "PlaceholderFormat", fmt).(SetBuilder)
}

type setData struct {
	Query             Sqlizer
	PlaceholderFormat PlaceholderFormat
}

func (d *setData) ToSql() (sql string, args []interface{}, err error) {
	return d.Query.ToSql()
}

func setToSql(queries []Sqlizer, op string) (sql string, args []interface{}, err error) {
	var (
		selArgs []interface{}
		selSql  string
		sqlBuf  = &bytes.Buffer{}
	)

	if len(queries) == 0 {
		err = errors.New("require a minimum of 1 query in set builder")
		return sql, args, err
	}

	for index, selector := range queries {
		switch s := selector.(type) {
		case SelectBuilder:
			selector = s.PlaceholderFormat(Question)
		}

		selSql, selArgs, err = selector.ToSql()

		if err != nil {
			return sql, args, err
		}

		switch selector.(type) {
		case Intersect, Union:
			if index == 0 {
				sqlBuf.WriteString("(" + selSql + ")")
			} else {
				sqlBuf.WriteString(" " + op + " (" + selSql + ")")
			}

		default:
			if index == 0 {
				sqlBuf.WriteString(selSql)
			} else {
				sqlBuf.WriteString(" " + op + " " + selSql)
			}
		}

		args = append(args, selArgs...)
	}

	return sqlBuf.String(), args, err
}
