package main

import (
	"log"
	"os"
	"time"

	"tender_service/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *gorm.DB

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки файла .env")
	}

	initDB()
	defer db.Close()

	r := gin.Default()

	r.GET("/api/ping", handlers.Ping)

	r.POST("/api/tenders/new", handlers.CreateTender)
	r.GET("/api/tenders", handlers.GetTenders)
	r.GET("/api/tenders/my", handlers.GetUserTenders)
	r.GET("/api/tenders/:tenderId/status", handlers.GetTenderStatus)
	r.PUT("/api/tenders/:tenderId/status", handlers.UpdateTenderStatus)
	r.PATCH("/api/tenders/:tenderId/edit", handlers.EditTender)
	r.PUT("/api/tenders/:tenderId/rollback/:version", handlers.RollbackTender)
	r.GET("/api/tenders/:tenderId/bids/reviews", handlers.GetBidReviews)
	r.GET("/api/tenders/:tenderId/bids/list", handlers.GetBidsForTender)

	r.POST("/api/bids/new", handlers.CreateBid)
	r.GET("/api/bids/my", handlers.GetUserBids)
	r.GET("/api/bids/:bidId/status", handlers.GetBidStatus)
	r.PUT("/api/bids/:bidId/status", handlers.UpdateBidStatus)
	r.PATCH("/api/bids/:bidId/edit", handlers.EditBid)
	r.PUT("/api/bids/:bidId/submit_decision", handlers.SubmitBidDecision)
	r.PUT("/api/bids/:bidId/feedback", handlers.SubmitBidFeedback)
	r.PUT("/api/bids/:bidId/rollback/:version", handlers.RollbackBid)

	r.Run(os.Getenv("SERVER_ADDRESS"))
}

func initDB() {
	var err error
	dsn := os.Getenv("POSTGRES_CONN")
	log.Printf("Попытка подключения с DSN: %s", dsn)

	for i := 0; i < 30; i++ {
		db, err = gorm.Open("postgres", dsn)
		if err == nil {
			break
		}
		log.Printf("Не удалось подключиться к базе данных: %v", err)
		log.Printf("Повторная попытка через 5 секунд...")
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatal("Не удалось подключиться к базе данных после 30 попыток:", err)
	}

	if err := db.DB().Ping(); err != nil {
		log.Fatal("Не удалось выполнить ping к базе данных:", err)
	}

	log.Println("Успешно подключено к базе данных")
	handlers.InitDB(db)
}
