package gapi

import (
	"context"
	"errors"
	"time"
	"tunelport-svc-totems/internal/domain/persistence"
	"tunelport-svc-totems/internal/domain/services"
	"tunelport-svc-totems/internal/domain/types"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type Server struct {
	UnimplementedTotemServiceServer
	repository persistence.TotemRepository
}

func mapDomainError(err error, internalMsg string) error {
	switch {
	case errors.Is(err, types.ErrDeviceAlreadyExists):
		return status.Errorf(codes.AlreadyExists, "totem already exists: %v", err)
	case errors.Is(err, types.ErrDeviceNotFound):
		return status.Errorf(codes.NotFound, "totem not found: %v", err)
	case errors.Is(err, types.ErrInvalidStatus):
		return status.Errorf(codes.InvalidArgument, "invalid status: %v", err)
	default:
		return status.Errorf(codes.Internal, "%s: %v", internalMsg, err)
	}
}
func NewServer(repository persistence.TotemRepository) *Server {
	return &Server{
		repository: repository,
	}
}

func (s *Server) CreateTotem(ctx context.Context, req *CreateTotemRequest) (*CreateTotemResponse, error) {
	if req == nil || req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	service := services.NewSaveTotemUseCase(s.repository)
	id, err := service.Do(req.Name, time.Now(), time.Now())
	if err != nil {
		return nil, mapDomainError(err, "failed to create totem")
	}

	return &CreateTotemResponse{
		Id:   id,
		Name: req.Name,
	}, nil
}

func (s *Server) GetTotem(ctx context.Context, req *GetTotemRequest) (*GetTotemResponse, error) {
	if req == nil || req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	service := services.NewGetTotemUseCase(s.repository)
	totem, err := service.Do(req.Id)
	if err != nil {
		return nil, mapDomainError(err, "failed to get totem")
	}

	return &GetTotemResponse{
		Id:   totem.ID.Value(),
		Name: totem.Name.Value(),
	}, nil
}

func (s *Server) UpdateTotem(ctx context.Context, req *UpdateTotemRequest) (*UpdateTotemResponse, error) {
	if req == nil || req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	service := services.NewUpdateTotemUseCase(s.repository)
	err := service.Do(req.Id, req.Name, time.Now(), time.Now())
	if err != nil {
		return nil, mapDomainError(err, "failed to update totem")
	}

	return &UpdateTotemResponse{
		Updated: time.Now().Format(time.RFC3339),
	}, nil
}

func (s *Server) DeleteTotem(ctx context.Context, req *DeleteTotemRequest) (*DeleteTotemResponse, error) {
	if req == nil || req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	service := services.NewDeleteTotemUseCase(s.repository)
	err := service.Do(req.Id)
	if err != nil {
		return nil, mapDomainError(err, "failed to delete totem")
	}

	return &DeleteTotemResponse{
		Deleted: time.Now().Format(time.RFC3339),
	}, nil
}

func (s *Server) ListTotems(ctx context.Context, req *ListTotemsRequest) (*ListTotemsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	service := services.NewListTotemUseCase(s.repository)
	totems, err := service.Do()
	if err != nil {
		return nil, mapDomainError(err, "failed to list totems")
	}

	items := make([]*Totem, 0, len(totems))
	for _, t := range totems {
		items = append(items, &Totem{
			Id:   t.ID.Value(),
			Name: t.Name.Value(),
		})
	}

	return &ListTotemsResponse{
		Totems: items,
	}, nil
}
