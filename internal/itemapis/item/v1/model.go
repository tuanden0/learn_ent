package v1

import "go.mongodb.org/mongo-driver/bson/primitive"

type ErrorResponse struct {
	Message string `json:"message"`
}

type CreateItemRequest struct {
	Owner    string   `json:"owner,omitempty" bson:"owner,omitempty"`
	Name     string   `json:"name,omitempty" bson:"name,omitempty"`
	Price    float64  `json:"price,omitempty" bson:"price,omitempty"`
	Quantity int      `json:"quantity,omitempty" bson:"quantity,omitempty"`
	Tags     []string `json:"tags,omitempty" bson:"tags,omitempty"`
}

type CreateItemResponse struct {
	ObjectID primitive.ObjectID `bson:"_id" json:"_id"`
	Owner    string             `json:"owner,omitempty" bson:"owner,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Price    float64            `json:"price,omitempty" bson:"price,omitempty"`
	Quantity int                `json:"quantity,omitempty" bson:"quantity,omitempty"`
	Tags     []string           `json:"tags,omitempty" bson:"tags,omitempty"`
}

type RetrieveItemResponse struct {
	ObjectID primitive.ObjectID `bson:"_id" json:"_id"`
	Owner    string             `json:"owner,omitempty" bson:"owner,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Price    float64            `json:"price,omitempty" bson:"price,omitempty"`
	Quantity int                `json:"quantity,omitempty" bson:"quantity,omitempty"`
	Tags     []string           `json:"tags,omitempty" bson:"tags,omitempty"`
}

type UpdateItemRequest struct {
	Owner    string   `json:"owner,omitempty" bson:"owner,omitempty"`
	Name     string   `json:"name,omitempty" bson:"name,omitempty"`
	Price    float64  `json:"price,omitempty" bson:"price,omitempty"`
	Quantity int      `json:"quantity,omitempty" bson:"quantity,omitempty"`
	Tags     []string `json:"tags,omitempty" bson:"tags,omitempty"`
}

type UpdateItemResponse struct {
	ObjectID primitive.ObjectID `bson:"_id" json:"_id"`
	Owner    string             `json:"owner,omitempty" bson:"owner,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Price    float64            `json:"price,omitempty" bson:"price,omitempty"`
	Quantity int                `json:"quantity,omitempty" bson:"quantity,omitempty"`
	Tags     []string           `json:"tags,omitempty" bson:"tags,omitempty"`
}
