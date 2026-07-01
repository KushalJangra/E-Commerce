package internal

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kushaljangra/e-commerce/order/models"
	paymentpb "github.com/kushaljangra/e-commerce/payment/proto/pb"
	"github.com/kushaljangra/e-commerce/pkg/auth"
)

var errUnauthorized = errors.New("unauthorized")

type paginationQuery struct {
	Skip uint64 `form:"skip"`
	Take uint64 `form:"take"`
}

func bounds(p paginationQuery) (uint64, uint64) {
	if p.Take == 0 {
		return p.Skip, 100
	}
	return p.Skip, p.Take
}

func requireUser(c *gin.Context) (int, error) {
	accountID, err := auth.GetUserIdInt(c.Request.Context(), false)
	if err != nil || accountID == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return 0, errUnauthorized
	}
	return accountID, nil
}

func (s *Server) Register(c *gin.Context) {
	var body struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	token, err := s.Account.Register(ctx, body.Name, body.Email, body.Password)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", token, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusCreated, gin.H{"token": token})
}

func (s *Server) Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	token, err := s.Account.Login(ctx, body.Email, body.Password)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", token, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (s *Server) ListAccounts(c *gin.Context) {
	var query paginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	skip, take := bounds(query)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	accounts, err := s.Account.GetAccounts(ctx, skip, take)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accounts": accounts})
}

func (s *Server) GetAccount(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	account, err := s.Account.GetAccount(ctx, id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (s *Server) GetAccountOrders(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	orders, err := s.Order.GetOrdersForAccount(ctx, id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

func (s *Server) ListProducts(c *gin.Context) {
	var query struct {
		paginationQuery
		Query string `form:"q"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	skip, take := bounds(query.paginationQuery)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	products, err := s.Product.GetProducts(ctx, skip, take, nil, query.Query)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

func (s *Server) GetProduct(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	product, err := s.Product.GetProduct(ctx, c.Param("id"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (s *Server) CreateProduct(c *gin.Context) {
	accountID, err := requireUser(c)
	if err != nil {
		return
	}

	var body struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description" binding:"required"`
		Price       float64 `json:"price" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	product, err := s.Product.PostProduct(ctx, body.Name, body.Description, body.Price, int64(accountID))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (s *Server) UpdateProduct(c *gin.Context) {
	accountID, err := requireUser(c)
	if err != nil {
		return
	}

	var body struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description" binding:"required"`
		Price       float64 `json:"price" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	product, err := s.Product.UpdateProduct(ctx, c.Param("id"), body.Name, body.Description, body.Price, int64(accountID))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (s *Server) DeleteProduct(c *gin.Context) {
	accountID, err := requireUser(c)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	if err := s.Product.DeleteProduct(ctx, c.Param("id"), int64(accountID)); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (s *Server) CreateOrder(c *gin.Context) {
	accountID, err := requireUser(c)
	if err != nil {
		return
	}

	var body struct {
		Products []struct {
			ID       string `json:"id" binding:"required"`
			Quantity int    `json:"quantity" binding:"required,gt=0"`
		} `json:"products" binding:"required,min=1,dive"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var products []*models.OrderedProduct
	for _, product := range body.Products {
		products = append(products, &models.OrderedProduct{
			ID:       product.ID,
			Quantity: uint32(product.Quantity),
		})
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	order, err := s.Order.PostOrder(ctx, uint64(accountID), products)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (s *Server) CreateCustomerPortalSession(c *gin.Context) {
	var body struct {
		AccountID int    `json:"accountId" binding:"required"`
		Email     string `json:"email" binding:"required,email"`
		Name      string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	url, err := s.Payment.CreateCustomerPortalSession(ctx, uint64(body.AccountID), body.Email, body.Name)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

func (s *Server) CreateCheckoutSession(c *gin.Context) {
	var body struct {
		AccountID   int    `json:"accountId" binding:"required"`
		Email       string `json:"email" binding:"required,email"`
		Name        string `json:"name" binding:"required"`
		RedirectURL string `json:"redirectUrl" binding:"required"`
		OrderID     int    `json:"orderId" binding:"required"`
		Products    []struct {
			ID       string `json:"id" binding:"required"`
			Quantity int    `json:"quantity" binding:"required,gt=0"`
		} `json:"products" binding:"required,min=1,dive"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var products []*paymentpb.CartItem
	for _, product := range body.Products {
		products = append(products, &paymentpb.CartItem{
			ProductId: product.ID,
			Quantity:  uint64(product.Quantity),
		})
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	url, err := s.Payment.CreateCheckoutSession(
		ctx,
		body.OrderID,
		body.AccountID,
		body.Email,
		body.Name,
		body.RedirectURL,
		products,
	)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}
