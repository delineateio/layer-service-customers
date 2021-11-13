package routes

import (
	"net/http"

	"delineate.io/customers/src/database"
	"delineate.io/customers/src/models"
	"github.com/gin-gonic/gin"
)

type Customers struct {
	Customers []Customer `json:"customers"`
}

// healthz godoc
// @Summary Get Customers
// @Description gets a list of all the customers
// @Accept plain
// @Produce json
// @Success 200 {object} routes.Customers "Customers returned"
// @Failure 500 "Internal error"
// @Router /customers [get]
func getCustomers(ctx *gin.Context) {
	var customers []models.Customer
	if db, err := database.OpenDB(); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	} else {
		db.Find(&customers)
		ctx.SecureJSON(http.StatusOK, convertCustomers(customers))
	}
}

func convertCustomers(items []models.Customer) Customers {
	customers := make([]Customer, 0)
	for i := 0; i < len(items); i++ {
		customers = append(customers, newCustomer(&items[i]))
	}
	return Customers{
		Customers: customers,
	}
}
