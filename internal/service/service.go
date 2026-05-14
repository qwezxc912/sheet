package service

import (
	"context"
	"fmt"

	models "github.com/qweq1232/dnd_form/internal/domane/models/char"
)

const (
	emptyValue = 0
)

type Service struct {
	cp CharProvider
	cs CharSaver
	cu CharUpdater
	cd CharDeleter
}

type CharProvider interface {
	AllChar(ctx context.Context, userID int32) ([]models.Char, error)
	Char(ctx context.Context, id, userID int32) (*models.Char, error)
}

type CharSaver interface {
	SaveChar(ctx context.Context, stats []byte, userID int32) (int32, error)
}

type CharUpdater interface {
	UpdateChar(ctx context.Context, stats []byte, id int32) error
}

type CharDeleter interface {
	DeleteChar(ctx context.Context, id int32) error
}

func New(cp CharProvider,
	cs CharSaver,
	cu CharUpdater,
	cd CharDeleter,
) *Service {
	return &Service{
		cp: cp,
		cs: cs,
		cu: cu,
		cd: cd,
	}
}

func (s *Service) CreateChar(
	ctx context.Context,
	char models.Char,
) (int32, error) {
	const op = "service.char.service.CreateChar"

	id, err := s.cs.SaveChar(ctx, char.Stats, char.UserID)
	if err != nil {
		return emptyValue, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Service) Char(
	ctx context.Context,
	id int32,
	userID int32,
) (*models.Char, error) {
	const op = "service.char.service.Char"

	char, err := s.cp.Char(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return char, nil
}

func (s *Service) AllChar(
	ctx context.Context,
	userId int32,
) ([]models.Char, error) {
	const op = "service.char.service.AllChar"

	chars, err := s.cp.AllChar(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return chars, nil
}

func (s *Service) UpdateChar(
	ctx context.Context,
	char models.Char,
) error {
	const op = "service.char.service.UpdateChar"

	if err := s.cu.UpdateChar(ctx, char.Stats, char.ID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) DeleteChar(
	ctx context.Context,
	id int32,
) error {
	const op = "service.char.service.UpdateChar"

	if err := s.cd.DeleteChar(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
