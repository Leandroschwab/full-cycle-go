package auction

import (
	"context"
	"os"
	"testing"
	"time"

	"fullcycle-auction_go/internal/entity/auction_entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestAuctionAutoClose(t *testing.T) {
	os.Setenv("AUCTION_INTERVAL", "2s")

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://admin:admin@localhost:27017/"))
	if err != nil {
		t.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("auction_test_db")
	defer db.Collection("auctions").Drop(ctx)

	repo := NewAuctionRepository(db)
	auction := &auction_entity.Auction{
		Id:          "test-auction-id",
		ProductName: "Test Product",
		Category:    "Test Category",
		Description: "Test Description",
		Condition:   auction_entity.New,
		Status:      auction_entity.Active,
		Timestamp:   time.Now(),
	}

	repo.CreateAuction(ctx, auction)

	time.Sleep(3 * time.Second)

	var result AuctionEntityMongo
	err = repo.Collection.FindOne(ctx, bson.M{"_id": auction.Id}).Decode(&result)
	if err != nil {
		t.Fatalf("error fetching auction: %v", err)
	}
	if result.Status != auction_entity.Completed {
		t.Errorf("expected auction status to be Completed, got %v", result.Status)
	}
}
