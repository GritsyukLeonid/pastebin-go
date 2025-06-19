package grpcimpl

import (
	"context"
	"time"

	"github.com/GritsyukLeonid/pastebin-go/internal/model"
	"github.com/GritsyukLeonid/pastebin-go/internal/pb"
	"github.com/GritsyukLeonid/pastebin-go/internal/repository"
	"github.com/google/uuid"
)

type Server struct {
	pb.UnimplementedPasteServiceServer
	pb.UnimplementedUserServiceServer
	pb.UnimplementedStatsServiceServer
	pb.UnimplementedShortURLServiceServer
}

// --- User ---

func (s *Server) CreateUser(ctx context.Context, user *pb.User) (*pb.User, error) {
	u := &model.User{
		ID:       user.Id,
		Username: user.Username,
	}
	err := repository.AddUser(u)
	if err != nil {
		return nil, err
	}
	return &pb.User{
		Id:       u.ID,
		Username: u.Username,
	}, nil
}

func (s *Server) GetUser(ctx context.Context, req *pb.IDRequestInt) (*pb.User, error) {
	u, err := repository.GetUserByID(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.User{
		Id:       u.ID,
		Username: u.Username,
	}, nil
}

func (s *Server) ListUsers(req *pb.Empty, stream pb.UserService_ListUsersServer) error {
	users := repository.GetAllUsers()
	for _, u := range users {
		err := stream.Send(&pb.User{
			Id:       u.ID,
			Username: u.Username,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) UpdateUser(ctx context.Context, user *pb.User) (*pb.User, error) {
	u := &model.User{
		ID:       user.Id,
		Username: user.Username,
	}
	err := repository.UpdateUser(u.ID, u)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *pb.IDRequestInt) (*pb.Status, error) {
	err := repository.DeleteUser(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Status{Message: "User deleted"}, nil
}

// --- Paste ---

func (s *Server) CreatePaste(ctx context.Context, req *pb.Paste) (*pb.Paste, error) {
	paste := &model.Paste{
		ID:        uuid.New().String(),
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	err := repository.StoreObject(paste)
	if err != nil {
		return nil, err
	}

	return &pb.Paste{
		Id:        paste.ID,
		Content:   paste.Content,
		CreatedAt: paste.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *Server) GetPaste(ctx context.Context, req *pb.IDRequest) (*pb.Paste, error) {
	paste, err := repository.GetPasteByID(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Paste{
		Id:        paste.ID,
		Title:     "",
		Content:   paste.Content,
		CreatedAt: paste.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *Server) ListPastes(req *pb.Empty, stream pb.PasteService_ListPastesServer) error {
	pastes := repository.GetAllPastes()
	for _, p := range pastes {
		err := stream.Send(&pb.Paste{
			Id:        p.ID,
			Title:     "",
			Content:   p.Content,
			CreatedAt: p.CreatedAt.Format(time.RFC3339),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) UpdatePaste(ctx context.Context, req *pb.Paste) (*pb.Paste, error) {
	paste := &model.Paste{
		ID:      req.Id,
		Content: req.Content,
	}

	err := repository.UpdatePaste(paste.ID, paste)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (s *Server) DeletePaste(ctx context.Context, req *pb.IDRequest) (*pb.Status, error) {
	err := repository.DeletePaste(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Status{Message: "Paste deleted successfully"}, nil
}

// --- Stats ---

func (s *Server) CreateStats(ctx context.Context, req *pb.Stats) (*pb.Stats, error) {
	stats := &model.Stats{
		ID:    req.Id,
		Views: int(req.Views),
	}
	err := repository.StoreObject(stats)
	if err != nil {
		return nil, err
	}
	return &pb.Stats{
		Id:    stats.ID,
		Views: int64(stats.Views),
	}, nil
}

func (s *Server) GetStats(ctx context.Context, req *pb.IDRequest) (*pb.Stats, error) {
	stats, err := repository.GetStatsByID(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Stats{
		Id:    stats.ID,
		Views: int64(stats.Views),
	}, nil
}

func (s *Server) ListStats(req *pb.Empty, stream pb.StatsService_ListStatsServer) error {
	allStats := repository.GetAllStats()
	for _, st := range allStats {
		err := stream.Send(&pb.Stats{
			Id:    st.ID,
			Views: int64(st.Views),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) UpdateStats(ctx context.Context, req *pb.Stats) (*pb.Status, error) {
	stats := &model.Stats{
		ID:    req.Id,
		Views: int(req.Views),
	}
	err := repository.UpdateStats(stats.ID, stats)
	if err != nil {
		return nil, err
	}
	return &pb.Status{Message: "Update successful"}, nil
}

func (s *Server) DeleteStats(ctx context.Context, req *pb.IDRequest) (*pb.Status, error) {
	err := repository.DeleteStats(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Status{Message: "Delete successful"}, nil
}

// --- ShortURL ---

func (s *Server) CreateShortURL(ctx context.Context, req *pb.ShortURL) (*pb.ShortURL, error) {
	shorturl := &model.ShortURL{
		ID:       req.Id,
		Original: req.OriginalUrl,
	}
	err := repository.StoreObject(shorturl)
	if err != nil {
		return nil, err
	}
	return &pb.ShortURL{
		Id:          shorturl.ID,
		OriginalUrl: shorturl.Original,
	}, nil
}

func (s *Server) GetShortURL(ctx context.Context, req *pb.IDRequest) (*pb.ShortURL, error) {
	shorturl, err := repository.GetShortURLByID(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.ShortURL{
		Id:          shorturl.ID,
		OriginalUrl: shorturl.Original,
	}, nil
}

func (s *Server) ListShortURLs(req *pb.Empty, stream pb.ShortURLService_ListShortURLsServer) error {
	allShortURLs := repository.GetAllShortURLs()
	for _, su := range allShortURLs {
		err := stream.Send(&pb.ShortURL{
			Id:          su.ID,
			OriginalUrl: su.Original,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) UpdateShortURL(ctx context.Context, req *pb.ShortURL) (*pb.Status, error) {
	shorturl := &model.ShortURL{
		ID:       req.Id,
		Original: req.OriginalUrl,
	}
	err := repository.UpdateShortURL(shorturl.ID, shorturl)
	if err != nil {
		return nil, err
	}
	return &pb.Status{Message: "Update successful"}, nil
}

func (s *Server) DeleteShortURL(ctx context.Context, req *pb.IDRequest) (*pb.Status, error) {
	err := repository.DeleteShortURL(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Status{Message: "Delete successful"}, nil
}
