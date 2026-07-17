package daily

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDailyRecordRepo struct {
	col *mongo.Collection
}

func NewDailyRecordRepository(db *mongo.Database) *MongoDailyRecordRepo {
	return &MongoDailyRecordRepo{col: db.Collection("dailyRecords")}
}

func (r *MongoDailyRecordRepo) EnsureIndexes(ctx context.Context) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "date", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := r.col.Indexes().CreateOne(ctx, indexModel)
	return err
}

func (r *MongoDailyRecordRepo) FindByDate(ctx context.Context, date string) (*DailyRecord, error) {
	var record DailyRecord
	err := r.col.FindOne(ctx, bson.M{"date": date}).Decode(&record)
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
	filter := bson.M{"date": date}
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
func (r *MongoDailyRecordRepo) SumTaskProgressByIDs(ctx context.Context, taskIDs []primitive.ObjectID) (map[primitive.ObjectID]int, error) {
	if len(taskIDs) == 0 {
		return make(map[primitive.ObjectID]int), nil
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"tasks.taskId": bson.M{"$in": taskIDs}}}},
		{{Key: "$unwind", Value: "$tasks"}},
		{{Key: "$match", Value: bson.M{"tasks.taskId": bson.M{"$in": taskIDs}}}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$tasks.taskId",
			"total": bson.M{"$sum": "$tasks.actualAmount"},
		}}},
	}

	cursor, err := r.col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		ID    primitive.ObjectID `bson:"_id"`
		Total int                `bson:"total"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	progressMap := make(map[primitive.ObjectID]int)
	for _, r := range results {
		progressMap[r.ID] = r.Total
	}
	return progressMap, nil
}

func (r *MongoDailyRecordRepo) RemoveTaskFromRecordsBeforeDate(ctx context.Context, taskID primitive.ObjectID, date string) error {
	filter := bson.M{"date": bson.M{"$lt": date}}
	update := bson.M{"$pull": bson.M{"tasks": bson.M{"taskId": taskID}}}

	_, err := r.col.UpdateMany(ctx, filter, update)
	return err
}

