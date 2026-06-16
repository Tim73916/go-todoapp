package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Tim73916/go-todoapp/internal/core/domain"
	core_logger "github.com/Tim73916/go-todoapp/internal/core/logger"
	core_http_request "github.com/Tim73916/go-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Tim73916/go-todoapp/internal/core/transport/http/response"
)

type GetStatisticsResponse struct {
	TasksCreated                int      `json:"tasks_created" example:"50"`
	TasksCompleted              int      `json:"tasks_completed" example:"13"`
	TasksCompletedRate          *float64 `json:"tasks_completeed_rate" example:"67"`
	TasksAveragweCompletionTime *string  `json:"tasks_average_completion_time" example:"12m45s"`
}

// GetStatistics godoc
// @Summary Get statistics
// @Description Get task statistics with optional filtering by user_id and/or time range
// @Tags statistics
// @Produce json
// @Param user_id query int false "Filter statistics by specific user"
// @Param from query string false "Start of the statistics period (inclusive), format: YYYY-MM-DD"
// @Param to query string false "End of the statistics period (exclusive), format: YYYY-MM-DD"
// @Success 200 {object} GetStatisticsResponse "Statistics successfully retrieved"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /statistics [get]
func (h *StatisticsHTTPHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPReponseHandler(log, rw)

	userID, from, to, err := getUserIDFromToQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID/from/to query params",
		)

		return
	}

	statistics, err := h.statisticsService.GetStatistics(ctx, userID, from, to)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get statistics",
		)

		return
	}

	response := toDTOFromDomain(statistics)

	responseHandler.JSONResponse(response, http.StatusOK)

}

func toDTOFromDomain(statistics domain.Statistics) GetStatisticsResponse {
	var avgTime *string
	if statistics.TasksAveragweCompletionTime != nil {
		duration := statistics.TasksAveragweCompletionTime.String()
		avgTime = &duration
	}

	return GetStatisticsResponse{
		TasksCreated:                statistics.TasksCreated,
		TasksCompleted:              statistics.TasksCompleted,
		TasksCompletedRate:          statistics.TasksCompletedRate,
		TasksAveragweCompletionTime: avgTime,
	}
}

func getUserIDFromToQueryParams(r *http.Request) (*int, *time.Time, *time.Time, error) {
	const (
		userIDQueryParamKey = "user_id"
		fromQueryParamKey   = "from"
		totQueryParamKey    = "to"
	)

	userID, err := core_http_request.GetIntQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}

	from, err := core_http_request.GetDateQueryParam(r, fromQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'from' query param: %w", err)
	}

	to, err := core_http_request.GetDateQueryParam(r, totQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'to' query param: %w", err)
	}

	return userID, from, to, nil
}
