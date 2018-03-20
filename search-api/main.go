package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	elastic "gopkg.in/olivere/elastic.v5"
)

type Document struct {
	Title     string                `json:"title"`
	CreatedAt time.Time             `json:"created_at"`
	Content   string                `json:"content"`
	Suggest   *elastic.SuggestField `json:"suggest_field"`
}

const mapping = `
{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0
  },
  "mappings": {
    "document" {
      "properties": {
        "title": {
          "type": "text"
        },
        "created_at": {
          "type": "date"
        },
        "content": {
          "type": "text",
          "store": true,
          "fielddata": true
        },
        "suggest_field": {
          "type": "completion"
        }
      }
    }
  }
}
`

type CreateDocuentRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type SearchRequest struct {
}

var (
	elasticClient *elastic.Client
)

func main() {
	var err error

	// Create Elastic client and wait for Elasticsearch to be ready
	for {
		elasticClient, err = elastic.NewClient(
			elastic.SetURL("http://elasticsearch:9200"),
			elastic.SetSniff(false),
		)
		if err != nil {
			log.Println(err)
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}

	r := gin.Default()
	r.POST("/documents", createDocumentEndpoint)
	r.GET("/search", searchEndpoint)
	if err = r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func createDocumentEndpoint(c *gin.Context) {
	var req CreateDocuentRequest
	if err := c.BindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "Malformed request body")
		return
	}
	c.JSON(http.StatusOK, req)
}

func searchEndpoint(c *gin.Context) {
}

func errorResponse(c *gin.Context, code int, err string) {
	c.JSON(code, gin.H{
		"error": err,
	})
}
