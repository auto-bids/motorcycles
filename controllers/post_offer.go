package controllers

import (
	"context"
	"github.com/go-playground/validator/v10"
	"motorcycles/models"
	"motorcycles/responses"
	"motorcycles/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func PostOffer(c *gin.Context) {
	result := make(chan responses.UserResponse)

	go func(cCp *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		defer close(result)
		var resultMotorcycle models.Motorcycle
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

		if err := cCp.ShouldBindJSON(&resultMotorcycle); err != nil {
			result <- responses.UserResponse{
				Status:  http.StatusBadRequest,
				Message: "Invalid request body",
				Data:    map[string]interface{}{"error": err.Error()},
			}
			return
		}

		if err := validate.Struct(resultMotorcycle); err != nil {
			result <- responses.UserResponse{
				Status:  http.StatusBadRequest,
				Message: "Error validation car",
				Data:    map[string]interface{}{"error": err.Error()},
			}
			return
		}

		var userCollection = service.GetCollection(service.DB)
		newOffer := models.PostOffer{
			UserEmail:  email.Email,
			Motorcycle: resultMotorcycle,
		}
		results, err := userCollection.InsertOne(ctx, newOffer)
		if err != nil {
			result <- responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error adding offer",
				Data:    map[string]interface{}{"error": err.Error()},
			}
			return
		}
		result <- responses.UserResponse{
			Status:  http.StatusCreated,
			Message: "Offer created",
			Data:    map[string]interface{}{"data": results},
		}
	}(c.Copy())
	res := <-result
	c.JSON(res.Status, res)
}
