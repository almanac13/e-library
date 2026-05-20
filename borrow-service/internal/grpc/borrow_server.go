package grpc

import (
	"context"

	"borrow-service/gen/borrowpb"
	"borrow-service/internal/service"
)

type BorrowGRPCServer struct {
	borrowpb.UnimplementedBorrowServiceServer
	service *service.BorrowService
}

func NewBorrowGRPCServer(service *service.BorrowService) *BorrowGRPCServer {
	return &BorrowGRPCServer{
		service: service,
	}
}

func (s *BorrowGRPCServer) GetBorrow(
	ctx context.Context,
	req *borrowpb.GetBorrowRequest,
) (*borrowpb.BorrowResponse, error) {

	borrow, err := s.service.GetBorrowByID(req.Id)
	if err != nil {
		return nil, err
	}

	return &borrowpb.BorrowResponse{
		Id:         borrow.ID,
		UserId:     borrow.UserID,
		BookId:     borrow.BookID,
		BorrowDate: borrow.BorrowDate.String(),
		DueDate:    borrow.DueDate.String(),
		Status:     borrow.Status,
		CreatedAt:  borrow.CreatedAt.String(),
		UpdatedAt:  borrow.UpdatedAt.String(),
	}, nil
}
func (s *BorrowGRPCServer) GetAllBorrows(
	ctx context.Context,
	req *borrowpb.GetAllBorrowsRequest,
) (*borrowpb.GetAllBorrowsResponse, error) {

	borrows, err := s.service.GetAllBorrows()
	if err != nil {
		return nil, err
	}

	response := &borrowpb.GetAllBorrowsResponse{
		Borrows: make([]*borrowpb.BorrowResponse, 0, len(borrows)),
	}

	for _, borrow := range borrows {
		item := &borrowpb.BorrowResponse{
			Id:         borrow.ID,
			UserId:     borrow.UserID,
			BookId:     borrow.BookID,
			BorrowDate: borrow.BorrowDate.String(),
			DueDate:    borrow.DueDate.String(),
			Status:     borrow.Status,
			CreatedAt:  borrow.CreatedAt.String(),
			UpdatedAt:  borrow.UpdatedAt.String(),
		}

		if borrow.ReturnDate != nil {
			item.ReturnDate = borrow.ReturnDate.String()
		}

		response.Borrows = append(response.Borrows, item)
	}

	return response, nil
}

func (s *BorrowGRPCServer) ReturnBorrow(
	ctx context.Context,
	req *borrowpb.ReturnBorrowRequest,
) (*borrowpb.BorrowResponse, error) {

	borrow, err := s.service.ReturnBorrow(req.Id)
	if err != nil {
		return nil, err
	}

	response := &borrowpb.BorrowResponse{
		Id:         borrow.ID,
		UserId:     borrow.UserID,
		BookId:     borrow.BookID,
		BorrowDate: borrow.BorrowDate.String(),
		DueDate:    borrow.DueDate.String(),
		Status:     borrow.Status,
		CreatedAt:  borrow.CreatedAt.String(),
		UpdatedAt:  borrow.UpdatedAt.String(),
	}

	if borrow.ReturnDate != nil {
		response.ReturnDate = borrow.ReturnDate.String()
	}

	return response, nil
}
func (s *BorrowGRPCServer) ExtendBorrowPeriod(
	ctx context.Context,
	req *borrowpb.ExtendBorrowPeriodRequest,
) (*borrowpb.BorrowResponse, error) {

	return &borrowpb.BorrowResponse{
		Id:     req.Id,
		Status: "extended",
	}, nil
}

func (s *BorrowGRPCServer) CancelBorrow(
	ctx context.Context,
	req *borrowpb.CancelBorrowRequest,
) (*borrowpb.BorrowResponse, error) {

	return &borrowpb.BorrowResponse{
		Id:     req.Id,
		Status: "cancelled",
	}, nil
}

func (s *BorrowGRPCServer) GetBorrowsByUserID(
	ctx context.Context,
	req *borrowpb.GetBorrowsByUserIDRequest,
) (*borrowpb.GetAllBorrowsResponse, error) {

	return &borrowpb.GetAllBorrowsResponse{}, nil
}

func (s *BorrowGRPCServer) GetBorrowsByBookID(
	ctx context.Context,
	req *borrowpb.GetBorrowsByBookIDRequest,
) (*borrowpb.GetAllBorrowsResponse, error) {

	return &borrowpb.GetAllBorrowsResponse{}, nil
}

func (s *BorrowGRPCServer) GetOverdueBorrows(
	ctx context.Context,
	req *borrowpb.GetOverdueBorrowsRequest,
) (*borrowpb.GetAllBorrowsResponse, error) {

	return &borrowpb.GetAllBorrowsResponse{}, nil
}

func (s *BorrowGRPCServer) GetActiveBorrows(
	ctx context.Context,
	req *borrowpb.GetActiveBorrowsRequest,
) (*borrowpb.GetAllBorrowsResponse, error) {

	return &borrowpb.GetAllBorrowsResponse{}, nil
}

func (s *BorrowGRPCServer) CountBorrows(
	ctx context.Context,
	req *borrowpb.CountBorrowsRequest,
) (*borrowpb.CountBorrowsResponse, error) {

	return &borrowpb.CountBorrowsResponse{
		Count: 0,
	}, nil
}

func (s *BorrowGRPCServer) CheckBorrowExists(
	ctx context.Context,
	req *borrowpb.CheckBorrowExistsRequest,
) (*borrowpb.CheckBorrowExistsResponse, error) {

	return &borrowpb.CheckBorrowExistsResponse{
		Exists: true,
	}, nil
}
