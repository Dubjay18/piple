package helpers

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// TimeToPgTimestamp converts an optional Go time into the pgtype.Timestamp
// sqlc generates for nullable timestamp columns.
func TimeToPgTimestamp(t *time.Time) pgtype.Timestamp {
	if t == nil {
		return pgtype.Timestamp{}
	}
	return pgtype.Timestamp{Time: *t, Valid: true}
}

// PgTimestampToTime is the inverse of TimeToPgTimestamp.
func PgTimestampToTime(t pgtype.Timestamp) *time.Time {
	if !t.Valid {
		return nil
	}
	value := t.Time
	return &value
}

// UUIDToPgUUID converts an optional uuid.UUID into the pgtype.UUID sqlc
// generates for nullable uuid columns.
func UUIDToPgUUID(id *uuid.UUID) pgtype.UUID {
	if id == nil {
		return pgtype.UUID{}
	}
	return pgtype.UUID{Bytes: [16]byte(*id), Valid: true}
}
