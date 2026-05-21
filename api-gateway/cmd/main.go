package main

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"api-gateway/gen/bookpb"
	"api-gateway/gen/borrowpb"
	"api-gateway/gen/userpb"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	userClient   userpb.UserServiceClient
	bookClient   bookpb.BookServiceClient
	borrowClient borrowpb.BorrowServiceClient
}

func main() {
	userConn := mustConnect(getEnv("USER_SERVICE_GRPC_ADDR", "user-service:50051"))
	bookConn := mustConnect(getEnv("BOOK_SERVICE_GRPC_ADDR", "book-service:50052"))
	borrowConn := mustConnect(getEnv("BORROW_SERVICE_GRPC_ADDR", "borrow-service:50053"))

	defer userConn.Close()
	defer bookConn.Close()
	defer borrowConn.Close()

	server := &Server{
		userClient:   userpb.NewUserServiceClient(userConn),
		bookClient:   bookpb.NewBookServiceClient(bookConn),
		borrowClient: borrowpb.NewBorrowServiceClient(borrowConn),
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "api-gateway is running"})
	})

	api := router.Group("/api")
	{
		api.POST("/users/register", server.RegisterUser)
		api.POST("/users/login", server.LoginUser)
		api.GET("/users", server.GetAllUsers)
		api.GET("/users/count", server.CountUsers)
		api.GET("/users/email/:email", server.GetUserByEmail)
		api.GET("/users/role/:role", server.GetUsersByRole)
		api.GET("/users/:id", server.GetUserByID)
		api.GET("/users/:id/exists", server.CheckUserExists)
		api.PUT("/users/:id/name", server.UpdateUserName)
		api.PUT("/users/:id/role", server.UpdateUserRole)
		api.PUT("/users/:id/password", server.ChangePassword)
		api.DELETE("/users/:id", server.DeleteUser)

		api.POST("/books", server.CreateBook)
		api.GET("/books", server.ListBooks)
		api.GET("/books/search", server.SearchBooks)
		api.GET("/books/stats", server.GetBookStats)
		api.GET("/books/author/:author", server.ListBooksByAuthor)
		api.GET("/books/category/:category", server.ListBooksByCategory)
		api.GET("/books/:id", server.GetBookByID)
		api.GET("/books/:id/availability", server.CheckBookAvailability)
		api.PUT("/books/:id", server.UpdateBook)
		api.PUT("/books/:id/available", server.MarkBookAvailable)
		api.PUT("/books/:id/unavailable", server.MarkBookUnavailable)
		api.DELETE("/books/:id", server.DeleteBook)

		api.POST("/borrows", server.CreateBorrow)
		api.GET("/borrows", server.GetAllBorrows)
		api.GET("/borrows/overdue", server.GetOverdueBorrows)
		api.GET("/borrows/active", server.GetActiveBorrows)
		api.GET("/borrows/count", server.CountBorrows)
		api.GET("/borrows/user/:userId", server.GetBorrowsByUserID)
		api.GET("/borrows/book/:bookId", server.GetBorrowsByBookID)
		api.GET("/borrows/:id", server.GetBorrow)
		api.GET("/borrows/:id/exists", server.CheckBorrowExists)
		api.PUT("/borrows/:id/return", server.ReturnBorrow)
		api.PUT("/borrows/:id/extend", server.ExtendBorrowPeriod)
		api.PUT("/borrows/:id/cancel", server.CancelBorrow)
	}

	port := getEnv("API_GATEWAY_PORT", "8080")
	if err := router.Run(":" + port); err != nil {
		panic(err)
	}
}

func mustConnect(address string) *grpc.ClientConn {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	return conn
}

func timeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func parseInt32(value string) (int32, error) {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return int32(parsed), nil
}

func handleError(c *gin.Context, err error) {
	c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
}

