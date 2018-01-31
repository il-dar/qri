package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/qri-io/qri/repo"
)

// QriCtxKey defines a distinct type for
// keys for context values should always use custom
// types to avoid collissions.
// see comment on context.WithValue for more info
type QriCtxKey string

// DatasetRefCtxKey is the key for adding a dataset reference
// to a context.Context
const DatasetRefCtxKey QriCtxKey = "datasetRef"

// DatasetRefFromReq examines the path element of a request URL
// to
func DatasetRefFromReq(r *http.Request) (*repo.DatasetRef, error) {
	if r.URL.String() == "" || r.URL.Path == "" {
		return nil, nil
	}
	return DatasetRefFromPath(r.URL.Path)
}

func DatasetRefFromPath(path string) (*repo.DatasetRef, error) {
	refstr := HttpPathToQriPath(path)
	ref, err := repo.ParseDatasetRef(refstr)
	if err != nil {
		ref = nil
	}
	return ref, err
}

// DatasetRefFromCtx extracts a Dataset reference from a given
// context if one is set, returning nil otherwise
func DatasetRefFromCtx(ctx context.Context) *repo.DatasetRef {
	iface := ctx.Value(DatasetRefCtxKey)
	if ref, ok := iface.(*repo.DatasetRef); ok {
		return ref
	}
	return nil
}

// HttpPathToQriPath converts a http path to a
// qri path
func HttpPathToQriPath(path string) string {
	paramIndex := strings.Index(path, "?")
	if paramIndex != -1 {
		path = path[:paramIndex]
	}
	path = strings.Replace(path, "/at/ipfs/", "@", 1)
	path = strings.Replace(path, "/at/", "@", 1)
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")
	return path
}
