package routes

import (
	"fmt"
	"net/http"
	"time"

	"delineate.io/customers/src/database"
	"delineate.io/customers/src/logging"
	"delineate.io/customers/src/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

type Customer struct {
	ID          int       `json:"id"`
	Forename    string    `json:"forename" binding:"required"`
	Surname     string    `json:"surname" binding:"required"`
	DateOfBirth time.Time `json:"dob" binding:"required"`
}

func newCustomer(item *models.Customer) Customer {
	return Customer{
		ID:          item.ID,
		Forename:    item.Forename,
		Surname:     item.Surname,
		DateOfBirth: item.DateOfBirth,
	}
}

// createCustomer godoc
// @Summary Create Customer
// @Description Creates a new customer
// @Accept json
// @param customer body Customer true "customer"
// @Produce  plain
// @Success 201 "Customer created"
// @Failure 400 "Bad request"
// @Failure 409 "Customer exists"
// @Failure 500 "Internal error"
// @Router /customer [post]
func createCustomer(ctx *gin.Context) {
	var request Customer
	if err := ctx.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		logging.Err(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if db, err := database.OpenDB(); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	} else {
		if customerExists(request, db) {
			logging.Warn("customer already exists")
			ctx.AbortWithStatus(http.StatusConflict)
			return
		}
		customer := *models.NewCustomer(request.Forename,
			request.Surname,
			request.DateOfBirth)

		if err := db.Create(&customer).Error; err != nil {
			logging.Err(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		ctx.AbortWithStatus(http.StatusCreated)
	}
}

func getCondition(request Customer) (search string, values []interface{}) {
	condition := "forename = ? AND surname = ? AND dob = ?"
	var params []interface{}
	params = append(params, request.Forename, request.Surname, request.DateOfBirth)
	return condition, params
}

func customerExists(request Customer, db *gorm.DB) bool {
	var customer models.Customer
	condition, params := getCondition(request)
	db.Where(condition, params...).Find(&customer)
	return customer.ID != 0
}

// getCustomer godoc
// @Summary Get Customer
// @Description Returns a single customer
// @Accept plain
// @Produce json
// @Param id path int true "Customer ID"
// @Success 200 {object} routes.Customer "Customer returned"
// @Failure 404 "Customer not found"
// @Failure 500 "Internal error"
// @Router /customer/{id} [get]
func getCustomerByID(ctx *gin.Context) {
	id, _ := ctx.Params.Get("id")
	var customer models.Customer
	if db, err := database.OpenDB(); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	} else {
		db.Limit(1).Find(&customer, id)
		if customer.ID == 0 {
			logging.Warn(fmt.Sprintf("customer '%s' not found", id))
			ctx.AbortWithStatus(http.StatusNotFound)
		} else {
			ctx.SecureJSON(http.StatusOK, newCustomer(&customer))
		}
	}
}
