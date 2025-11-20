package notes

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrNotFound = errors.New("note not found")

type Repo struct {
	col *mongo.Collection
}

func NewRepo(db *mongo.Database) (*Repo, error) {
	col := db.Collection("notes")

	_, err := col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "title", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil { return nil, err }

	_, err = col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "title", Value: "text"},
			{Key: "content", Value: "text"},
		},
	})
	if err != nil { return nil, err }

	ttl := int32(0)
	_, err = col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "expiresAt", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(ttl),
	})
	if err != nil { return nil, err }

	return &Repo{col: col}, nil
}

func (r *Repo) Create(ctx context.Context, title, content string, expiresAt *time.Time) (Note, error) {
	now := time.Now()
	n := Note{Title: title, Content: content, CreatedAt: now, UpdatedAt: now}
	res, err := r.col.InsertOne(ctx, n)
	if err != nil { return Note{}, err }
	n.ID = res.InsertedID.(primitive.ObjectID)
	return n, nil
}

func (r *Repo) ByID(ctx context.Context, idHex string) (Note, error) {
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil { return Note{}, ErrNotFound }
	var n Note
	if err := r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&n); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) { return Note{}, ErrNotFound }
		return Note{}, err
	}
	return n, nil
}

func (r *Repo) List(ctx context.Context, q string, limit, skip int64) ([]Note, error) {
	filter := bson.M{}
	if q != "" {
		filter = bson.M{"$text": bson.M{"$search": q}}
	}
	opts := options.Find().SetLimit(limit).SetSkip(skip).SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cur, err := r.col.Find(ctx, filter, opts)
	if err != nil { return nil, err }
	defer cur.Close(ctx)

	var out []Note
	for cur.Next(ctx) {
		var n Note
		if err := cur.Decode(&n); err != nil { return nil, err }
		out = append(out, n)
	}
	return out, cur.Err()
}

func (r *Repo) Update(ctx context.Context, idHex string, title, content *string) (Note, error) {
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil { return Note{}, ErrNotFound }

	set := bson.M{"updatedAt": time.Now()}
	if title != nil   { set["title"] = *title }
	if content != nil { set["content"] = *content }

	after := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated Note
	if err := r.col.FindOneAndUpdate(ctx, bson.M{"_id": oid}, bson.M{"$set": set}, after).Decode(&updated); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) { return Note{}, ErrNotFound }
		return Note{}, err
	}
	return updated, nil
}

func (r *Repo) Delete(ctx context.Context, idHex string) error {
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil { return ErrNotFound }
	res, err := r.col.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil { return err }
	if res.DeletedCount == 0 { return ErrNotFound }
	return nil
}

func (r *Repo) Stats(ctx context.Context) (count int64, avgLen float64, err error) {
	// pipeline: [{ $group: { _id: null, count: {$sum:1}, avgLen: {$avg: { $strLenCP: "$content" } } } }]
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "avgLen", Value: bson.D{{Key: "$avg", Value: bson.D{{Key: "$strLenCP", Value: "$content"}}}}},
		}}},
	}

	cur, err := r.col.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, 0, err
	}
	defer cur.Close(ctx)

	var res []bson.M
	if err := cur.All(ctx, &res); err != nil {
		return 0, 0, err
	}
	if len(res) == 0 {
		return 0, 0, nil
	}
	doc := res[0]
	if v, ok := doc["count"].(int32); ok {
		count = int64(v)
	} else if v64, ok := doc["count"].(int64); ok {
		count = v64
	} else if f, ok := doc["count"].(float64); ok {
		count = int64(f)
	}

	if a, ok := doc["avgLen"].(float64); ok {
		avgLen = a
	} else {
		avgLen = 0
	}

	return count, avgLen, nil
}
