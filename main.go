package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Модель товара
type Product struct {
    ID      uint    `json:"id"`      // Уникальный идентификатор товара
    Product string  `json:"product"` // Название товара
    Price   float64 `json:"price"`   // Цена товара
}

var db *gorm.DB
var err error

// Инициализация базы данных
func initDatabase() {
    // Подключение к базе данных
    db, err = gorm.Open("sqlite3", "./store.db")
    if err != nil {
        log.Fatalf("Ошибка подключения к базе данных: %s", err)
    }
    fmt.Println("Успешное подключение к базе данных!")

    // Автоматическое создание таблицы для модели Product
    db.AutoMigrate(&Product{})
}

// Получение списка всех товаров
func getProducts(c *gin.Context) {
    var products []Product
    if err := db.Find(&products).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить список товаров"})
        return
    }
    c.JSON(http.StatusOK, products)
}

// Добавление нового товара
func addProduct(c *gin.Context) {
    var newProduct Product
    if err := c.BindJSON(&newProduct); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
        return
    }
    if err := db.Create(&newProduct).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось добавить товар"})
        return
    }
    c.JSON(http.StatusCreated, newProduct)
}

// Обновление информации о товаре
func updateProduct(c *gin.Context) {
    var product Product
    id := c.Param("id")

    if err := db.First(&product, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Товар не найден"})
        return
    }

    var updatedProduct Product
    if err := c.BindJSON(&updatedProduct); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
        return
    }

    product.Product = updatedProduct.Product
    product.Price = updatedProduct.Price

    if err := db.Save(&product).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить товар"})
        return
    }
    c.JSON(http.StatusOK, product)
}

// Удаление товара
func deleteProduct(c *gin.Context) {
    var product Product
    id := c.Param("id")

    if err := db.First(&product, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Товар не найден"})
        return
    }

    if err := db.Delete(&product).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить товар"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Товар удален"})
}

func main() {
    initDatabase()
    defer func() {
        if err := db.Close(); err != nil {
            log.Fatalf("Ошибка при закрытии базы данных: %s", err)
        }
    }()

    r := gin.Default()

    // Маршруты
    r.GET("/products", getProducts)           // Получение списка всех товаров
    r.POST("/products", addProduct)           // Добавление нового товара
    r.PUT("/products/:id", updateProduct)     // Обновление товара
    r.DELETE("/products/:id", deleteProduct)  // Удаление товара

    // Запуск сервера на порту 8080
    r.Run(":8080")
}
