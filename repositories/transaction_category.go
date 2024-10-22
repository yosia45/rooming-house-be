package repositories

import (
	"rooming-house-cms-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionCategoryRepository interface {
	CreateTransactionCategory(transactionCategory *models.TransactionCategory) error
	FindTransactionCategoryByID(id uuid.UUID) (*models.TransactionCategory, error)
	FindAllTransactionCategories() (*[]models.TransactionCategory, error)
	UpdateTransactionCategoryByID(transaction *models.TransactionCategory, id uuid.UUID) error
	DeleteTransactionCategoryByID(id uuid.UUID) error
}

type transactionCategoryRepository struct {
	db *gorm.DB
}

func NewTransactionCategoryRepository(db *gorm.DB) TransactionCategoryRepository {
	return &transactionCategoryRepository{db: db}
}

func (r *transactionCategoryRepository) CreateTransactionCategory(transactionCategory *models.TransactionCategory) error {
	if err := r.db.Create(transactionCategory).Error; err != nil {
		return err
	}
	return nil
}

func (r *transactionCategoryRepository) FindTransactionCategoryByID(id uuid.UUID) (*models.TransactionCategory, error) {
	var transactionCategory models.TransactionCategory
	if err := r.db.Where("id = ?", id).First(&transactionCategory).Error; err != nil {
		return nil, err
	}
	return &transactionCategory, nil
}

func (r *transactionCategoryRepository) FindAllTransactionCategories() (*[]models.TransactionCategory, error) {
	var transactionCategories []models.TransactionCategory
	if err := r.db.Find(&transactionCategories).Error; err != nil {
		return nil, err
	}
	return &transactionCategories, nil
}

func (r *transactionCategoryRepository) UpdateTransactionCategoryByID(transactionCategory *models.TransactionCategory, id uuid.UUID) error {
	res := r.db.Where("id = ?", id).Updates(transactionCategory)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *transactionCategoryRepository) DeleteTransactionCategoryByID(id uuid.UUID) error {
	res := r.db.Where("id = ?", id).Delete(&models.TransactionCategory{})
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
