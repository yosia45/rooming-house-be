package repositories

import (
	"rooming-house-cms-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(transaction *models.Transaction) error
	FindAllTransactions() (*[]models.Transaction, error)
	FindTransactionByID(id uuid.UUID) (*models.Transaction, error)
	DeleteTransactionByID(id uuid.UUID) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (t *transactionRepository) CreateTransaction(transaction *models.Transaction) error {
	if err := t.db.Create(transaction).Error; err != nil {
		return err
	}
	return nil
}

func (t *transactionRepository) FindAllTransactions() (*[]models.Transaction, error) {
	var transactions []models.Transaction
	if err := t.db.Find(&transactions).Error; err != nil {
		return nil, err
	}
	return &transactions, nil
}

func (t *transactionRepository) FindTransactionByID(id uuid.UUID) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := t.db.Where("id = ?", id).First(&transaction).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (t *transactionRepository) DeleteTransactionByID(id uuid.UUID) error {
	if err := t.db.Where("id = ?", id).Delete(&models.Transaction{}).Error; err != nil {
		return err
	}
	return nil
}
