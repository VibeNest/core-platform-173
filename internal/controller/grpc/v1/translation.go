package v1

import (
	"context"

	"gitverse.ru/apavlov-systems/core-platform/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// translationRoutes связывает сгенерированный gRPC сервер и твой UseCase
type translationRoutes struct {
	UnimplementedTranslationServer // Обязательно для совместимости
	t                              usecase.Translation
}

func NewTranslationRoutes(t usecase.Translation) *translationRoutes {
	return &translationRoutes{t: t}
}

// Translate реализует метод из .proto файла
func (s *translationRoutes) Translate(ctx context.Context, req *TranslateRequest) (*TranslateResponse, error) {
	res, err := s.t.Translate(ctx, req.Source, req.Destination, req.Original)
	if err != nil {
		return nil, status.Error(codes.Internal, "translate error")
	}

	return &TranslateResponse{Translation: res.Translation}, nil
}

// GetHistory реализует метод получения истории
func (s *translationRoutes) GetHistory(ctx context.Context, req *GetHistoryRequest) (*GetHistoryResponse, error) {
	history, err := s.t.History(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "history error")
	}

	var items []*HistoryItem
	for _, h := range history {
		items = append(items, &HistoryItem{
			Source:      h.Source,
			Destination: h.Destination,
			Original:    h.Original,
			Translation: h.Translation,
		})
	}

	return &GetHistoryResponse{History: items}, nil
}
