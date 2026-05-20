package main

import (
	"context"
	"log"
	"time"

	"borrow-service/gen/borrowpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:50053",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := borrowpb.NewBorrowServiceClient(conn)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	response, err := client.GetAllBorrows(
		ctx,
		&borrowpb.GetAllBorrowsRequest{},
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Borrows count:", len(response.Borrows))

	for _, b := range response.Borrows {
		log.Printf(
			"ID=%s User=%s Book=%s Status=%s",
			b.Id,
			b.UserId,
			b.BookId,
			b.Status,
		)
	}
}
