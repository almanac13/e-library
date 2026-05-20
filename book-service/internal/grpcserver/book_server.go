package grpcserver

import (
	"context"
	"log"

	"github.com/almanac13/e-library/book-service/gen/bookpb"
	"github.com/almanac13/e-library/book-service/internal/model"
	"github.com/almanac13/e-library/book-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BookGRPCServer struct {
	bookpb.UnimplementedBookServiceServer
	service *service.BookService
}

func NewBookGRPCServer(service *service.BookService) *BookGRPCServer {
	return &BookGRPCServer{
		service: service,
	}
}

func (s *BookGRPCServer) CreateBook(ctx context.Context, req *bookpb.CreateBookRequest) (*bookpb.BookResponse, error) {
	book, err := s.service.Create(model.CreateBookRequest{
		Title:     req.GetTitle(),
		Author:    req.GetAuthor(),
		Category:  req.GetCategory(),
		Available: req.GetAvailable(),
	})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &bookpb.BookResponse{Book: toProtoBook(book)}, nil
}

func (s *BookGRPCServer) GetBookByID(ctx context.Context, req *bookpb.GetBookByIDRequest) (*bookpb.BookResponse, error) {
	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "book id must be positive")
	}

	book, err := s.service.GetByID(int(req.GetId()))
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &bookpb.BookResponse{Book: toProtoBook(book)}, nil
}

func (s *BookGRPCServer) ListBooks(ctx context.Context, req *bookpb.ListBooksRequest) (*bookpb.ListBooksResponse, error) {
	books, err := s.service.GetAll()
	if err != nil {
		log.Println("failed to get books:", err)
		return nil, status.Error(codes.Internal, "failed to get books")
	}

	return toProtoBookList(books), nil
}

func (s *BookGRPCServer) UpdateBook(ctx context.Context, req *bookpb.UpdateBookRequest) (*bookpb.BookResponse, error) {
	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "book id must be positive")
	}

	book, err := s.service.Update(int(req.GetId()), model.UpdateBookRequest{
		Title:     req.GetTitle(),
		Author:    req.GetAuthor(),
		Category:  req.GetCategory(),
		Available: req.GetAvailable(),
	})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &bookpb.BookResponse{Book: toProtoBook(book)}, nil
}

func (s *BookGRPCServer) DeleteBook(ctx context.Context, req *bookpb.DeleteBookRequest) (*bookpb.DeleteBookResponse, error) {
	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "book id must be positive")
	}

	if err := s.service.Delete(int(req.GetId())); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &bookpb.DeleteBookResponse{
		Message: "book deleted successfully",
	}, nil
}

func (s *BookGRPCServer) SearchBooks(ctx context.Context, req *bookpb.SearchBooksRequest) (*bookpb.ListBooksResponse, error) {
	books, err := s.service.Search(req.GetQuery())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return toProtoBookList(books), nil
}

func (s *BookGRPCServer) ListBooksByAuthor(ctx context.Context, req *bookpb.ListBooksByAuthorRequest) (*bookpb.ListBooksResponse, error) {
	books, err := s.service.ListByAuthor(req.GetAuthor())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return toProtoBookList(books), nil
}

func (s *BookGRPCServer) ListBooksByCategory(ctx context.Context, req *bookpb.ListBooksByCategoryRequest) (*bookpb.ListBooksResponse, error) {
	books, err := s.service.ListByCategory(req.GetCategory())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return toProtoBookList(books), nil
}

func (s *BookGRPCServer) CheckBookAvailability(ctx context.Context, req *bookpb.CheckBookAvailabilityRequest) (*bookpb.AvailabilityResponse, error) {
	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "book id must be positive")
	}

	available, err := s.service.CheckAvailability(int(req.GetId()))
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &bookpb.AvailabilityResponse{
		BookId:    req.GetId(),
		Available: available,
	}, nil
}

func (s *BookGRPCServer) MarkBookAvailable(ctx context.Context, req *bookpb.MarkBookAvailableRequest) (*bookpb.BookResponse, error) {
	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "book id must be positive")
	}

	book, err := s.service.MarkAvailable(int(req.GetId()))
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &bookpb.BookResponse{Book: toProtoBook(book)}, nil
}

func (s *BookGRPCServer) MarkBookUnavailable(ctx context.Context, req *bookpb.MarkBookUnavailableRequest) (*bookpb.BookResponse, error) {
	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "book id must be positive")
	}

	book, err := s.service.MarkUnavailable(int(req.GetId()))
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &bookpb.BookResponse{Book: toProtoBook(book)}, nil
}

func (s *BookGRPCServer) GetBookStats(ctx context.Context, req *bookpb.GetBookStatsRequest) (*bookpb.BookStatsResponse, error) {
	stats, err := s.service.GetStats()
	if err != nil {
		log.Println("failed to get book stats:", err)
		return nil, status.Error(codes.Internal, "failed to get book stats")
	}

	return &bookpb.BookStatsResponse{
		TotalBooks:       int32(stats.TotalBooks),
		AvailableBooks:   int32(stats.AvailableBooks),
		UnavailableBooks: int32(stats.UnavailableBooks),
	}, nil
}

func toProtoBookList(books []model.Book) *bookpb.ListBooksResponse {
	protoBooks := make([]*bookpb.Book, 0, len(books))

	for _, book := range books {
		bookCopy := book
		protoBooks = append(protoBooks, toProtoBook(&bookCopy))
	}

	return &bookpb.ListBooksResponse{
		Books: protoBooks,
	}
}

func toProtoBook(book *model.Book) *bookpb.Book {
	return &bookpb.Book{
		Id:        int32(book.ID),
		Title:     book.Title,
		Author:    book.Author,
		Category:  book.Category,
		Available: book.Available,
		CreatedAt: book.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
