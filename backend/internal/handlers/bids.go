package handlers

import (
	"net/http"
	"tender_service/internal/models"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateBid(c *gin.Context) {
	var bid models.Bid
	if err := c.ShouldBindJSON(&bid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	bid.ID = uuid.New()
	if err := db.Table("bid").Create(&bid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать предложение"})
		return
	}

	c.JSON(http.StatusOK, bid)
}

func GetUserBids(c *gin.Context) {
	username := c.Query("username")
	var bids []models.Bid

	if err := db.Table("bid").Where("author_id = (SELECT id FROM employee WHERE username = ?)", username).Find(&bids).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не существует или некорректен."})
		return
	}

	c.JSON(http.StatusOK, bids)
}

func GetBidsForTender(c *gin.Context) {
	tenderId := c.Param("tenderId")
	username := c.Query("username")
	var bids []models.Bid

	if err := db.Table("bid").Where("tender_id = ? AND author_id = (SELECT id FROM employee WHERE username = ?)", tenderId, username).Find(&bids).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Тендер или предложение не найдены."})
		return
	}

	c.JSON(http.StatusOK, bids)
}

func GetBidStatus(c *gin.Context) {
	bidId := c.Param("bidId")
	username := c.Query("username")
	var bid models.Bid

	if err := db.Table("bid").Where("id = ? AND author_id = (SELECT id FROM employee WHERE username = ?)", bidId, username).First(&bid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Предложение не найдено."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": bid.Status})
}

func UpdateBidStatus(c *gin.Context) {
	bidId := c.Param("bidId")
	status := c.Query("status")
	username := c.Query("username")

	if err := db.Table("bid").Where("id = ? AND author_id = (SELECT id FROM employee WHERE username = ?)", bidId, username).Updates(models.Bid{Status: status}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Предложение не найдено."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Статус предложения успешно изменен."})
}

func EditBid(c *gin.Context) {
	bidId := c.Param("bidId")
	username := c.Query("username")
	var bid models.Bid

	// Проверка существования предложения
	if err := db.Table("bid").Where("id = ? AND author_id = (SELECT id FROM employee WHERE username = ?)", bidId, username).First(&bid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Предложение не найдено."})
		return
	}

	// Обновление полей предложения
	var updatedBid models.Bid
	if err := c.ShouldBindJSON(&updatedBid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	// Обновление только тех полей, которые были переданы
	if updatedBid.Name != "" {
		bid.Name = updatedBid.Name
	}
	if updatedBid.Description != "" {
		bid.Description = updatedBid.Description
	}

	if err := db.Save(&bid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить предложение"})
		return
	}

	c.JSON(http.StatusOK, bid)
}

func SubmitBidDecision(c *gin.Context) {
	bidId := c.Param("bidId")
	decision := c.Query("decision")
	username := c.Query("username")

	if decision == "" || username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Имя пользователя и решение обязательны"})
		return
	}

	var user models.User
	if err := db.Table("user").Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не существует или некорректен"})
		return
	}

	var bid models.Bid
	if err := db.Table("bid").Where("id = ?", bidId).First(&bid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Предложение не найдено"})
		return
	}

	if decision == "Approved" {
		bid.Status = "Approved"
	} else if decision == "Rejected" {
		bid.Status = "Rejected"
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверное значение решения"})
		return
	}

	if err := db.Save(&bid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить предложение"})
		return
	}

	c.JSON(http.StatusOK, bid)
}

func SubmitBidFeedback(c *gin.Context) {
	bidId := c.Param("bidId")
	feedback := c.Query("bidFeedback")
	username := c.Query("username")

	if feedback == "" || username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Имя пользователя и отзыв обязательны"})
		return
	}

	var user models.User
	if err := db.Table("user").Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не существует или некорректен"})
		return
	}

	var bid models.Bid
	if err := db.Table("bid").Where("id = ?", bidId).First(&bid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Предложение не найдено"})
		return
	}

	// Создание отзыва
	var bidReview models.BidReview
	bidReview.ID = uuid.New()
	bidReview.BidID = bid.ID
	bidReview.Description = feedback

	if err := db.Table("bid_review").Create(&bidReview).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось отправить отзыв"})
		return
	}

	c.JSON(http.StatusOK, bid)
}

func RollbackBid(c *gin.Context) {
	bidId := c.Param("bidId")
	versionStr := c.Param("version")
	username := c.Query("username")

	// Проверка существования пользователя
	var user models.User
	if err := db.Table("user").Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не существует или некорректен"})
		return
	}

	// Преобразование версии из строки в целое число
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат версии"})
		return
	}

	// Логика отката версии
	var bidHistory models.BidHistory
	if err := db.Table("bid_history").Where("bid_id = ? AND version = ?", bidId, version).First(&bidHistory).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Предложение или версия не найдены."})
		return
	}

	// Обновление текущего предложения
	var bid models.Bid
	if err := db.Table("bid").Where("id = ?", bidId).First(&bid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Предложение не найдено."})
		return
	}

	bid.Name = bidHistory.Name
	bid.Description = bidHistory.Description
	bid.Status = bidHistory.Status
	bid.TenderID = bidHistory.TenderID
	bid.AuthorType = bidHistory.AuthorType
	bid.AuthorID = bidHistory.AuthorID
	bid.Version++

	if err := db.Save(&bid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось откатить предложение"})
		return
	}

	c.JSON(http.StatusOK, bid)
}

func GetBidReviews(c *gin.Context) {
	tenderId := c.Param("tenderId")
	authorUsername := c.Query("authorUsername")
	requesterUsername := c.Query("requesterUsername")

	// Проверка существования пользователя
	var requester models.User
	if err := db.Table("user").Where("username = ?", requesterUsername).First(&requester).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не существует или некорректен"})
		return
	}

	// Получение отзывов
	var reviews []models.BidReview
	if err := db.Table("bid_review").Where("bid_id IN (SELECT id FROM bid WHERE tender_id = ? AND author_id = (SELECT id FROM employee WHERE username = ?))", tenderId, authorUsername).Find(&reviews).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Тендер или отзывы не найдены."})
		return
	}

	c.JSON(http.StatusOK, reviews)
}
