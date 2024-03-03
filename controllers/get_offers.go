package controllers

import (
	"cars/models"
	"cars/queries"
	"cars/responses"
	"cars/service"
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetOffers(c *gin.Context) {
	result := make(chan responses.UserResponse)

	go func(cCp *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		defer close(result)
		var resultModel models.CheckOffer
		validate := validator.New(validator.WithRequiredStructEnabled())

		pageStr := cCp.Param("page")

		page, err := strconv.ParseInt(pageStr, 10, 64)
		if err != nil {
			result <- responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error parse page",
				Data:    map[string]interface{}{"error": err.Error()},
			}
			return
		}

		if err := cCp.ShouldBind(&resultModel); err != nil {
			result <- responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error model get offer",
				Data:    map[string]interface{}{"error": err.Error()},
			}
			return
		}

		fmt.Println(resultModel)

		if err := validate.Struct(resultModel); err != nil {
			result <- responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error validation offers",
				Data:    map[string]interface{}{"error": err.Error()},
			}
			return
		}

		limit := int64(10)
		var userCollection = service.GetCollection(service.DB, "cars")

		filter := queries.GetOfferQuery(resultModel)
		opts := options.Find().SetSkip((page - 1) * 10).SetLimit(limit)
		results, err := userCollection.Find(ctx, filter, opts)
		if err != nil {
			result <- responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error finding offers",
				Data:    map[string]interface{}{"error": err.Error()},
			}
			return
		}

		var offers []models.Offer
		if err := results.All(ctx, &offers); err != nil {
			result <- responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error decoding offers",
				Data:    map[string]interface{}{"error": err.Error()},
			}
			return
		}

		result <- responses.UserResponse{
			Status:  http.StatusOK,
			Message: "ok",
			Data:    map[string]interface{}{"data": offers, "number_of_pages": (len(offers) + 10 - 1) / 10},
		}

	}(c.Copy())
	res := <-result
	c.JSON(res.Status, res)
}
