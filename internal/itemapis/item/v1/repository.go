package v1

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	Create(ctx context.Context, in *CreateItemRequest) (*CreateItemResponse, error)
	Retrieve(ctx context.Context, id string) (*RetrieveItemResponse, error)
	Update(ctx context.Context, id string, in *UpdateItemRequest) (*UpdateItemResponse, error)
	Delete(ctx context.Context, id string) error
}

type repoManager struct {
	items *mongo.Collection
}

func NewRepoManager(items *mongo.Collection) Repository {
	return &repoManager{
		items: items,
	}
}

func (r *repoManager) Create(ctx context.Context, in *CreateItemRequest) (*CreateItemResponse, error) {

	// Convert struct to Mongo Doc
	doc, err := toDoc(in)
	if err != nil {
		return nil, err
	}

	oid, err := r.items.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}

	out := &CreateItemResponse{
		ObjectID: oid.InsertedID.(primitive.ObjectID),
		Owner:    in.Owner,
		Name:     in.Name,
		Price:    in.Price,
		Quantity: in.Quantity,
		Tags:     in.Tags,
	}

	return out, nil
}

func (r *repoManager) Retrieve(ctx context.Context, id string) (*RetrieveItemResponse, error) {

	// Convert hex id to mongo ObjectID
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	out := &RetrieveItemResponse{}
	filter := bson.D{
		{Key: "_id", Value: oid},
	}
	if err := r.items.FindOne(ctx, filter).Decode(out); err != nil {
		return nil, err
	}

	return out, nil
}

func (r *repoManager) Update(ctx context.Context, id string, in *UpdateItemRequest) (*UpdateItemResponse, error) {

	// Convert id string to ObjectID
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Filter ObjectID
	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	// Prepare update doc
	update := bson.D{{Key: "$set", Value: in}}

	// Using After option to get return updated doc
	after := options.After
	opt := &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	// Find One and Update data
	// Using FindOneAndUpdate instead of UpdateOne or UpdateByID
	// because FindOneAndUpdate return updated docs while
	// UpdateOne and UpdateByID return ObjectID
	out := &UpdateItemResponse{}
	if err := r.items.FindOneAndUpdate(ctx, filter, update, opt).Decode(out); err != nil {
		return nil, err
	}

	return out, nil
}

func (r *repoManager) Delete(ctx context.Context, id string) error {

	// Convert id string to ObjectID
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Filter ObjectID
	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	if _, err := r.items.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}
