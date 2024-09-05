package service

import (
	"context"
	"fmt"

	"github.com/CP-Payne/ecomstore/internal/config"
	"github.com/CP-Payne/ecomstore/internal/database"
	"github.com/CP-Payne/ecomstore/internal/models"
	"github.com/CP-Payne/ecomstore/internal/utils/apperrors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ProductService struct {
	logger *zap.Logger
	db     *database.Queries
}

func NewProductService(db *database.Queries) *ProductService {
	return &ProductService{
		logger: config.GetLogger(),
		db:     db,
	}
}

func (s *ProductService) ProductExists(ctx context.Context, id uuid.UUID) (bool, error) {
	exists, err := s.db.ProductExists(ctx, id)
	if err != nil {
		s.logger.Error("failed to check if product exists", zap.Error(err), zap.String("productID", id.String()))
		return false, fmt.Errorf("failed to check if product exists: %w", err)
	}

	return exists, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id uuid.UUID) (models.Product, error) {
	product, err := s.db.GetProduct(ctx, id)
	if err != nil {
		if apperrors.IsNoRowsError(err) {
			return models.Product{}, apperrors.ErrNotFound
		}
		s.logger.Error("failed to retrieve product from db", zap.Error(err), zap.String("productID", id.String()))
		return models.Product{}, fmt.Errorf("failed to retrieve product: %w", err)

	}

	p, ok := models.DatabaseProductToProduct(product, false).(models.Product)
	if !ok {
		return models.Product{}, fmt.Errorf("failed to process interface into product: %w", err)
	}

	return p, nil
}

func (s *ProductService) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	products, err := s.db.GetAllProducts(ctx)
	if err != nil {
		s.logger.Error("failed to retrieve products from db", zap.Error(err))
		return nil, fmt.Errorf("failed to retrieve and process products: %w", err)
	}

	pl, ok := models.DatabaseProductsToProducts(products, false).([]models.Product)
	if !ok {
		return []models.Product{}, fmt.Errorf("failed to process interface into product slice: %w", err)
	}

	return pl, nil
}

func (s *ProductService) GetProductCategories(ctx context.Context) ([]models.Category, error) {
	categories, err := s.db.GetProductCategories(ctx)
	if err != nil {
		s.logger.Error("failed to retrieve product categories from db", zap.Error(err))
		return nil, fmt.Errorf("failed to retrieve and process products: %w", err)
	}
	return models.DatabaseCategoriesToCategories(categories), nil
}

func (s *ProductService) GetProductsByCategory(ctx context.Context, categoryID uuid.UUID) ([]models.Product, error) {
	products, err := s.db.GetProductsByCategory(ctx, categoryID)
	if err != nil {
		if apperrors.IsNoRowsError(err) {
			return []models.Product{}, fmt.Errorf("no products found for category: %w", apperrors.ErrNotFound)
		}

		s.logger.Error("failed to retrieve products by category from db", zap.Error(err), zap.String("CategoryID", categoryID.String()))
		return []models.Product{}, fmt.Errorf("failed to retrieve products by category: %w", err)
	}

	pl, ok := models.DatabaseProductsToProducts(products, false).([]models.Product)
	if !ok {
		return []models.Product{}, fmt.Errorf("failed to process interface into product slice: %w", err)
	}
	return pl, nil
}
