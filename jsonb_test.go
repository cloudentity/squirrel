package squirrel

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJsonbString(t *testing.T) {
	j := JsonbContains("payload", []string{"name"}, "john")

	sql, args, err := j.ToSql()
	require.NoError(t, err)

	assert.Equal(t, "payload->'name' @> '?'::JSONB", sql)
	assert.Equal(t, []interface{}{`"john"`}, args)
}

func TestJsonbNumber(t *testing.T) {
	j := JsonbContains("payload", []string{"age"}, 30)

	sql, args, err := j.ToSql()
	require.NoError(t, err)

	assert.Equal(t, "payload->'age' @> '?'::JSONB", sql)
	assert.Equal(t, []interface{}{`30`}, args)
}

func TestJsonbBool(t *testing.T) {
	j := JsonbContains("payload", []string{"active"}, true)

	sql, args, err := j.ToSql()
	require.NoError(t, err)

	assert.Equal(t, "payload->'active' @> '?'::JSONB", sql)
	assert.Equal(t, []interface{}{`true`}, args)
}

func TestJsonbNested(t *testing.T) {
	j := JsonbContains("payload", []string{"address", "street"}, "abbey road")

	sql, args, err := j.ToSql()
	require.NoError(t, err)

	assert.Equal(t, "payload->'address'->'street' @> '?'::JSONB", sql)
	assert.Equal(t, []interface{}{`"abbey road"`}, args)
}

func TestJsonbArray(t *testing.T) {
	j := JsonbContains("payload", []string{"roles"}, []string{"admin"})

	sql, args, err := j.ToSql()
	require.NoError(t, err)

	assert.Equal(t, "payload->'roles' @> '?'::JSONB", sql)
	assert.Equal(t, []interface{}{`["admin"]`}, args)
}

func TestJsonbComplex(t *testing.T) {
	j := JsonbContains("payload", []string{}, map[string]interface{}{
		"address": map[string]interface{}{
			"street": "abbey road",
		},
		"age": 30,
	})

	sql, args, err := j.ToSql()
	require.NoError(t, err)

	assert.Equal(t, "payload @> '?'::JSONB", sql)
	assert.Equal(t, []interface{}{`{"address":{"street":"abbey road"},"age":30}`}, args)
}
