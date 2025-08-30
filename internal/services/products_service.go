package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Yuhribrp/gobid/internal/services"
	"github.com/Yuhribrp/gobid/internal/store/pgstore"
)

type ProductsService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewProductsService(pool *pgxpool.Pool) ProductsService {
	return ProductsService{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

func (ps *ProductsService) CreateProduct(
	ctx context.Context,
	sellerId uuid.UUID,
	productName string,
	description string,
	baseprice float64,
	auctionEnd time.Time,
) (uuid.UUID, error) {
	id, err := ps.queries.CreateProduct(ctx, pgstore.CreateProductParams{
		SellerID:    sellerId,
		ProductName: productName,
		Description: description,
		Baseprice:   baseprice,
		AuctionEnd:  auctionEnd,
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	auctionRoom := services.NewAuctionRoom(ctx, id, services.NewBidsService(ps.pool))
	go auctionRoom.Run()

	ps.pool.Exec(ctx, "INSERT INTO auction_rooms (id) VALUES ($1)", id)

	return id, nil
}

var ErrProductNotFound = errors.New("product not found")

func (ps *ProductsService) GetProductById(
	ctx context.Context,
	productId uuid.UUID,
) (pgstore.Product, error) {
	product, err := ps.queries.GetProductByID(ctx, productId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgstore.Product{}, ErrProductNotFound
		}
		return pgstore.Product{}, err
	}

	return product, nil
}
