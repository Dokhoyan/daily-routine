package habit

import (
	"context"
	"fmt"
)

func (s *serv) Delete(ctx context.Context, id int64) error {
	if err := s.habitRepo.DeleteHabit(ctx, id); err != nil {
		return fmt.Errorf("failed to delete habit: %w", err)
	}
	return nil
}
