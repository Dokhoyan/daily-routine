package sprint

import (
	"context"
	"fmt"
)

func (s *serv) Delete(ctx context.Context, id int64) error {
	if ctx == nil {
		ctx = context.Background()
	}

	if err := s.sprintRepo.DeleteSprint(ctx, id); err != nil {
		return fmt.Errorf("failed to delete sprint: %w", err)
	}

	return nil
}

