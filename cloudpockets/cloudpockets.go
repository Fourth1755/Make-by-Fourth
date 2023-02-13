package cloudpockets

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CloudPocket struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
	Account  string  `json:"account"`
}
type handler struct {
	db *sql.DB
}

func NewApplication(db *sql.DB) *handler {
	return &handler{db}
}
func (h *handler) CreateCloudPockets(c echo.Context) error {
	var cp CloudPocket
	err := c.Bind(&cp)
	if err != nil {
		//logger.Error("bad request body", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body", err.Error())
	}
	row := h.db.QueryRow("INSERT INTO cloud_pockets (ID, Name, category, Currency, Balance , Account) values ($1, $2, $3, $4, $5 , $6)  RETURNING id;", cp.ID, cp.Name, cp.Category, cp.Currency, cp.Balance, cp.Account)
	err = row.Scan(&cp.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "can't craete could pockets", err.Error())
	}
	return c.JSON(http.StatusCreated, cp)
}
func (h *handler) GetAllCloudPockets(c echo.Context) error {
	//logger := mlog.L(c)

	var cp CloudPocket
	err := c.Bind(&cp)
	if err != nil {
		//logger.Error("bad request body", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body", err.Error())
	}

	rows, err := h.db.Query("SELECT * FROM cloud_pockets")
	if err != nil {
		//logger.Error("error", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "can't get all could pockets", err.Error())
	}
	var cloudPockets = []CloudPocket{}
	for rows.Next() {
		var c CloudPocket
		err := rows.Scan(&c.ID, &c.Name, &c.Category, &c.Currency, &c.Balance, &c.Account)
		if err != nil {
			//logger.Error("can't scan query all cloud_pockets", zap.Error(err))
			return echo.NewHTTPError(http.StatusInternalServerError, "can't get all could pockets", err.Error())
		}
		c = CloudPocket{
			ID: c.ID, Name: c.Name, Category: c.Category, Currency: c.Currency, Balance: c.Balance, Account: c.Account,
		}
		cloudPockets = append(cloudPockets, c)
	}

	return c.JSON(http.StatusOK, cloudPockets)
}
