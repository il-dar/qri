package fsrepo

import (
	"encoding/json"
	"github.com/ipfs/go-datastore"
	"io/ioutil"
	"os"
	"path/filepath"
)

type basepath string

func (bp basepath) filepath(f File) string {
	return filepath.Join(string(bp), Filepath(f))
}

func (bp basepath) readBytes(f File) ([]byte, error) {
	return ioutil.ReadFile(bp.filepath(f))
}

func (bp basepath) saveFile(d interface{}, f File) error {
	data, err := json.Marshal(d)
	if err != nil {
		log.Debug(err.Error())
		return err
	}
	return ioutil.WriteFile(bp.filepath(f), data, os.ModePerm)
}

// File represents a type file in a qri repository
type File int

const (
	// FileUnknown makes the default file value invalid
	FileUnknown File = iota
	// FileLockfile is the on-disk mutex lock
	FileLockfile
	// FileInfo stores information about this repository
	// like version number, size of repo, etc.
	FileInfo
	// FileConfig holds configuration specific to this repo
	FileConfig
	// FileDatasets holds the list of datasets
	FileDatasets
	// FileEventLogs is a log of all queries in order they're run
	FileEventLogs
	// FileRefstore is a file for the user's local namespace
	FileRefstore
	// FilePeers holds peer repositories
	// Ideally this won't stick around for long
	FilePeers
	// FileAnalytics holds analytics data
	FileAnalytics
	// FileSearchIndex is the path to a search index
	FileSearchIndex
	// FileSelectedRefs is the path to the current ref selection
	FileSelectedRefs
	// FileChangeRequests is a file of change requests
	FileChangeRequests
)

var paths = map[File]string{
	FileUnknown:        "",
	FileLockfile:       "/repo.lock",
	FileInfo:           "/info.json",
	FileConfig:         "/config.json",
	FileDatasets:       "/datasets.json",
	FileEventLogs:      "/events.json",
	FileRefstore:       "/ds_refs.json",
	FilePeers:          "/peers.json",
	FileAnalytics:      "/analytics.json",
	FileSearchIndex:    "/index.bleve",
	FileSelectedRefs:   "/selected_refs.json",
	FileChangeRequests: "/change_requests.json",
}

// Filepath gives the relative filepath to a repofile
// in a given repository
func Filepath(rf File) string {
	return paths[rf]
}

// FileKey returns a datastore.Key reference for a
// given File
func FileKey(rf File) datastore.Key {
	return datastore.NewKey(Filepath(rf))
}
