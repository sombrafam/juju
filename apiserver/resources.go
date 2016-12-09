// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package apiserver

import (
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/juju/errors"
	charmresource "gopkg.in/juju/charm.v6-unstable/resource"

	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/resource"
	"github.com/juju/juju/state"
)

// resourceUploadHandler handles resources uploads for model migrations.
type resourceUploadHandler struct {
	ctxt          httpContext
	stateAuthFunc func(*http.Request) (*state.State, error)
}

func (h *resourceUploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Validate before authenticate because the authentication is dependent
	// on the state connection that is determined during the validation.
	st, err := h.stateAuthFunc(r)
	if err != nil {
		if err := sendError(w, err); err != nil {
			logger.Errorf("%v", err)
		}
		return
	}
	defer h.ctxt.release(st)

	switch r.Method {
	case "POST":
		res, err := h.processPost(r, st)
		if err != nil {
			if err := sendError(w, err); err != nil {
				logger.Errorf("%v", err)
			}
			return
		}
		if err := sendStatusAndJSON(w, http.StatusOK, &params.ResourceUploadResult{
			ID:        res.ID,
			Timestamp: res.Timestamp,
		}); err != nil {
			logger.Errorf("%v", err)
		}
	default:
		if err := sendError(w, errors.MethodNotAllowedf("unsupported method: %q", r.Method)); err != nil {
			logger.Errorf("%v", err)
		}
	}
}

// processPost handles resources upload POST request after
// authentication.
func (h *resourceUploadHandler) processPost(r *http.Request, st *state.State) (resource.Resource, error) {
	var empty resource.Resource
	query := r.URL.Query()

	target, isUnit, err := getUploadTarget(query)
	if err != nil {
		return empty, errors.Trace(err)
	}

	userID := query.Get("user") // Is allowed to be blank
	res, err := queryToResource(query)
	if err != nil {
		return empty, errors.Trace(err)
	}
	rSt, err := st.Resources()
	if err != nil {
		return empty, errors.Trace(err)
	}

	outRes, err := setResource(isUnit, target, userID, res, r.Body, rSt)
	if err != nil {
		return empty, errors.Annotate(err, "resource upload failed")
	}
	return outRes, nil
}

func setResource(isUnit bool, target, user string, res charmresource.Resource, r io.Reader, rSt state.Resources) (
	resource.Resource, error,
) {
	if isUnit {
		return rSt.SetUnitResource(target, user, res)
	}
	return rSt.SetResource(target, user, res, r)

}

func getUploadTarget(query url.Values) (string, bool, error) {
	appName := query.Get("application")
	unitName := query.Get("unit")
	switch {
	case appName == "" && unitName == "":
		return "", false, errors.BadRequestf("missing application/unit")
	case appName != "" && unitName != "":
		return "", false, errors.BadRequestf("application and unit can't be set at the same time")
	case appName != "":
		return appName, false, nil
	default:
		return unitName, true, nil
	}
}

func queryToResource(query url.Values) (charmresource.Resource, error) {
	var err error
	empty := charmresource.Resource{}

	res := charmresource.Resource{
		Meta: charmresource.Meta{
			Name:        query.Get("name"),
			Path:        query.Get("path"),
			Description: query.Get("description"),
		},
	}
	if res.Name == "" {
		return empty, errors.BadRequestf("missing name")
	}
	if res.Path == "" {
		return empty, errors.BadRequestf("missing path")
	}
	if res.Description == "" {
		return empty, errors.BadRequestf("missing description")
	}
	res.Type, err = charmresource.ParseType(query.Get("type"))
	if err != nil {
		return empty, errors.BadRequestf("invalid type")
	}
	res.Origin, err = charmresource.ParseOrigin(query.Get("origin"))
	if err != nil {
		return empty, errors.BadRequestf("invalid origin")
	}
	res.Revision, err = strconv.Atoi(query.Get("revision"))
	if err != nil {
		return empty, errors.BadRequestf("invalid revision")
	}
	res.Size, err = strconv.ParseInt(query.Get("size"), 10, 64)
	if err != nil {
		return empty, errors.BadRequestf("invalid size")
	}
	res.Fingerprint, err = charmresource.ParseFingerprint(query.Get("fingerprint"))
	if err != nil {
		return empty, errors.BadRequestf("invalid fingerprint")
	}
	return res, nil
}
