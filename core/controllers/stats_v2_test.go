package controllers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestGetStatsDefaultQuery(t *testing.T) {
	now := time.Date(2026, time.March, 11, 12, 0, 0, 0, time.UTC)
	query := getStatsDefaultQuery(now)

	orClauses, ok := query["$or"].(bson.A)
	if !assert.True(t, ok) {
		return
	}

	if !assert.Len(t, orClauses, 2) {
		return
	}

	since := now.Add(-statsWindow)

	assert.Equal(t, bson.M{"$gte": since}, orClauses[0].(bson.M)["created_ts"])
	assert.Equal(t, bson.M{"$gte": since}, orClauses[1].(bson.M)["create_ts"])
}

func TestGetRecentTimeQuerySingleField(t *testing.T) {
	since := time.Date(2026, time.March, 1, 0, 0, 0, 0, time.UTC)

	assert.Equal(t, bson.M{
		"created_ts": bson.M{
			"$gte": since,
		},
	}, getRecentTimeQuery(since, "created_ts"))
}
