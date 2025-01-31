package helper

import (
	"fmt"
	"log/slog"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func CreateArgsForIn(qtyArgs int) string {
	var placeholders []string

	for i := 0; i < qtyArgs; i++ {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1)) // начинаем с $1
	}

	return strings.Join(placeholders, ", ")
}

func GeneratePromoCode(n int) string {
	const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var sb strings.Builder
	sb.Grow(n)

	rand.NewSource(time.Now().UnixNano())

	// Выбираем случайные символы из letterBytes.
	for i := 0; i < n; i++ {
		sb.WriteByte(letterBytes[rand.Intn(len(letterBytes))])
	}

	return sb.String()
}

func GetIdsSlice(idsString string) []uint64 {
	ids := make([]uint64, 0, 100)

	for _, id := range strings.Split(idsString, ",") {
		parseUint, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			slog.Error("GetIdsSlice", err.Error())
		}

		ids = append(ids, parseUint)
	}

	return ids
}
