package dtos

import (
	"bytes"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/reitmas32/rkit/core/customctx"
	"github.com/reitmas32/rkit/core/kerrors"
	"github.com/reitmas32/rkit/core/result"
	basehttp "github.com/reitmas32/rkit/infrastructure/http"

	"github.com/gin-gonic/gin"
)

type DTO interface {
	Validate() error
}

type ErrorDTO struct {
	kerrors.KError
}

func GetDTO[K DTO](ctx *gin.Context, cc *customctx.CustomContext) *K {

	entry := cc.Logger()

	var dto K
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		entry.Error(err.Error())

		response := makeResponseError(err, dto)

		cc.NewError(kerrors.NewKError(response.Error.Error(), http.StatusUnprocessableEntity, map[string]any{"scope": "dto.validate." + reflect.TypeOf(dto).Name()}))
		ctx.JSON(response.StatusCode, response.ToMap())
		return nil
	}
	if err := dto.Validate(); err != nil {
		entry.Error(err.Error())

		response := makeResponseError(err, dto)

		cc.NewError(kerrors.NewKError(response.Error.Error(), http.StatusUnprocessableEntity, map[string]any{"scope": "dto.validate." + reflect.TypeOf(dto).Name()}))
		//ctx.JSON(response.StatusCode, response.ToMapWithCustomContext(cc))
		return nil
	}
	return &dto
}

func GetDTOWithResponse[K DTO](ctx *gin.Context, cc *customctx.CustomContext) basehttp.Response[K] {

	entry := cc.Logger()

	var dto K
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		entry.Error(err.Error())

		response := makeResponseError(err, dto)

		cc.NewError(kerrors.NewKError(response.Error.Error(), http.StatusUnprocessableEntity, map[string]any{"scope": "dto.validate." + reflect.TypeOf(dto).Name()}))
		return response
	}
	if err := dto.Validate(); err != nil {
		entry.Error(err.Error())

		response := makeResponseError(err, dto)

		cc.NewError(kerrors.NewKError(response.Error.Error(), http.StatusUnprocessableEntity, map[string]any{"scope": "dto.validate." + reflect.TypeOf(dto).Name()}))
		//ctx.JSON(response.StatusCode, response.ToMapWithCustomContext(cc))
		return response
	}
	return basehttp.Response[K]{
		Data:       dto,
		StatusCode: http.StatusOK,
		Success:    true,
	}
}

func makeResponseError[K DTO](err error, dto K) basehttp.Response[K] {

	typ := reflect.TypeOf(dto)

	// Si es un puntero, obten el elemento apuntado
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	response := basehttp.Response[K]{
		Error: &basehttp.ErrorInfo{
			Code:    http.StatusUnprocessableEntity,
			Message: err.Error(),
			Scope:   "dto.validate." + typ.Name(),
		},
		StatusCode: http.StatusUnprocessableEntity,
		Success:    false,
	}

	return response
}

func GetAuthToken(ctx *gin.Context, cc *customctx.CustomContext) result.Result[string] {

	entry := cc.Logger()

	token := ctx.GetHeader("Authorization")

	if token == "" {
		entry.Error("No se encontró un token de autenticación")
		return result.NewErrResult[string](kerrors.NewKError("No se encontró un token de autenticación", http.StatusUnauthorized, nil))
	}

	token = strings.TrimPrefix(token, "Bearer ")

	if token == "" {
		entry.Error("No se encontró un token de autenticación")
		return result.NewErrResult[string](kerrors.NewKError("No se encontró un token de autenticación", http.StatusUnauthorized, nil))
	}

	return result.NewOkResult[string](token)
}

func GetAuthTokenWithEarlyResponse(ctx *gin.Context, cc *customctx.CustomContext) result.Result[string] {

	entry := cc.Logger()

	token := ctx.GetHeader("Authorization")

	if token == "" {
		entry.Error("Not Found Authorization Header")

		err := kerrors.NewKError("Not Found Authorization Header", http.StatusUnauthorized, nil)
		cc.NewError(err)

		response := basehttp.Response[string]{
			StatusCode: http.StatusUnauthorized,
		}

		ctx.JSON(response.StatusCode, response.ToMapWithCustomContext(cc))
		ctx.Abort()
		return result.NewErrResult[string](err)
	}

	token = strings.TrimPrefix(token, "Bearer ")

	if token == "" {
		entry.Error("Not Found Token in Authorization Header")

		err := kerrors.NewKError("Not Found Token in Authorization Header", http.StatusUnauthorized, nil)
		cc.NewError(err)

		response := basehttp.Response[string]{
			StatusCode: http.StatusUnauthorized,
		}
		ctx.JSON(response.StatusCode, response.ToMap())
		ctx.Abort()
		return result.NewErrResult[string](kerrors.NewKError("Not Found Token in Authorization Header", http.StatusUnauthorized, nil))
	}

	return result.NewOkResult[string](token)
}

// GetDTONonDestructive extrae un DTO del contexto sin modificar el estado del request body,
// permitiendo que se pueda volver a extraer el DTO en otro lugar
func GetDTONonDestructive[K DTO](ctx *gin.Context, cc *customctx.CustomContext) result.Result[*K] {
	entry := cc.Logger()

	// Leer el body completo sin consumirlo
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		entry.Error("error reading request body", err)
		return result.Err[*K](kerrors.NewKError(err.Error(), http.StatusUnprocessableEntity, nil))
	}

	// Restaurar el body para que pueda ser leído nuevamente
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Parsear el JSON en el DTO
	var dto K
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		entry.Error("error binding JSON", err)
		// Restaurar el body antes de retornar el error
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		return result.Err[*K](kerrors.NewKError(err.Error(), http.StatusUnprocessableEntity, map[string]any{"scope": "dto.validate." + reflect.TypeOf(dto).Name()}))
	}

	// Validar el DTO
	if err := dto.Validate(); err != nil {
		entry.Error("error validating DTO", err)
		// Restaurar el body antes de retornar el error
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		return result.Err[*K](kerrors.NewKError(err.Error(), http.StatusUnprocessableEntity, map[string]any{"scope": "dto.validate." + reflect.TypeOf(dto).Name()}))
	}

	// Restaurar el body para que pueda ser leído nuevamente
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return result.NewOkResult[*K](&dto)
}
