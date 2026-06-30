package helpers

import "github.com/google/uuid"

func ScanNullableUUID(src uuid.NullUUID) *uuid.UUID {
	if src.Valid {
		return &src.UUID
	}
	return nil
}
