package squirrel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetBuilder(t *testing.T) {
	builder := StatementBuilder.PlaceholderFormat(Dollar)

	fromIdentifiers := builder.Select("user_id").
		From("user_identifiers").
		Where(Eq{
			"tenant_id":    "default",
			"user_pool_id": "pool",
		}).Where(Like{
		"identifier_value": "%",
	})

	fromAddresses := builder.Select("user_id").
		From("user_verifiable_addresses").
		Where(Eq{
			"tenant_id":    "default",
			"user_pool_id": "pool",
		}).Where(Like{
		"identifier_value": "%",
	})

	fromPayload := builder.Select("user_id").
		From("users").
		Where(Eq{
			"tenant_id":    "default",
			"user_pool_id": "pool",
		}).Where(Eq{
		"payload#>'name'": "andrzej",
	})

	fromMetadata := builder.Select("user_id").
		From("users").
		Where(Eq{
			"tenant_id":    "default",
			"user_pool_id": "pool",
		}).Where(Eq{
		"metadata#>'name'": "andrzej",
	})

	b := builder.
		Select("u.user_id").
		FromSet(builder.Set(Intersect{
			Union{fromIdentifiers, fromAddresses},
			Intersect{fromPayload, fromMetadata},
		}), "u").Where(Gt{
		"user_id": "user100",
	}).OrderBy("u.user_id").
		Limit(10)

	sql, args, err := b.ToSql()
	assert.NoError(t, err)

	expectedSql := "SELECT u.user_id FROM ((" +
		"SELECT user_id FROM user_identifiers WHERE tenant_id = $1 AND user_pool_id = $2 AND identifier_value LIKE $3 " +
		"UNION SELECT user_id FROM user_verifiable_addresses WHERE tenant_id = $4 AND user_pool_id = $5 AND identifier_value LIKE $6) " +
		"INTERSECT (SELECT user_id FROM users WHERE tenant_id = $7 AND user_pool_id = $8 AND payload#>'name' = $9 " +
		"INTERSECT SELECT user_id FROM users WHERE tenant_id = $10 AND user_pool_id = $11 AND metadata#>'name' = $12)) " +
		"AS u WHERE user_id > $13 ORDER BY u.user_id LIMIT 10"
	assert.Equal(t, expectedSql, sql)

	expectedArgs := []interface{}{"default", "pool", "%", "default", "pool", "%", "default", "pool", "andrzej", "default", "pool", "andrzej", "user100"}
	assert.Equal(t, expectedArgs, args)
}
