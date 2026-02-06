package handlers

import (
	"e-wallet-go/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	walletService services.WalletService
}

func NewWalletHandler(walletService services.WalletService) *WalletHandler {
	return &WalletHandler{walletService: walletService}
}

type AmountRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

func (h *WalletHandler) GetMyBalance(c *gin.Context) {
	userID := c.GetString("userID")

	wallet, err := h.walletService.GetUserBalance(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch balance"})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

func (h *WalletHandler) UserTopUp(c *gin.Context) {
	userID := c.GetString("userID")

	var req AmountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := h.walletService.UserTopUp(userID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Top up successful",
		"user_id":       wallet.UserID,
		"wallet_number": wallet.WalletNumber,
		"balance":       wallet.Balance,
	})
}

func (h *WalletHandler) UserWithdraw(c *gin.Context) {
	userID := c.GetString("userID")

	var req AmountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := h.walletService.UserWithdraw(userID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Withdraw successful",
		"user_id":       wallet.UserID,
		"wallet_number": wallet.WalletNumber,
		"balance":       wallet.Balance,
	})
}

func (h *WalletHandler) GetTransactionHistory(c *gin.Context) {
	userID := c.GetString("userID")

	searchQuery := c.Query("search")

	transactions, err := h.walletService.GetTransactionHistory(userID, searchQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": transactions,
	})
}

func (h *WalletHandler) GetAdminGlobalTransactions(c *gin.Context) {
	searchQuery := c.Query("search")

	transactions, err := h.walletService.GetAdminGlobalTransactions(searchQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success fetching global transactions",
		"total":   len(transactions),
		"data":    transactions,
	})
}

func (h *WalletHandler) AdminTopUpUser(c *gin.Context) {
	targetUserID := c.Param("userID")

	var req AmountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := h.walletService.AdminTopUpUser(targetUserID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Admin Top-up successful",
		"target_user_id": wallet.UserID,
		"new_balance":    wallet.Balance,
	})
}

func (h *WalletHandler) AdminDeductUser(c *gin.Context) {
	targetUserID := c.Param("userID")

	var req AmountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet, err := h.walletService.AdminDeductUser(targetUserID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Admin Deduction successful",
		"target_user_id": wallet.UserID,
		"new_balance":    wallet.Balance,
	})
}
