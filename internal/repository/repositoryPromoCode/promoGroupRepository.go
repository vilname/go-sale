package repositoryPromoCode

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"telemetry-sale/internal/config/storage"
	"telemetry-sale/internal/model"
)

type PromoGroupRepository struct {
	db  *pgxpool.Pool
	ctx *gin.Context
}

func NewPromoGroupRepository(ctx *gin.Context) *PromoGroupRepository {
	return &PromoGroupRepository{
		db:  storage.GetDB(),
		ctx: ctx,
	}
}

func (repository PromoGroupRepository) SavePromoGroup(name string, organizationId uint64) (model.IdResult, error) {
	var result model.IdResult
	var groupId uint64

	query := `insert into groups (
					name,
					organization_id,
					created
				) values ($1, $2, NOW()) returning id`

	row := repository.db.QueryRow(repository.ctx, query, name, organizationId)

	if err := row.Scan(&groupId); err != nil {
		return result, err
	}

	result.Id = groupId

	return result, nil
}

func (repository PromoGroupRepository) GetQtyElementGroup(id uint64) (uint16, error) {
	var qtyElement uint16

	query := `select count(*) from promo_codes where group_id=$1`
	err := repository.db.QueryRow(repository.ctx, query, id).Scan(&qtyElement)
	if err != nil {
		return qtyElement, err
	}

	return qtyElement, nil
}

func (repository PromoGroupRepository) GetList(organizationId uint64) ([]model.PromoGroupResult, error) {
	var promoGroupsResult []model.PromoGroupResult

	query := `select id, name from groups where organization_id=$1`
	rows, err := repository.db.Query(repository.ctx, query, organizationId)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var promoGroupResult model.PromoGroupResult

		rows.Scan(&promoGroupResult.Id, &promoGroupResult.Name)

		promoGroupsResult = append(promoGroupsResult, promoGroupResult)
	}

	rows.Close()

	return promoGroupsResult, nil
}