func (s *Server) RegisterUser(c *gin.Context) {
	var req userpb.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.userClient.RegisterUser(ctx, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (s *Server) LoginUser(c *gin.Context) {
	var req userpb.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.userClient.LoginUser(ctx, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) GetAllUsers(c *gin.Context) {
	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.userClient.GetAllUsers(ctx, &userpb.GetAllUsersRequest{})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) GetUserByID(c *gin.Context) {
	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.userClient.GetUserByID(ctx, &userpb.GetUserByIDRequest{
		Id: c.Param("id"),
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) GetUserByEmail(c *gin.Context) {
	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.userClient.GetUserByEmail(ctx, &userpb.GetUserByEmailRequest{
		Email: c.Param("email"),
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) UpdateUserName(c *gin.Context) {
	var req userpb.UpdateUserNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Id = c.Param("id")

	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.userClient.UpdateUserName(ctx, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) UpdateUserRole(c *gin.Context) {
	var req userpb.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Id = c.Param("id")

	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.userClient.UpdateUserRole(ctx, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) ChangePassword(c *gin.Context) {
	var req userpb.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Id = c.Param("id")

	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.userClient.ChangePassword(ctx, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) DeleteUser(c *gin.Context) {
	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.userClient.DeleteUser(ctx, &userpb.DeleteUserRequest{
		Id: c.Param("id"),
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) CheckUserExists(c *gin.Context) {
	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.userClient.CheckUserExists(ctx, &userpb.CheckUserExistsRequest{
		Id: c.Param("id"),
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) CountUsers(c *gin.Context) {
	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.userClient.CountUsers(ctx, &userpb.CountUsersRequest{})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) GetUsersByRole(c *gin.Context) {
	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.userClient.GetUsersByRole(ctx, &userpb.GetUsersByRoleRequest{
		Role: c.Param("role"),
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) CreateBook(c *gin.Context) {
	var req bookpb.CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.bookClient.CreateBook(ctx, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (s *Server) ListBooks(c *gin.Context) {
	page, _ := parseInt32(c.DefaultQuery("page", "1"))
	limit, _ := parseInt32(c.DefaultQuery("limit", "10"))

	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.bookClient.ListBooks(ctx, &bookpb.ListBooksRequest{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) GetBookByID(c *gin.Context) {
	id, err := parseInt32(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}

	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.bookClient.GetBookByID(ctx, &bookpb.GetBookByIDRequest{
		Id: id,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) UpdateBook(c *gin.Context) {
	id, err := parseInt32(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}

	var req bookpb.UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Id = id

	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.bookClient.UpdateBook(ctx, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) DeleteBook(c *gin.Context) {
	id, err := parseInt32(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}

	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.bookClient.DeleteBook(ctx, &bookpb.DeleteBookRequest{
		Id: id,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) SearchBooks(c *gin.Context) {
	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.bookClient.SearchBooks(ctx, &bookpb.SearchBooksRequest{
		Query: c.Query("q"),
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) ListBooksByAuthor(c *gin.Context) {
	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.bookClient.ListBooksByAuthor(ctx, &bookpb.ListBooksByAuthorRequest{
		Author: c.Param("author"),
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) ListBooksByCategory(c *gin.Context) {
	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.bookClient.ListBooksByCategory(ctx, &bookpb.ListBooksByCategoryRequest{
		Category: c.Param("category"),
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) CheckBookAvailability(c *gin.Context) {
	id, err := parseInt32(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}

	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.bookClient.CheckBookAvailability(ctx, &bookpb.CheckBookAvailabilityRequest{
		Id: id,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) MarkBookAvailable(c *gin.Context) {
	id, err := parseInt32(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}

	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.bookClient.MarkBookAvailable(ctx, &bookpb.MarkBookAvailableRequest{
		Id: id,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) MarkBookUnavailable(c *gin.Context) {
	id, err := parseInt32(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid book id"})
		return
	}

	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.bookClient.MarkBookUnavailable(ctx, &bookpb.MarkBookUnavailableRequest{
		Id: id,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) GetBookStats(c *gin.Context) {
	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.bookClient.GetBookStats(ctx, &bookpb.GetBookStatsRequest{})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) CreateBorrow(c *gin.Context) {
	var req borrowpb.CreateBorrowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.borrowClient.CreateBorrow(ctx, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (s *Server) GetAllBorrows(c *gin.Context) {
	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.borrowClient.GetAllBorrows(ctx, &borrowpb.GetAllBorrowsRequest{})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) GetBorrow(c *gin.Context) {
	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.borrowClient.GetBorrow(ctx, &borrowpb.GetBorrowRequest{
		Id: c.Param("id"),
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) ReturnBorrow(c *gin.Context) {
	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.borrowClient.ReturnBorrow(ctx, &borrowpb.ReturnBorrowRequest{
		Id: c.Param("id"),
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) ExtendBorrowPeriod(c *gin.Context) {
	var req borrowpb.ExtendBorrowPeriodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Id = c.Param("id")

	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.borrowClient.ExtendBorrowPeriod(ctx, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) CancelBorrow(c *gin.Context) {
	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.borrowClient.CancelBorrow(ctx, &borrowpb.CancelBorrowRequest{
		Id: c.Param("id"),
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) GetBorrowsByUserID(c *gin.Context) {
	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.borrowClient.GetBorrowsByUserID(ctx, &borrowpb.GetBorrowsByUserIDRequest{
		UserId: c.Param("userId"),
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) GetBorrowsByBookID(c *gin.Context) {
	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.borrowClient.GetBorrowsByBookID(ctx, &borrowpb.GetBorrowsByBookIDRequest{
		BookId: c.Param("bookId"),
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) GetOverdueBorrows(c *gin.Context) {
	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.borrowClient.GetOverdueBorrows(ctx, &borrowpb.GetOverdueBorrowsRequest{})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) GetActiveBorrows(c *gin.Context) {
	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.borrowClient.GetActiveBorrows(ctx, &borrowpb.GetActiveBorrowsRequest{})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) CountBorrows(c *gin.Context) {
	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.borrowClient.CountBorrows(ctx, &borrowpb.CountBorrowsRequest{})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Server) CheckBorrowExists(c *gin.Context) {
	ctx, cancel := timeoutContext()
	defer cancel()

	res, err := s.borrowClient.CheckBorrowExists(ctx, &borrowpb.CheckBorrowExistsRequest{
		Id: c.Param("id"),
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}
