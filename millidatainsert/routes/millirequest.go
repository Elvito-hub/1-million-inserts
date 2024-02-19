package routes

import (
	"context"
	"fmt"
	"log"
	"millidatainsert/db"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func HandleMilliRequest(c *gin.Context) {
	var newProduct Product
	err := c.ShouldBindJSON(&newProduct)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println(newProduct, "new product")
	InsertData(newProduct)

	c.JSON(http.StatusCreated, gin.H{})

}

func InsertData(p Product) {

	productId := p.ProductId
	product := p.Product
	shop := p.Shop
	location := p.Location
	Type := p.Type
	price := p.Price
	country := p.Country
	registeredAt := p.Registeredat
	conn, err := db.DB.Conn(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}

	query := fmt.Sprintf("INSERT INTO products (%s) VALUES (%s)",
		strings.Join(DataHeaders, ","),
		strings.Join(generateQuestionsMark(len(DataHeaders)), ","),
	)
	// values = trimSpaces(values)

	_, err = conn.ExecContext(context.Background(), query, productId, product, shop, location, Type, price, country, registeredAt)

	if err != nil {
		log.Fatal(err.Error())
	}
	defer conn.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
}
