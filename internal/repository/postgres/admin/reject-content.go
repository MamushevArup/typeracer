package admin

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/Masterminds/squirrel"
)

func (r *repo) RejectContent(ctx context.Context, reject models.ModerationRejectToRepo) error {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.
		Update("moderation").
		Set("status", -1).
		Where(squirrel.Eq{"moderation_id": reject.ModerationID}).
		ToSql()

	if err != nil {
		r.lg.Errorf("can't build query for reject content, mod_id=%v err=%v", reject.ModerationID, err)
		return fmt.Errorf("fail build query err=%w", err)
	}

	b, err := r.db.Begin(ctx)
	if err != nil {
		r.lg.Errorf("can't start transaction, mod_id=%v err=%v", reject.ModerationID, err)
		return fmt.Errorf("fail start transaction err=%w", err)
	}

	upd, err := b.Exec(ctx, sql, args...)
	if err != nil {
		r.lg.Errorf("can't reject content, mod_id=%v err=%v", reject.ModerationID, err)
		return fmt.Errorf("fail reject content err=%w", err)
	}

	if !upd.Update() {
		r.lg.Errorf("unable update content, mod_id=%v", reject.ModerationID)
		return fmt.Errorf("error with updating content")
	}

	sql, args, err = sq.
		Insert("rejected").
		Columns("moderation_id", "response").
		Values(reject.ModerationID, reject.Reason).
		ToSql()

	if err != nil {
		r.lg.Errorf("can't build query for insert rejected, mod_id=%v err=%v", reject.ModerationID, err)
		return fmt.Errorf("fail build query err=%w", err)
	}

	inst, err := b.Exec(ctx, sql, args...)
	if err != nil {
		r.lg.Errorf("can't insert rejected, mod_id=%v err=%v", reject.ModerationID, err)
		return fmt.Errorf("fail insert rejected err=%w", err)
	}

	if !inst.Insert() {
		r.lg.Errorf("unable insert rejected, mod_id=%v", reject.ModerationID)
		return fmt.Errorf("error with inserting rejected")
	}

	err = b.Commit(ctx)
	if err != nil {
		r.lg.Errorf("can't commit transaction, mod_id=%v err=%v", reject.ModerationID, err)
		return fmt.Errorf("fail commit transaction err=%w", err)
	}

	return nil
}
