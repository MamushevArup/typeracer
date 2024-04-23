package admin

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (r *repo) ModerationById(ctx context.Context, modId uuid.UUID) (models.ModerationApprove, error) {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	details := models.ModerationApprove{}
	details.ModerationID = modId

	sql, args, err := sq.
		Select("racer_id", "content", "author", "length", "source", "source_title", "sent_at").
		From("moderation").
		Where(squirrel.Eq{"moderation_id": modId}).
		ToSql()

	if err != nil {
		r.lg.Errorf("can't build query for select content details, mod_id=%v err=%v", modId, err)
		return details, fmt.Errorf("fail build query err=%w", err)
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(&details.RacerID, &details.Content, &details.Author,
		&details.Length, &details.Source, &details.SourceTitle, &details.SentAt)

	if err != nil {
		if err.Error() == "no rows in result set" {
			r.lg.Errorf("content not found, mod_id=%v", modId)
			return details, fmt.Errorf("content not found")
		}
		r.lg.Errorf("can't select content details, mod_id=%v err=%v", modId, err)
		return details, fmt.Errorf("fail select content details err=%w", err)
	}

	return details, nil

}

func (r *repo) ApproveContent(ctx context.Context, transaction models.TextAcceptTransaction) error {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	b, err := r.db.Begin(ctx)
	if err != nil {
		r.lg.Errorf("can't start transaction, user_id=%v err=%v", transaction.Text.ContributorID, err)
		return fmt.Errorf("fail start transaction err=%w", err)

	}

	approve := transaction.Text

	sql, args, err := sq.
		Insert("text").
		Columns("id", "contributor_id", "content", "author", "length", "source", "source_title", "accepted_at").
		Values(approve.TextID, approve.ContributorID, approve.Content, approve.Author, approve.Length, approve.Source, approve.SourceTitle, approve.AcceptedAt).
		ToSql()

	if err != nil {
		r.lg.Errorf("can't build query for insert text, err=%v", err)
		return fmt.Errorf("fail build query err=%w", err)
	}

	_, err = b.Exec(ctx, sql, args...)
	if err != nil {
		r.lg.Errorf("can't insert text, text_id=%v, user_id=%v, err=%v", approve.TextID, approve.ContributorID, err)
		return fmt.Errorf("fail insert text err=%w", err)
	}

	cont := transaction.Contributor

	sql, args, err = sq.
		Insert("contributor").
		Columns("user_id", "text_id", "sent_at").
		Values(cont.ContributorID, cont.TextID, cont.SentAt).
		ToSql()

	if err != nil {
		r.lg.Errorf("can't build query for insert contributor, err=%v", err)
		return fmt.Errorf("fail build query err=%w", err)
	}

	_, err = b.Exec(ctx, sql, args...)
	if err != nil {
		r.lg.Errorf("can't insert contributor, user_id=%v, text_id=%, err=%v", cont.ContributorID, cont.TextID, err)
		return fmt.Errorf("fail insert contributor err=%w", err)
	}

	err = b.Commit(ctx)
	if err != nil {
		r.lg.Errorf("can't commit transaction, err=%v", err)
		return fmt.Errorf("fail commit transaction err=%w", err)
	}

	return nil
}
