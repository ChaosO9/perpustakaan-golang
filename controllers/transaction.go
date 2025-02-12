package controllers

import (
	"net/http"
	"time"

	"perpustakaan-golang/models"
	"perpustakaan-golang/types"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BorrowBookRequest struct {
    IDAnggota    int64 `json:"id_anggota" binding:"required"`
    IDBuku       string `json:"id_buku" binding:"required"`
    LoanDuration string    `json:"loan_duration" binding:"required"` 
}

type ReturnBookRequest struct {
	TransactionID uuid.UUID `json:"transaction_id" binding:"required"`
}


type BorrowBookResponse struct {
	Message string `json:"message"`
	ID uuid.UUID `json:"id"` // or int64 if that's your ID type
}

type ReturnBookResponse struct {
	Message string `json:"message"`
}

type GetTransactionsByMemberIDResponse struct {
	Transactions []models.TransactionSwagger `json:"transactions"`
}

type TransactionController struct{}

var transactionModel = models.TransactionModel{}

// Borrow handles borrowing a book.
// @Summary Borrow Book
// @Description Borrow a book by member and book IDs.
// @Tags Transaction
// @Accept json
// @Produce json
// @Param request body controllers.BorrowBookRequest true "Borrow Request"
// @Success 200 {object} types.StandardResponse{data=controllers.BorrowBookResponse}
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /transactions/borrow [post]
func (tc TransactionController) Borrow(c *gin.Context) {
	var borrowRequest BorrowBookRequest
	if err := c.ShouldBindJSON(&borrowRequest); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	idBuku, err := uuid.Parse(borrowRequest.IDBuku)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid IDBuku"})
		return
	}

	loanDuration, err := time.ParseDuration(borrowRequest.LoanDuration)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid loan duration"})
		return
	}

	transactionID, err := transactionModel.Borrow(borrowRequest.IDAnggota, idBuku, loanDuration)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, types.StandardResponse{
		Code:    http.StatusOK,
		Message: "Book borrowed successfully",
		Data:    &BorrowBookResponse{Message: "Book borrowed successfully", ID: transactionID},
	})
}

// Return handles returning a book.
// @Summary Return Book
// @Description Return a borrowed book by transaction ID.
// @Tags Transaction
// @Accept json
// @Produce json
// @Param request body controllers.ReturnBookRequest true "Return Request" 
// @@Success 200 {object} types.StandardResponse{data=controllers.ReturnBookResponse}
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /transactions/return [put]
func (tc TransactionController) Return(c *gin.Context) {
	var returnRequest ReturnBookRequest
	if err := c.ShouldBindJSON(&returnRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	denda, err := transactionModel.Return(returnRequest.TransactionID, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, types.StandardResponse{
		Code:    http.StatusOK,
		Message: "Book returned successfully",
		Data:    gin.H{"denda": denda}, // Include the denda in the response
	})
}

//GetTransactionsByMemberID gets all the transactions by member ID
// @Summary Get Transactions by Member
// @Description Retrieve transactions for a specific member.
// @Tags Transaction
// @Produce json
// @Param id path string true "Member ID (UUID)"
// @Success 200 {object} types.StandardResponse{data=models.TransactionSwagger}
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /user/{id}/transactions [get]
func (tc TransactionController) GetTransactionsByMemberID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Invalid member ID"})
		return
	}
	transactions, err := transactionModel.GetTransactionsByMemberID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	transactionsSwagger := make([]models.TransactionSwagger, len(transactions))
    for i, t := range transactions {
        transactionsSwagger[i] = models.TransactionSwagger{
            ID:               t.ID,
            IDAnggota:        t.IDAnggota,
            IDBuku:           t.IDBuku,
            TanggalPinjam:    t.TanggalPinjam,
            TanggalJatuhTempo: t.TanggalJatuhTempo,
			TanggalKembali:   t.TanggalKembali.Time.Format(time.RFC3339), // Or t.TanggalKembali.String,
            Denda:            t.Denda.Int64, 							// Or t.Denda.String,
			Status:			t.Status,
        }
    }
	c.JSON(http.StatusOK, types.StandardResponse{
		Code:    http.StatusOK,
		Message: "Successfully retrieved transactions",
		Data:    transactionsSwagger,
	})
}

