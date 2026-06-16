package users_transport_http

import (
	"net/http"

	"github.com/Tim73916/go-todoapp/internal/core/domain"
	core_logger "github.com/Tim73916/go-todoapp/internal/core/logger"
	core_http_request "github.com/Tim73916/go-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Tim73916/go-todoapp/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=50" example:"Brandy Werson"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+" example:"+131349581245"`
}

type CreateUserResponse UserDTOResponse

// CreateUser 		godoc
// @Summary 		Create a user
// @Description 	Create a new user in the system
// @Tags 			users
// @Accept 			json
// @Produce 		json
// @Param 			request body CreateUserRequest true "CreateUser requset body"
// @Success 		201 {object} CreateUserResponse "Successfully created user"
// @Failure 		400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 		500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 			/users [post]
func (h *UsersHTTPHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPReponseHandler(log, rw)

	var request CreateUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")

		return
	}
	userDomain := domainFromDTO(request)
	userDomain, err := h.usersService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")

		return
	}

	response := CreateUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(dto CreateUserRequest) domain.User {
	return domain.NewUserUninitialized(dto.FullName, dto.PhoneNumber)
}
