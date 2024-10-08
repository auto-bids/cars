package controllers

import (
	"cars/models"
	"cars/responses"
	"cars/service"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

// DeleteAllUserOffer godoc
// @Summary Delete all offers
// @Description Delete all offers by user email
// @ID delete-all-offers
// @Produce json
// @Param email path string true "Email address of the user to delete all offers"
// @Success 200
// @Failure 400 {object} responses.UserResponse
// @Failure 404 {object} responses.UserResponse
// @Failure 500 {object} responses.UserResponse
// @Router /cars/delete/all/{email} [delete]
func DeleteAllUserOffer(c *gin.Context) {
	result := make(chan responses.UserResponse)

	go func(cCp *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		defer close(result)
		validate := validator.New(validator.WithRequiredStructEnabled())
		email := models.Email{Email: cCp.Param("email")}

		if err := validate.Struct(email); err != nil {
			result <- responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error validation email",
				Data:    map[string]interface{}{"error": err.Error()},
			}
			return
		}

		var userCollection = service.GetCollection(service.DB)

		filter := bson.D{{"user_email", email.Email}}
		results, err := userCollection.DeleteMany(ctx, filter)
		if err != nil {
			result <- responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error finding offers",
				Data:    map[string]interface{}{"error": err.Error()},
			}
			return
		}
		if results.DeletedCount < 1 {
			result <- responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "Offer does not exist",
				Data:    map[string]interface{}{"error": results},
			}
			return
		}
		result <- responses.UserResponse{
			Status:  http.StatusOK,
			Message: "ok",
			Data:    map[string]interface{}{"data": results},
		}
	}(c.Copy())
	res := <-result
	c.JSON(res.Status, res)
}
