package repository

import (
	"context"

	"daily-seed/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DailyRecordRepository interface {
	FindByDate(ctx context.Context, date string) (*model.DailyRecord, error)
	Upsert(ctx context.Context, record *model.DailyRecord) error
	PatchByDate(ctx context.Context, date string, setFields bson.M) error
}

type mongoDailyRecordRepo struct {
	col *mongo.Collection
}

func NewDailyRecordRepository(db *mongo.Database) DailyRecordRepository {
	return &mongoDailyRecordRepo{col: db.Collection("dailyRecords")}
}

func (r *mongoDailyRecordRepo) FindByDate(ctx context.Context, date string) (*model.DailyRecord, error) {
	var record model.DailyRecord
	err := r.col.FindOne(ctx, bson.M{"_id": date}).Decode(&record)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

func (r *mongoDailyRecordRepo) Upsert(ctx context.Context, record *model.DailyRecord) error {
	filter := bson.M{"_id": record.ID}
	update := bson.M{"$set": record}
	opts := options.Update().SetUpsert(true)

	_, err := r.col.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *mongoDailyRecordRepo) PatchByDate(ctx context.Context, date string, setFields bson.M) error {
	filter := bson.M{"_id": date}
	update := bson.M{"$set": setFields}

	result, err := r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
