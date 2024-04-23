package admin

import (
	"context"
	"fmt"
	"github.com/MamushevArup/typeracer/internal/models"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (r *repo) SelectModeration(ctx context.Context, limit, offset int, sort string) ([]models.ModerationRepoResponse, error) {

	var resp []models.ModerationRepoResponse

	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.
		Select("moderation_id", "sent_at", "username").
		From("moderation").
		Join("racer on moderation.racer_id=racer.id").
		Where(squirrel.Eq{"moderation.status": 0}).
		OrderBy("sent_at " + sort).
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()

	if err != nil {
		r.lg.Errorf("can't build query for select content, err=%v", err)
		return resp, fmt.Errorf("fail build query err=%w", err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		r.lg.Errorf("can't select content, err=%v", err)
		return resp, fmt.Errorf("fail select content err=%w", err)
	}

	for rows.Next() {
		var currResponse models.ModerationRepoResponse
		err = rows.Scan(&currResponse.ModerationID, &currResponse.SentAt, &currResponse.Username)
		if err != nil {
			r.lg.Errorf("can't scan content, err=%v", err)
			return resp, fmt.Errorf("fail scan content err=%w", err)
		}
		resp = append(resp, currResponse)
	}

	defer rows.Close()

	if len(resp) == 0 {
		return resp, fmt.Errorf("no content to moderate")
	}

	return resp, nil
}

func (r *repo) ContentDetails(ctx context.Context, modId uuid.UUID) (models.ModerationTextDetails, error) {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	details := models.ModerationTextDetails{}

	sql, args, err := sq.
		Select("content", "author", "source", "source_title").
		From("moderation").
		Where(squirrel.Eq{"moderation_id": modId}).
		ToSql()

	if err != nil {
		r.lg.Errorf("can't build query for select content details, mod_id=%v err=%v", modId, err)
		return details, fmt.Errorf("fail build query err=%w", err)
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(&details.Content, &details.Author, &details.Source, &details.SourceTitle)
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

func (r *repo) DeleteModerationText(ctx context.Context, modId uuid.UUID) error {
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	sql, args, err := sq.
		Delete("moderation").
		Where(squirrel.Eq{"moderation_id": modId}).
		ToSql()

	if err != nil {
		r.lg.Errorf("can't build query for delete moderation text, mod_id=%v err=%v", modId, err)
		return fmt.Errorf("fail build query err=%w", err)
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		r.lg.Errorf("can't delete moderation text, mod_id=%v err=%v", modId, err)
		return fmt.Errorf("fail delete moderation text err=%w", err)
	}

	return nil
}
