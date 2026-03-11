package spider

import (
	"errors"

	models "github.com/crawlab-team/crawlab/core/models/models"
	models2 "github.com/crawlab-team/crawlab/core/models/models/v2"
	"github.com/crawlab-team/crawlab/core/models/service"
	"github.com/crawlab-team/crawlab/core/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
)

func EnsureDataCollectionForSpiderV2(s *models2.SpiderV2) error {
	if s == nil {
		return nil
	}

	dc, err := ensureDataCollection(s.ColId, s.ColName, s.Name)
	if err != nil || dc == nil {
		return err
	}

	if s.ColId == dc.Id && s.ColName == dc.Name {
		return nil
	}

	s.ColId = dc.Id
	s.ColName = dc.Name

	if s.Id.IsZero() {
		return nil
	}

	return service.NewModelServiceV2[models2.SpiderV2]().UpdateById(s.Id, bson.M{
		"$set": bson.M{
			"col_id":   s.ColId,
			"col_name": s.ColName,
		},
	})
}

func EnsureDataCollectionForSpider(s *models.Spider) error {
	if s == nil {
		return nil
	}

	dc, err := ensureDataCollection(s.ColId, s.ColName, s.Name)
	if err != nil || dc == nil {
		return err
	}

	if s.ColId == dc.Id && s.ColName == dc.Name {
		return nil
	}

	s.ColId = dc.Id
	s.ColName = dc.Name

	if s.Id.IsZero() {
		return nil
	}

	return service.NewModelServiceV2[models2.SpiderV2]().UpdateById(s.Id, bson.M{
		"$set": bson.M{
			"col_id":   s.ColId,
			"col_name": s.ColName,
		},
	})
}

func ensureDataCollection(colId primitive.ObjectID, colName, spiderName string) (*models2.DataCollectionV2, error) {
	dcSvc := service.NewModelServiceV2[models2.DataCollectionV2]()

	if !colId.IsZero() {
		dc, err := dcSvc.GetById(colId)
		if err == nil {
			return dc, nil
		}
		if !errors.Is(err, mongo2.ErrNoDocuments) {
			return nil, err
		}
	}

	name := utils.GetSpiderCol(colName, spiderName)
	if name == "" {
		return nil, nil
	}

	dc, err := dcSvc.GetOne(bson.M{"name": name}, nil)
	if err == nil {
		return dc, nil
	}
	if !errors.Is(err, mongo2.ErrNoDocuments) {
		return nil, err
	}

	dc = &models2.DataCollectionV2{Name: name}
	id, err := dcSvc.InsertOne(*dc)
	if err != nil {
		return nil, err
	}
	dc.SetId(id)

	return dc, nil
}
