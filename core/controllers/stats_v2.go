package controllers

import (
	"time"

	"github.com/crawlab-team/crawlab/core/stats"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

const statsWindow = 30 * 24 * time.Hour

func getStatsDefaultQuery(now time.Time) bson.M {
	return getRecentTimeQuery(now.Add(-statsWindow), "created_ts", "create_ts")
}

func getRecentTimeQuery(since time.Time, fieldNames ...string) bson.M {
	if len(fieldNames) == 1 {
		return bson.M{
			fieldNames[0]: bson.M{
				"$gte": since,
			},
		}
	}

	clauses := make(bson.A, 0, len(fieldNames))
	for _, fieldName := range fieldNames {
		clauses = append(clauses, bson.M{
			fieldName: bson.M{
				"$gte": since,
			},
		})
	}

	return bson.M{
		"$or": clauses,
	}
}

func GetStatsOverview(c *gin.Context) {
	data, err := stats.GetStatsService().GetOverviewStats(getStatsDefaultQuery(time.Now()))
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccessWithData(c, data)
}

func GetStatsDaily(c *gin.Context) {
	data, err := stats.GetStatsService().GetDailyStats(getStatsDefaultQuery(time.Now()))
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccessWithData(c, data)
}

func GetStatsTasks(c *gin.Context) {
	data, err := stats.GetStatsService().GetTaskStats(getStatsDefaultQuery(time.Now()))
	if err != nil {
		HandleErrorInternalServerError(c, err)
		return
	}
	HandleSuccessWithData(c, data)
}
