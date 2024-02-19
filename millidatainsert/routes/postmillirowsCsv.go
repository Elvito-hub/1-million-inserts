package routes

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"millidatainsert/db"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var DataHeaders = []string{"productid", "product", "shop", "location", "type", "price", "country", "registeredat"}

var totalWorker = 150

func HandlePostMilliData(c *gin.Context) {
	start := time.Now()
	db1, err := db.InitDb()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Get the CSV file from the POST request
	file, err := c.FormFile("csv")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Open the file
	f, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer f.Close()

	reader := csv.NewReader(f)

	chunks, rows := divideWorkers(reader)

	DispatchWorkers(chunks, *rows, db1)

	duration := time.Since(start)
	log.Printf("Done in %d seconds, to insert %d rows", int(math.Ceil(duration.Seconds())), len(*rows))

	c.JSON(200, gin.H{"message": "CSV file received successfully"})
}

func processCsvFile(reader *csv.Reader, jobs chan<- []interface{}, wg *sync.WaitGroup) {
	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}

		if len(DataHeaders) == 0 {
			fmt.Println(row, "all row")
			DataHeaders = row
			continue
		}

		rowOrdered := make([]interface{}, len(row))
		for i, v := range row {
			rowOrdered[i] = v
		}

		wg.Add(1)

		jobs <- rowOrdered
	}
	close(jobs)

}

func dispatchWorkers(db *sql.DB, jobs <-chan []interface{}, wg *sync.WaitGroup) {
	for workerIndex := 0; workerIndex <= totalWorker; workerIndex++ {
		go func(wIndex int, db *sql.DB, jobs <-chan []interface{}, wg *sync.WaitGroup) {
			counter := 0
			for job := range jobs {
				PerformInsert(wIndex, counter, db, job)
				wg.Done()
				counter++
			}
		}(workerIndex, db, jobs, wg)
	}
}

func PerformInsert(workerIndex, counter int, db *sql.DB, values []interface{}) {
	// for {
	productId := values[0]
	product := values[1]
	shop := values[2]
	location := values[3]
	Type := values[4]
	price := 0
	n, ok := values[5].(int)
	if ok {
		price = n
	}
	country := values[6]
	registeredAt := values[7]
	conn, err := db.Conn(context.Background())
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

	if counter%100 == 0 {
		log.Printf("=> Worker %d inserted %d data\n", workerIndex, counter)
	}

}
func generateQuestionsMark(n int) []string {
	marks := make([]string, n)
	for i := 1; i <= n; i++ {
		marks[i-1] = "$" + fmt.Sprint(i)
	}
	return marks
}

func trimSpaces(vals []interface{}) []interface{} {
	result := make([]interface{}, len(vals))
	for i, v := range vals {
		str, ok := v.(string)
		if ok {
			str = strings.TrimSpace(str)
		}
		result[i] = str
	}
	return result
}

func divideWorkers(reader *csv.Reader) ([]int, *[][]interface{}) {
	nChunks := 50

	var rows [][]interface{}

	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}

		rowOrdered := make([]interface{}, len(row))
		for i, v := range row {
			rowOrdered[i] = v
		}
		rows = append(rows, rowOrdered)
	}

	chunkSize := len(rows) / nChunks

	chunks := make([]int, 0, nChunks)

	offset := 0

	for offset < len(rows) {
		offset += chunkSize
		if offset >= len(rows) {
			chunks = append(chunks, len(rows))
			break
		} else {
			chunks = append(chunks, offset)
		}
	}

	return chunks, &rows

}

func DispatchWorkers(chunks []int, rows [][]interface{}, db *sql.DB) {
	var wg sync.WaitGroup
	wg.Add(len(chunks))
	for i := 0; i < len(chunks); i++ {
		var startInd int
		var endInd int
		if i == 0 {
			startInd = 0
			endInd = chunks[i]
		} else {
			startInd = chunks[i-1]
			endInd = chunks[i]
		}
		go func(workerInd int, rows [][]interface{}, startInd int, endInd int) {
			//var wg1 sync.WaitGroup

			for j := startInd; j < endInd; j++ {
				// wg1.Add(1)
				// go func(wg *sync.WaitGroup, j int) {
				PerformInsert(workerInd, j, db, rows[j])
				// 	wg1.Done()
				// }(&wg1, j)
			}
			//wg1.Wait()
			wg.Done()
		}(i, rows, startInd, endInd)
	}
	wg.Wait()
}
