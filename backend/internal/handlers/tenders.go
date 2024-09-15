package handlers

import (
	"net/http"
	"strconv"
	"time"

	"tender_service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

// Функция для инициализации db
func InitDB(database *gorm.DB) {
	db = database
}

func GetTenders(c *gin.Context) {
	var tenders []models.Tender
	serviceTypes := c.QueryArray("service_type")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit < 0 || limit > 50 {
		c.JSON(http.StatusBadRequest, gin.H{"reason": "Неверное значение лимита"})
		return
	}

	if offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"reason": "Неверное значение смещения"})
		return
	}

	query := db.Table("tender").Order("name")

	if len(serviceTypes) > 0 {
		query = query.Where("service_type IN (?)", serviceTypes)
	}

	var total int64
	query.Table("tender").Count(&total)

	if err := query.Table("tender").Limit(limit).Offset(offset).Find(&tenders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"reason": "Не удалось получить тендеры"})
		return
	}

	c.JSON(http.StatusOK, tenders)
}

func CreateTender(c *gin.Context) {
	var tender models.Tender
	if err := c.ShouldBindJSON(&tender); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка существования пользователя
	var user models.User
	if err := db.Table("user").Where("username = ?", tender.CreatorUsername).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не существует или некорректен"})
		return
	}

	// Проверка прав пользователя
	var orgResponsible models.OrganizationResponsible
	if err := db.Table("organization_responsible").Where("organization_id = ? AND user_id = ?", tender.OrganizationID, user.ID).First(&orgResponsible).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Недостаточно прав для выполнения действия"})
		return
	}

	tender.ID = uuid.New()
	tender.CreatedAt = time.Now()
	tender.UpdatedAt = time.Now()
	tender.Status = "Created"

	if err := db.Table("tender").Create(&tender).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tender)
}

func GetUserTenders(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"reason": "Имя пользователя обязательно"})
		return
	}

	var user models.User
	if err := db.Table("user").Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"reason": "Пользователь не существует или некорректен"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit < 0 || limit > 50 {
		c.JSON(http.StatusBadRequest, gin.H{"reason": "Неверное значение лимита"})
		return
	}

	if offset < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"reason": "Неверное значение смещения"})
		return
	}

	var tenders []models.Tender
	query := db.Table("tender").Where("creator_username = ?", username).Order("name")

	var total int64
	query.Count(&total)

	if err := query.Limit(limit).Offset(offset).Find(&tenders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"reason": "Не удалось получить тендеры пользователя"})
		return
	}

	c.JSON(http.StatusOK, tenders)
}

func GetTenderStatus(c *gin.Context) {
	tenderId := c.Param("tenderId")
	username := c.Query("username")

	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"reason": "Имя пользователя обязательно"})
		return
	}

	var user models.User
	if err := db.Table("user").Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"reason": "Пользователь не существует или некорректен"})
		return
	}

	var tender models.Tender
	if err := db.Table("tender").Where("id = ?", tenderId).First(&tender).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"reason": "Тендер не найден"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"reason": "Не удалось получить тендер"})
		}
		return
	}

	var orgResponsible models.OrganizationResponsible
	if err := db.Table("organization_responsible").Where("organization_id = ? AND user_id = ?", tender.OrganizationID, user.ID).First(&orgResponsible).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"reason": "Недостаточно прав для выполнения действия"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": tender.Status})
}

func UpdateTenderStatus(c *gin.Context) {
	tenderId := c.Param("tenderId")
	newStatus := c.Query("status")
	username := c.Query("username")

	if username == "" || newStatus == "" {
		c.JSON(http.StatusBadRequest, gin.H{"reason": "Имя пользователя и статус обязательны"})
		return
	}

	var user models.User
	if err := db.Table("user").Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"reason": "Пользователь не существует или некорректен"})
		return
	}

	var tender models.Tender
	if err := db.Table("tender").Where("id = ?", tenderId).First(&tender).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"reason": "Тендер не найден"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"reason": "Не удалось получить тендер"})
		}
		return
	}

	var orgResponsible models.OrganizationResponsible
	if err := db.Table("organization_responsible").Where("organization_id = ? AND user_id = ?", tender.OrganizationID, user.ID).First(&orgResponsible).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"reason": "Недостаточно прав для выполнения действия"})
		return
	}

	if newStatus != "Created" && newStatus != "Published" && newStatus != "Closed" {
		c.JSON(http.StatusBadRequest, gin.H{"reason": "Неверное значение статуса"})
		return
	}

	tender.Status = newStatus
	tender.UpdatedAt = time.Now()
	tender.Version++

	if err := db.Table("tender").Save(&tender).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"reason": "Не удалось обновить статус тендера"})
		return
	}

	c.JSON(http.StatusOK, tender)
}

