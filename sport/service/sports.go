package service

import (
	"github.com/AndrewTheJavaGuy/entain/sport/db"
	"github.com/AndrewTheJavaGuy/entain/sport/proto/sport"
	"golang.org/x/net/context"
)

type Sports interface {
	// ListSports will return a collection of sport.
	ListSports(ctx context.Context, in *sport.ListSportsRequest) (*sport.ListSportsResponse, error)
}


// sportService implements the Racing interface.
type sportsService struct {
	sportsRepo db.SportsRepo
}

// NewRacingService instantiates and returns a new sportService.
func NewSportsService(sportsRepo db.SportsRepo) Sports {
	return &sportsService{sportsRepo}
}

func (s *sportsService) ListSports(ctx context.Context, in *sport.ListSportsRequest) (*sport.ListSportsResponse, error) {
	sportsResp, err := s.sportsRepo.List(in.Filter)
	if err != nil {
		return nil, err
	}

	return &sport.ListSportsResponse{Sport: sportsResp}, nil
}
