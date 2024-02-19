package routes

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"millidatainsert/db"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
)

type Product struct {
	ProductId    string `json:"productid,omitempty" binding:"required"`
	Product      string `json:"product,omitempty" binding:"required"`
	Price        int64  `json:"price,omitempty"`
	Shop         string `json:"shop,omitempty" binding:"required"`
	Location     string `json:"location,omitempty" binding:"required"`
	Type         string `json:"type,omitempty" binding:"required"`
	Country      string `json:"country,omitempty" binding:"required"`
	Registeredat string `json:"registeredat,omitempty" binding:"required"`
}

func GenerateCsv(c *gin.Context) {

	start := time.Now()

	db1, err := db.InitDb()
	if err != nil {
		log.Fatal(err.Error())
	}

	myfake := make([]Product, 5)
	for i := 0; i < 5; i++ {
		myfake[i] = Product{}
		err := faker.FakeData(&myfake[i])

		myfake[i].Registeredat = time.Now().UTC().Format("2006-01-02T15:04:05-0700")
		if err != nil {
			fmt.Println(err)
		}
	}

	products, err := GetAllProducts(db1)
	if err != nil {
		log.Println("error")
	}

	tookSecond := time.Since(start)

	log.Printf("Done in %d seconds", int(math.Ceil(tookSecond.Seconds())))

	c.JSON(http.StatusAccepted, gin.H{"products": products})
}

func GetAllProducts(db *sql.DB) (*[]Product, error) {
	query := "SELECT * FROM products"

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err, "error")
		return nil, err
	}

	fmt.Println(rows, "the rows")

	defer rows.Close()

	var products []Product

	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ProductId, &product.Country, &product.Location, &product.Price, &product.Product, &product.Registeredat, &product.Shop, &product.Type)

		fmt.Println(product, "the product")
		if err != nil {
			fmt.Println(err, "the error")
			return nil, err
		}
		products = append(products, product)
	}

	return &products, nil
}
