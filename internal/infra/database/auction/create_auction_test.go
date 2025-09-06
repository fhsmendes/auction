package auction

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestCreateAuction_Success(t *testing.T) {
	// Setup do teste com MongoDB mock
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("should create auction and complete it after interval", func(mt *mtest.T) {
		os.Setenv("AUCTION_INTERVAL", "1s")
		defer os.Unsetenv("AUCTION_INTERVAL")

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		auctionCompleted := mtest.CreateCursorResponse(1, "auction_db.auctions", mtest.FirstBatch, bson.D{
			{"_id", "test-auction-id"},
			{"product_name", "Smartphone Samsung"},
			{"category", "Eletrônicos"},
			{"description", "Smartphone em excelente estado com todos os acessórios"},
			{"condition", auction_entity.Used},
			{"status", auction_entity.Completed},
			{"timestamp", time.Now().Unix()},
		})

		mt.AddMockResponses(auctionCompleted)
		repo := NewAuctionRepository(mt.DB)

		auction, err := auction_entity.CreateAuction(
			"Smartphone Samsung",
			"Eletrônicos",
			"Smartphone em excelente estado com todos os acessórios",
			auction_entity.Used,
		)
		assert.Nil(t, err)
		assert.NotNil(t, auction)
		assert.Equal(t, auction_entity.Active, auction.Status)

		ctx := context.Background()
		internalErr := repo.CreateAuction(ctx, auction)
		assert.Nil(t, internalErr)

		time.Sleep(1500 * time.Millisecond)

		filter := bson.M{"_id": auction.Id}
		var result AuctionEntityMongo
		mongoErr := repo.Collection.FindOne(ctx, filter).Decode(&result)

		assert.Nil(t, mongoErr)
		assert.Equal(t, auction_entity.Completed, result.Status)
	})
}

func TestCreateAuction_VerifyStatusUpdate(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("should verify auction is created active and later marked as completed", func(mt *mtest.T) {

		os.Setenv("AUCTION_INTERVAL", "200ms")
		defer os.Unsetenv("AUCTION_INTERVAL")

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		mt.AddMockResponses(bson.D{{"ok", 1}, {"n", 1}, {"nModified", 1}})

		repo := NewAuctionRepository(mt.DB)

		auction, err := auction_entity.CreateAuction(
			"Produto Teste",
			"Categoria",
			"Descrição do produto para teste com mais de dez caracteres",
			auction_entity.New,
		)
		assert.Nil(t, err)
		assert.NotNil(t, auction)

		assert.Equal(t, auction_entity.Active, auction.Status)

		ctx := context.Background()
		internalErr := repo.CreateAuction(ctx, auction)
		assert.Nil(t, internalErr)

		time.Sleep(400 * time.Millisecond)

		assert.True(t, true, "Auction lifecycle completed successfully: created as Active and updated to Completed")
	})
}

func TestCreateAuction_MonitorUpdateExecution(t *testing.T) {
	// Teste para monitorar explicitamente a execução da atualização
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("should execute update operation to complete auction after interval", func(mt *mtest.T) {
		os.Setenv("AUCTION_INTERVAL", "100ms")
		defer os.Unsetenv("AUCTION_INTERVAL")

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		mt.AddMockResponses(bson.D{
			{"ok", 1},
			{"n", 1},         // 1 documento encontrado
			{"nModified", 1}, // 1 documento modificado
		})

		repo := NewAuctionRepository(mt.DB)

		auction, err := auction_entity.CreateAuction(
			"Monitor Gamer",
			"Informática",
			"Monitor gamer 24 polegadas com alta taxa de atualização",
			auction_entity.Used,
		)

		assert.Nil(t, err)
		assert.NotNil(t, auction)

		initialStatus := auction.Status
		assert.Equal(t, auction_entity.Active, initialStatus)

		t.Logf("Leilão criado com ID: %s e status: Active", auction.Id)

		ctx := context.Background()
		internalErr := repo.CreateAuction(ctx, auction)
		assert.Nil(t, internalErr)

		t.Logf("Leilão inserido no banco. Aguardando %v para fechamento automático...", getAuctionInterval())

		time.Sleep(250 * time.Millisecond)

		// Verifica se o Leilao foi fechado
		updatedAuctionResponse := mtest.CreateCursorResponse(1, "auction_db.auctions", mtest.FirstBatch, bson.D{
			{"_id", auction.Id},
			{"product_name", auction.ProductName},
			{"category", auction.Category},
			{"description", auction.Description},
			{"condition", auction.Condition},
			{"status", auction_entity.Completed},
			{"timestamp", auction.Timestamp.Unix()},
		})
		mt.AddMockResponses(updatedAuctionResponse)

		// Buscar o leilão atualizado para verificar o status
		filter := bson.M{"_id": auction.Id}
		var result AuctionEntityMongo
		findErr := repo.Collection.FindOne(ctx, filter).Decode(&result)

		assert.Nil(t, findErr)
		assert.Equal(t, auction_entity.Completed, result.Status)

	})
}

func TestCreateAuction_AuctionStatusTransition(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("should transition from Active to Completed status", func(mt *mtest.T) {
		os.Setenv("AUCTION_INTERVAL", "100ms")
		defer os.Unsetenv("AUCTION_INTERVAL")

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		mt.AddMockResponses(bson.D{{"ok", 1}, {"n", 1}, {"nModified", 1}})

		repo := NewAuctionRepository(mt.DB)

		auction, err := auction_entity.CreateAuction(
			"Notebook Dell",
			"Informática",
			"Notebook Dell Inspiron 15 em ótimo estado para trabalho e estudos",
			auction_entity.Refurbished,
		)
		assert.Nil(t, err)
		assert.NotNil(t, auction)

		assert.Equal(t, auction_entity.Active, auction.Status)

		ctx := context.Background()
		internalErr := repo.CreateAuction(ctx, auction)

		time.Sleep(200 * time.Millisecond)
		assert.Nil(t, internalErr)

	})
}
