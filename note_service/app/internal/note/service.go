package note

import (
	"context"
	"errors"
	"fmt"
	"github.com/theartofdevel/notes_system/note_service/internal/apperror"
	"github.com/theartofdevel/notes_system/note_service/pkg/logging"
)

var _ Service = &service{}

type service struct {
	storage Storage
	logger  logging.Logger
}

func NewService(noteStorage Storage, logger logging.Logger) (Service, error) {
	return &service{
		storage: noteStorage,
		logger:  logger,
	}, nil
}

type Service interface {
	Create(ctx context.Context, dto CreateNoteDTO) (string, error)
	GetOne(ctx context.Context, uuid string) (Note, error)
	GetByCategoryUUID(ctx context.Context, uuid string) ([]Note, error)
	Update(ctx context.Context, uuid string, dto UpdateNoteDTO, tagsUpdate bool) error
	Delete(ctx context.Context, uuid string) error
}

func (s service) Create(ctx context.Context, dto CreateNoteDTO) (noteUUID string, err error) {
	dto.GenerateShortBody()
	noteUUID, err = s.storage.Create(ctx, dto)

	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return noteUUID, err
		}
		return noteUUID, fmt.Errorf("failed to create note. error: %w", err)
	}

	return noteUUID, nil
}

func (s service) GetOne(ctx context.Context, uuid string) (n Note, err error) {
	n, err = s.storage.FindOne(ctx, uuid)

	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return n, err
		}
		return n, fmt.Errorf("failed to find note by uuid. error: %w", err)
	}
	return n, nil
}

func (s service) GetByCategoryUUID(ctx context.Context, uuid string) (notes []Note, err error) {
	notes, err = s.storage.FindByCategoryUUID(ctx, uuid)

	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return notes, err
		}
		return notes, fmt.Errorf("failed to get notes by ids. error: %w", err)
	}
	if len(notes) == 0 {
		return notes, apperror.ErrNotFound
	}
	return notes, nil
}

func (s service) Update(ctx context.Context, uuid string, dto UpdateNoteDTO, tagsUpdate bool) error {
	if dto.Body == "" && dto.Header == "" && dto.CategoryUUID == "" && !tagsUpdate {
		return apperror.BadRequestError("nothing to update")
	}
	err := s.storage.Update(ctx, uuid, dto, tagsUpdate)

	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to update note. error: %w", err)
	}
	return nil
}

func (s service) Delete(ctx context.Context, uuid string) error {
	err := s.storage.Delete(ctx, uuid)

	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to delete note. error: %w", err)
	}
	return err
}
