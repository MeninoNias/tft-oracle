package player

import (
	"context"

	"connectrpc.com/connect"

	tftv1 "github.com/MeninoNias/tft-oracle/backend/gen/tft/v1"
	"github.com/MeninoNias/tft-oracle/backend/gen/tft/v1/tftv1connect"
)

var _ tftv1connect.PlayerServiceHandler = (*Service)(nil)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetPlayerProfile(
	ctx context.Context,
	req *connect.Request[tftv1.GetPlayerProfileRequest],
) (*connect.Response[tftv1.GetPlayerProfileResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, nil)
}

func (s *Service) GetMatchHistory(
	ctx context.Context,
	req *connect.Request[tftv1.GetMatchHistoryRequest],
) (*connect.Response[tftv1.GetMatchHistoryResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, nil)
}
