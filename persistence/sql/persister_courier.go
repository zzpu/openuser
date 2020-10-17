package sql

import (
	"context"
	"database/sql"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/ory/x/sqlcon"

	"github.com/zzpu/openuser/courier"
)

var _ courier.Persister = new(Persister)

func (p *Persister) AddMessage(ctx context.Context, m *courier.Message) error {
	m.Status = courier.MessageStatusQueued
	return sqlcon.HandleError(p.GetConnection(ctx).Create(m)) // do not create eager to avoid identity injection.
}

func (p *Persister) NextMessages(ctx context.Context, limit uint8) ([]courier.Message, error) {
	var m []courier.Message
	if err := p.GetConnection(ctx).
		Eager().
		Where("status != ?", courier.MessageStatusSent).
		Order("created_at ASC").Limit(int(limit)).All(&m); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, errors.WithStack(courier.ErrQueueEmpty)
		}
		return nil, sqlcon.HandleError(err)
	}

	if len(m) == 0 {
		return nil, errors.WithStack(courier.ErrQueueEmpty)
	}

	return m, nil
}

func (p *Persister) LatestQueuedMessage(ctx context.Context) (*courier.Message, error) {
	var m courier.Message
	if err := p.GetConnection(ctx).
		Eager().
		Where("status != ?", courier.MessageStatusSent).
		Order("created_at DESC").First(&m); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, errors.WithStack(courier.ErrQueueEmpty)
		}
		return nil, sqlcon.HandleError(err)
	}

	return &m, nil
}

func (p *Persister) SetMessageStatus(ctx context.Context, id uuid.UUID, ms courier.MessageStatus) error {
	count, err := p.GetConnection(ctx).RawQuery("UPDATE courier_messages SET status = ? WHERE id = ?", ms, id).ExecWithCount()
	if err != nil {
		return sqlcon.HandleError(err)
	}

	if count == 0 {
		return errors.WithStack(sqlcon.ErrNoRows)
	}

	return nil
}
