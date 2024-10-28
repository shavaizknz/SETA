package controller

import (
	"fmt"
	"seta/pkg/model"
	"seta/pkg/service"

	"github.com/labstack/echo/v4"
)

type TransactionController struct {
	TransactionService service.ITransactionService
}

func TransactionControllerProvider(transactionService service.ITransactionService) model.IController {
	return &TransactionController{TransactionService: transactionService}
}

func (tc *TransactionController) SetupRoutes(r *echo.Group) {
	r.POST("/deposit", tc.CreateDeposit)
	r.POST("/withdraw", tc.CreateWithdraw)
	r.PUT("/transaction", tc.UpdateTransaction)
	r.GET("/transaction/:transaction_id", tc.GetTransaction)
}

//------------------Controller Methods------------------//

// @BasePath /
// Create Deposits POST
// @Summary API To create a deposit transaction
// @Schemes
// @Description Api will return status 200 if the transaction is successful, 400 if the request is invalid and 500 if there is an internal server error
// @Tags Transaction
// @Accept json
// @Produce json
// @Success 200 {object} model.TransactionResponse
// @Failure 400 {object} model.DefaultError{error=string}
// @Failure 500 {object} model.DefaultError{error=string}
// @Param TransactionRequest body DepositRequest true "Transaction Request"
// @Router /api/v1/deposit [post]
func (tc *TransactionController) CreateDeposit(c echo.Context) error {
	params, err := tc.ValidateTransactionRequest(c)
	if err != nil {
		return c.JSON(400, model.DefaultError{Error: err.Error()})
	}

	transactionResponse, err := tc.TransactionService.CreateTransaction(c.Request().Context(), params.AccountID, params.Amount, model.TransactionTypeDeposit)
	if err != nil {
		return c.JSON(500, model.DefaultError{Error: err.Error()})
	}

	return c.JSON(200, transactionResponse)
}

// @BasePath /
// Create Withdraw POST
// @Summary API To create a withdraw transaction
// @Schemes
// @Description Api will return status 200 if the transaction is successful, 400 if the request is invalid and 500 if there is an internal server error
// @Tags Transaction
// @Accept json
// @Produce json
// @Success 200 {object} model.TransactionResponse
// @Failure 400 {object} model.DefaultError{error=string}
// @Failure 500 {object} model.DefaultError{error=string}
// @Param TransactionRequest body DepositRequest true "Transaction Request"
// @Router /api/v1/withdraw [post]
func (tc *TransactionController) CreateWithdraw(c echo.Context) error {
	params, err := tc.ValidateTransactionRequest(c)
	if err != nil {
		return c.JSON(400, model.DefaultError{Error: err.Error()})
	}

	transactionResponse, err := tc.TransactionService.CreateTransaction(c.Request().Context(), params.AccountID, params.Amount, model.TransactionTypeWithdraw)
	if err != nil {
		return c.JSON(500, model.DefaultError{Error: err.Error()})
	}

	return c.JSON(200, transactionResponse)
}

// @BasePath /
// Get Transaction GET
// @Summary API To get a transaction
// @Schemes
// @Description Api will return status 200 if the transaction is found, 404 if the transaction is not found and 500 if there is an internal server error
// @Tags Transaction
// @Accept json
// @Produce json
// @Success 200 {object} model.TransactionResponse
// @Failure 404 {object} model.DefaultError{error=string}
// @Failure 500 {object} model.DefaultError{error=string}
// @Param transaction_id path string true "Transaction ID"
// @Router /api/v1/transaction/{transaction_id} [get]
func (tc *TransactionController) GetTransaction(c echo.Context) error {
	transactionID := c.Param("transaction_id")
	transactionResponse, err := tc.TransactionService.GetTransaction(c.Request().Context(), transactionID)
	if err != nil {
		return c.JSON(500, model.DefaultError{Error: err.Error()})
	}

	if transactionResponse == nil {
		return c.JSON(404, model.DefaultError{Error: "transaction not found"})
	}

	return c.JSON(200, transactionResponse)
}

// @BasePath /
// Update Transaction PUT
// @Summary API To update a transaction
// @Schemes
// @Description Api will return status 200 if the transaction is updated, 400 if the request is invalid and 500 if there is an internal server error
// @Tags Transaction
// @Accept json
// @Produce json
// @Success 200 {object} model.DefaultResponse{data=string}
// @Failure 404 {object} model.DefaultError{error=string}
// @Failure 500 {object} model.DefaultError{error=string}
// @Param TransactionRequest body UpdateTransactionRequest true "Transaction Request"
// @Router /api/v1/transaction [put]
func (tc *TransactionController) UpdateTransaction(c echo.Context) error {
	params, err := tc.ValidateTransactionUpdateRequest(c)
	if err != nil {
		return c.JSON(400, model.DefaultError{Error: err.Error()})
	}

	err = tc.TransactionService.UpdateTransaction(c.Request().Context(), params.AccountID, params.TransactionID, params.Status)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return c.JSON(404, model.DefaultError{Error: "transaction not found"})
		}
		return c.JSON(500, model.DefaultError{Error: err.Error()})
	}

	return c.JSON(200, model.DefaultResponse{Data: "success"})
}

// ------------------Validation Methods------------------//
func (tc *TransactionController) ValidateTransactionRequest(c echo.Context) (*DepositRequest, error) {
	// use echo.Bind to bind the request body to the DepositRequest struct
	params := new(DepositRequest)
	if err := c.Bind(params); err != nil {
		return nil, fmt.Errorf("invalid request body: %v", err)
	}

	// validate the request body
	if params.AccountID == "" {
		return nil, fmt.Errorf("account_id is required")
	}

	if params.Amount.IsZero() {
		return nil, fmt.Errorf("amount is required")
	}

	if params.Amount.IsNegative() {
		return nil, fmt.Errorf("amount must be positive")
	}

	return params, nil
}

func (tc *TransactionController) ValidateTransactionUpdateRequest(c echo.Context) (*UpdateTransactionRequest, error) {
	// use echo.Bind to bind the request body to the UpdateTransactionRequest struct
	params := new(UpdateTransactionRequest)
	if err := c.Bind(params); err != nil {
		return nil, fmt.Errorf("invalid request body: %v", err)
	}

	// validate the request body
	if params.AccountID == "" {
		return nil, fmt.Errorf("account_id is required")
	}

	if params.TransactionID == "" {
		return nil, fmt.Errorf("transaction_id is required")
	}

	if params.Status == "" {
		return nil, fmt.Errorf("status is required")
	}

	if params.Status != model.TransactionStatusSuccess && params.Status != model.TransactionStatusFailed && params.Status != model.TransactionStatusPending {
		return nil, fmt.Errorf("invalid status value")
	}

	return params, nil
}
