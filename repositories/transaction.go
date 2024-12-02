package repositories

import (
	"rooming-house-cms-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(transaction *models.Transaction) error
	FindAllTransactions(roomingHouseIDs []uuid.UUID, year int) (*[]models.TransactionResponse, error)
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

func (t *transactionRepository) FindAllTransactions(roomingHouseIDs []uuid.UUID, year int) (*[]models.TransactionResponse, error) {
	var transactions []models.TransactionResponse

	query := t.db.Table("transactions t").
		Select("t.id, t.day, t.month, t.year, t.amount, t.rooming_house_id AS rooming_house_id, rh.name AS rooming_house_name, tc.name AS transaction_category_name, tc.is_expense AS transaction_category_is_expense").
		Joins("JOIN rooming_houses rh ON t.rooming_house_id = rh.id").
		Joins("JOIN transaction_categories tc ON t.transaction_category_id = tc.id").
		Where("t.rooming_house_id IN (?)", roomingHouseIDs)

	if year != 0 {
		query = query.Where("t.year = ?", year)
	}

	if err := query.Find(&transactions).Error; err != nil {
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
