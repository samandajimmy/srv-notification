package svcError

import (
	"github.com/nbs-go/errx"
	"net/http"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp"
)

var b = errx.NewBuilder("srv-notification", errx.FallbackError(
	errx.NewError("500", "An error has occurred, please try again later",
		nhttp.WithStatus(http.StatusInternalServerError),
	),
))

var ResourceNotFound = b.NewError("E_COMM_1", "Resource not found",
	nhttp.WithStatus(http.StatusNotFound))

var StaleResource = b.NewError("E_COMM_2", "Cannot update stale resource",
	nhttp.WithStatus(http.StatusConflict))

var DuplicatedResource = b.NewError("E_COMM_6", "Duplicated resource",
	nhttp.WithStatus(http.StatusConflict))
