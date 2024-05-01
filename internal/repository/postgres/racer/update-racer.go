package racer

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/Masterminds/squirrel"
	"log"
)

func (r *repo) UpdateRacer(ctx context.Context, racer models.RacerUpdateRepo) error {
	updateCond := make(map[string]interface{}, 2)
	if racer.Email != "" {
		updateCond["email"] = racer.Email
	}
	if racer.Username != "" {
		updateCond["username"] = racer.Username
	}

	sql, args, err := squirrel.Update("racer").
		SetMap(updateCond).
		Where(squirrel.Eq{"id": racer.Id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		r.lg.Errorf("can't build query %v", err)
		return ErrQueryBuild
	}

	flag, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		r.lg.Errorf("can't update racer id=%v, err=%v", racer.Id, err)
		return fmt.Errorf("can't update racer %w", err)
	}
	if flag.RowsAffected() == 0 {
		r.lg.Errorf("can't update racer id=%v, err=%v", racer.Id, err)
		return ErrUpdate
	}

	log.Println("Racer updated successfully")
	return nil

}
