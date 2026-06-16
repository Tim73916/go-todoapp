package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Tim73916/go-todoapp/internal/core/domain"
	core_logger "github.com/Tim73916/go-todoapp/internal/core/logger"
	core_http_request "github.com/Tim73916/go-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Tim73916/go-todoapp/internal/core/transport/http/response"
	core_http_types "github.com/Tim73916/go-todoapp/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name" swaggertype:"string" example:"Brandy Werson"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number" swaggertype:"string" example:"+131349581245"`
}

func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("`FullName` can't be NULL")
		}

		fullNameLen := len([]rune(*r.FullName.Value))
		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("`FullName` must be between 3 and 100 symbols")
		}
	}

	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			phoneNumberLen := len([]rune(*r.PhoneNumber.Value))
			if phoneNumberLen < 10 || phoneNumberLen > 15 {
				return fmt.Errorf("`PhoneNumber` must be between 10 and 15 symbols")
			}

			if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
				return fmt.Errorf("`PhoneNumber` must start with '+' symbol")
			}
		}
	}

	return nil
}

type PatchUserResponse UserDTOResponse

// PatchUser godoc
// @Summary Update user
// @Description Update information of an existing user in the system
// @Description ### Field update logic (Three-state logic):
// @Description 1. *Field not provided*: `phone_number` is ignored, database value does not change
// @Description 2. *Value explicitly provided*: `"phone_number": "+711122233344"` – sets new phone number in DB
// @Description 3. *Null explicitly provided*: `"phone_number": null` – clears the field in DB (sets to NULL)
// @Description Restrictions: `full_name` cannot be set to null
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body PatchUserRequest true "Patch user request body"
// @Success 200 {object} PatchUserResponse "User successfully updated"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 404 {object} core_http_response.ErrorResponse "User not found"
// @Failure 409 {object} core_http_response.ErrorResponse "Conflict"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /users/{id} [patch]
func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPReponseHandler(log, rw)

	userID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID path value",
		)

		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err, "failed to decode an validate HTTP request",
		)

		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch user",
		)

		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusOK)

	log.Debug(
		fmt.Sprintf(
			"PatchUserRequest fields: \nFullName: '%v'\nPhoneNumber: '%v`'",
			request.FullName,
			request.PhoneNumber,
		),
	)
	rw.WriteHeader(http.StatusOK)
}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.FullName.ToDomain(),
		request.PhoneNumber.ToDomain(),
	)
}