func EditTender(c *gin.Context) {
	tenderId := c.Param("tenderId")
	username := c.Query("username")

	var updateData struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		ServiceType string `json:"serviceType"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"reason": "Данные некорректны"})
		return
	}

	var user models.User
	if err := db.Table("user").Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"reason": "Имя пользователя некорректно"})
		return
	}

	var tender models.Tender
	if err := db.Table("tender").Where("id = ?", tenderId).First(&tender).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"reason": "Тендер не найден"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"reason": "Не удалось получить тендер"})
		}
		return
	}

	var orgResponsible models.OrganizationResponsible
	if err := db.Table("organization_responsible").Where("organization_id = ? AND user_id = ?", tender.OrganizationID, user.ID).First(&orgResponsible).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"reason": "Недостаточно прав для выполнения действия"})
		return
	}

	tenderHistory := models.TenderHistory{
		TenderID:        tender.ID,
		Name:            tender.Name,
		Description:     tender.Description,
		ServiceType:     tender.ServiceType,
		Status:          tender.Status,
		OrganizationID:  tender.OrganizationID,
		CreatorUsername: tender.CreatorUsername,
		Version:         tender.Version,
		CreatedAt:       tender.CreatedAt,
		UpdatedAt:       tender.UpdatedAt,
	}

	if err := db.Table("tender_history").Create(&tenderHistory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"reason": "Не удалось создать историю тендера"})
		return
	}

	if updateData.Name != "" {
		tender.Name = updateData.Name
	}
	if updateData.Description != "" {
		tender.Description = updateData.Description
	}
	if updateData.ServiceType != "" {
		tender.ServiceType = updateData.ServiceType
	}

	tender.Version++
	tender.UpdatedAt = time.Now()

	if err := db.Table("tender").Save(&tender).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"reason": "Не удалось обновить тендер"})
		return
	}

	c.JSON(http.StatusOK, tender)
}

func RollbackTender(c *gin.Context) {
	tenderId := c.Param("tenderId")
	version, err := strconv.Atoi(c.Param("version"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"reason": "Неверный формат версии"})
		return
	}
	username := c.Query("username")

	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"reason": "Имя пользователя обязательно"})
		return
	}

	var user models.User
	if err := db.Table("user").Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"reason": "Пользователь не существует или некорректен"})
		return
	}

	var tender models.Tender
	if err := db.Table("tender").Where("id = ?", tenderId).First(&tender).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"reason": "Тендер не найден"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"reason": "Не удалось получить тендер"})
		}
		return
	}

	var orgResponsible models.OrganizationResponsible
	if err := db.Table("organization_responsible").Where("organization_id = ? AND user_id = ?", tender.OrganizationID, user.ID).First(&orgResponsible).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"reason": "Недостаточно прав для выполнения действия"})
		return
	}

	var tenderHistory models.TenderHistory
	if err := db.Table("tender_history").Where("tender_id = ? AND version = ?", tenderId, version).First(&tenderHistory).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{"reason": "Версия тендера не найдена"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"reason": "Не удалось получить историю тендера"})
		}
		return
	}

	tender.Name = tenderHistory.Name
	tender.Description = tenderHistory.Description
	tender.ServiceType = tenderHistory.ServiceType
	tender.Status = tenderHistory.Status
	tender.Version++
	tender.UpdatedAt = time.Now()

	if err := db.Table("tender").Save(&tender).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"reason": "Не удалось обновить тендер"})
		return
	}

	c.JSON(http.StatusOK, tender)
}
