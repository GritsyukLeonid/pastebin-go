package main

import (
	"context"
	"log"
	"time"

	pb "github.com/GritsyukLeonid/pastebin-go/internal/pb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("не удалось подключиться: %v", err)
	}
	defer conn.Close()

	userClient := pb.NewUserServiceClient(conn)
	pasteClient := pb.NewPasteServiceClient(conn)
	statsClient := pb.NewStatsServiceClient(conn)
	shortURLClient := pb.NewShortURLServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ==== USER ====
	userResp, err := userClient.CreateUser(ctx, &pb.User{
		Username: "test_user",
		Email:    "test@example.com",
		Password: "123456",
	})
	if err != nil {
		log.Fatalf("CreateUser error: %v", err)
	}
	log.Printf("User created: %v", userResp)

	userID := userResp.Id

	getUser, err := userClient.GetUser(ctx, &pb.IDRequestInt{Id: userID})
	if err != nil {
		log.Fatalf("GetUser error: %v", err)
	}
	log.Printf("Fetched user by ID: %v", getUser)

	allUsersStream, err := userClient.ListUsers(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("ListUsers error: %v", err)
	}
	log.Println("All users:")
	for {
		user, err := allUsersStream.Recv()
		if err != nil {
			break
		}
		log.Printf("- %v", user)
	}

	_, err = userClient.UpdateUser(ctx, &pb.User{
		Id:       userID,
		Username: "updated_user",
		Email:    "updated@example.com",
		Password: "654321",
	})
	if err != nil {
		log.Fatalf("UpdateUser error: %v", err)
	}

	_, err = userClient.DeleteUser(ctx, &pb.IDRequestInt{Id: userID})
	if err != nil {
		log.Fatalf("DeleteUser error: %v", err)
	}

	// ==== PASTE ====
	pasteResp, err := pasteClient.CreatePaste(ctx, &pb.Paste{
		Title:   "Test title",
		Content: "Hello world!",
	})
	if err != nil {
		log.Fatalf("CreatePaste error: %v", err)
	}
	log.Printf("Paste created: %v", pasteResp)

	pasteID := pasteResp.Id

	pasteByID, err := pasteClient.GetPaste(ctx, &pb.IDRequest{Id: pasteID})
	if err != nil {
		log.Fatalf("GetPaste error: %v", err)
	}
	log.Printf("Fetched paste: %v", pasteByID)

	_, err = pasteClient.UpdatePaste(ctx, &pb.Paste{
		Id:      pasteID,
		Content: "Updated content",
		UserId:  userID,
	})
	if err != nil {
		log.Fatalf("UpdatePaste error: %v", err)
	}

	_, err = pasteClient.DeletePaste(ctx, &pb.IDRequest{Id: pasteID})
	if err != nil {
		log.Fatalf("DeletePaste error: %v", err)
	}

	// ==== SHORT URL ====
	shortResp, err := shortURLClient.CreateShortURL(ctx, &pb.ShortURL{
		OriginalUrl: "https://example.com",
		ShortCode:   "exmpl",
	})
	if err != nil {
		log.Fatalf("CreateShortURL error: %v", err)
	}
	shortID := shortResp.Id
	log.Printf("ShortURL created: %v", shortResp)

	_, err = shortURLClient.UpdateShortURL(ctx, &pb.ShortURL{
		Id:          shortID,
		OriginalUrl: "https://updated.com",
		ShortCode:   "updcd",
	})
	if err != nil {
		log.Fatalf("UpdateShortURL error: %v", err)
	}

	_, err = shortURLClient.DeleteShortURL(ctx, &pb.IDRequest{Id: shortID})
	if err != nil {
		log.Fatalf("DeleteShortURL error: %v", err)
	}

	// ==== STATS ====

	statsResp, err := statsClient.CreateStats(ctx, &pb.Stats{
		Id:      "stat1",
		PasteId: pasteID,
		Views:   5,
	})
	if err != nil {
		log.Fatalf("CreateStats error: %v", err)
	}
	log.Printf("Stats created: %v", statsResp)

	statsID := statsResp.Id

	_, err = statsClient.UpdateStats(ctx, &pb.Stats{
		Id:      statsID,
		PasteId: pasteID,
		Views:   10,
	})
	if err != nil {
		log.Fatalf("UpdateStats error: %v", err)
	}

	_, err = statsClient.DeleteStats(ctx, &pb.IDRequest{Id: statsID})
	if err != nil {
		log.Fatalf("DeleteStats error: %v", err)
	}
}
