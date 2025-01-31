package repositoryPromoCode

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"telemetry-sale/internal/config/storage"
	"telemetry-sale/internal/util/helper"
)

type PromoGenerationRepository struct {
	db  *pgxpool.Pool
	ctx *gin.Context
	tr  *helper.Transaction
}

func NewPromoGenerationRepository(ctx *gin.Context) *PromoGenerationRepository {
	return &PromoGenerationRepository{
		db:  storage.GetDB(),
		ctx: ctx,
	}
}

func (r *PromoGenerationRepository) GetByCodeAndOrganizationId(code string, organizationId uint64) (uint64, error) {
	var id uint64

	query := `select id from promo_codes where code=$1 and organization_id=$2`

	err := r.db.QueryRow(r.ctx, query, code, organizationId).Scan(&id)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return id, err
	}

	return id, nil
}
