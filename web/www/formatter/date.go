package formatter

import (
	"fmt"
	"time"
)

// AsDashedDate returns the "YYYY-mm-dd" format of a time object.
func AsDashedDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}
