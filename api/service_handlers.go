package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	db "github.com/JekaTka/services-api/db/sqlc"
)

type getServicesRequest struct {
	Search string `form:"search"`
	SortBy string `form:"sort_by" binding:"required,service_sort"`
	Limit  int32  `form:"limit" binding:"required,min=1"`
	Offset int32  `form:"offset" binding:"min=0"`
}

type getServicesMetadata struct {
	TotalItems int64 `json:"total_items"`
}

type getServicesItem struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Versions    int64     `json:"versions"`
}

type getServicesResponse struct {
	Metadata getServicesMetadata `json:"metadata"`
	Services []getServicesItem   `json:"services"`
}

type sqlcService struct {
	db.Service
}

func (srv sqlcService) toResponse() getServicesItem {
	return getServicesItem{
		ID:          srv.ID,
		Name:        srv.Name,
		Description: srv.Description,
		CreatedAt:   srv.CreatedAt.Time,
		UpdatedAt:   srv.UpdatedAt.Time,
	}
}

func (server *Server) getServices(ctx *gin.Context) {
	var req getServicesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListServicesParams{
		Limit:   req.Limit,
		Offset:  req.Offset,
		OrderBy: req.SortBy,
		Search:  req.Search,
	}

	services, err := server.store.ListServices(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	servicesResponse := make([]getServicesItem, 0, len(services))
	for _, srv := range services {
		serviceModel := sqlcService{srv}
		serviceResponse := serviceModel.toResponse()
		versionCount, err := server.store.GetServiceVersionsCount(ctx, srv.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		serviceResponse.Versions = versionCount
		servicesResponse = append(servicesResponse, serviceResponse)
	}

	totalCount, err := server.store.GetServicesCount(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := getServicesResponse{
		Metadata: getServicesMetadata{totalCount},
		Services: servicesResponse,
	}

	ctx.JSON(http.StatusOK, response)
}

type getServiceVersionsRequest struct {
	ID string `uri:"id" binding:"required"`
}

func (server *Server) getServiceVersions(ctx *gin.Context) {
	var req getServiceVersionsRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	serviceVersions, err := server.store.GetVersionsByServiceID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, serviceVersions)
}
