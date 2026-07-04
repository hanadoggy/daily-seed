package daily

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDailyRecordRepo struct {
	col *mongo.Collection
}

func NewDailyRecordRepository(db *mongo.Database) *MongoDailyRecordRepo {
	return &MongoDailyRecordRepo{col: db.Collection("dailyRecords")}
}

func (r *MongoDailyRecordRepo) FindByDate(ctx context.Context, date string) (*DailyRecord, error) {
	var record DailyRecord
	err := r.col.FindOne(ctx, bson.M{"_id": date}).Decode(&record)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

func (r *MongoDailyRecordRepo) Upsert(ctx context.Context, record *DailyRecord) error {
	filter := bson.M{"_id": record.ID}
	update := bson.M{"$set": record}
	opts := options.Update().SetUpsert(true)

	_, err := r.col.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *MongoDailyRecordRepo) PatchByDate(ctx context.Context, date string, setFields bson.M) error {
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

// SumTaskProgress aggregates the total actualAmount for a given taskId across
// all daily records. This satisfies the task.TaskProgressAggregator interface.
func (r *MongoDailyRecordRepo) SumTaskProgress(ctx context.Context, taskID string) (int, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: "$tasks"}},
		{{Key: "$match", Value: bson.M{"tasks.taskId": taskID}}},
		{{Key: "$group", Value: bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$tasks.actualAmount"},
		}}},
	}

	cursor, err := r.col.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		Total int `bson:"total"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return 0, err
	}

	if len(results) == 0 {
		return 0, nil
	}
	return results[0].Total, nil
}

