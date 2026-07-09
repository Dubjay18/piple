package helpers

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ParseUUID(raw string) (id pgtype.UUID, ok bool) {
	if err := id.Scan(raw); err != nil {
		return id, false
	}
	return id, true
}

func UUIDToString(id pgtype.UUID) string {
	s, _ := id.Value()
	if str, ok := s.(string); ok {
		return str
	}
	return ""
}

func TimestampToTimePtr(ts pgtype.Timestamp) *time.Time {
	if !ts.Valid {
		return nil
	}
	t := ts.Time
	return &t
}