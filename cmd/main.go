package main

import (
	"context"
	"go-stock/backend/data"
	"go-stock/backend/db"
	"go-stock/backend/logger"
	"go-stock/backend/models"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func AutoMigrate() {
	db.Dao.AutoMigrate(&data.StockInfo{})
	db.Dao.AutoMigrate(&data.StockBasic{})
	db.Dao.AutoMigrate(&data.FollowedStock{})
	db.Dao.AutoMigrate(&data.IndexBasic{})
	db.Dao.AutoMigrate(&data.Settings{})
	db.Dao.AutoMigrate(&models.AIResponseResult{})
	db.Dao.AutoMigrate(&models.StockInfoHK{})
	db.Dao.AutoMigrate(&models.StockInfoUS{})
	db.Dao.AutoMigrate(&data.FollowedFund{})
	db.Dao.AutoMigrate(&data.FundBasic{})
	db.Dao.AutoMigrate(&models.PromptTemplate{})
	db.Dao.AutoMigrate(&data.Group{})
	db.Dao.AutoMigrate(&data.GroupStock{})
	db.Dao.AutoMigrate(&models.Tags{})
	db.Dao.AutoMigrate(&models.Telegraph{})
	db.Dao.AutoMigrate(&models.TelegraphTags{})
	db.Dao.AutoMigrate(&models.LongTigerRankData{})
	db.Dao.AutoMigrate(&data.AIConfig{})
	db.Dao.AutoMigrate(&models.BKDict{})
	db.Dao.AutoMigrate(&models.WordAnalyze{})
	db.Dao.AutoMigrate(&models.SentimentResultAnalyze{})
	db.Dao.AutoMigrate(&models.AiRecommendStocks{})

	//updateMultipleModel()
}

func main() {
	// 初始化日志
	// 注意：这里直接使用 logger，不调用 Init 函数
	logger.SugaredLogger.Info("Backend service starting...")

	// 初始化数据库
	db.Init("data/stock.db")
	logger.SugaredLogger.Info("Database initialized successfully")

	data.InitAnalyzeSentiment()
	go AutoMigrate()

	// 初始化数据API
	stockDataApi := data.NewStockDataApi()
	// 示例：获取股票基础信息
	logger.SugaredLogger.Info("Fetching stock base info...")
	stockDataApi.GetStockBaseInfo()

	// 示例：获取实时股票数据
	logger.SugaredLogger.Info("Fetching real-time stock data...")
	stocks := []string{"sh600000", "sz000001"}
	stockInfos, err := stockDataApi.GetStockCodeRealTimeData(stocks...)
	if err != nil {
		logger.SugaredLogger.Error("Error fetching stock data:", err)
	} else {
		logger.SugaredLogger.Infof("Fetched data for %d stocks", len(*stockInfos))
		for _, stock := range *stockInfos {
			logger.SugaredLogger.Infof("Stock: %s - %s, Price: %s",
				stock.Code, stock.Name, stock.Price)
		}
	}

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.SugaredLogger.Info("Shutting down backend service...")

	// 优雅退出
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.SugaredLogger.Info("Backend service stopped")
}
