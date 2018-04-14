package storage

import (
	"github.com/donutloop/chn/internal/storage/object"
	"github.com/donutloop/chn/internal/storage/handler/mongo"
	"github.com/pkg/errors"
	"github.com/donutloop/chn/internal/api"
)

type handlerName string

const (
	mongodbHandler handlerName = "mongodb"
)

type (
	// Interface defines the base functionality which any storage handler
	// should implement to become valid storage handler
	Interface interface {
		Prepare() error
		Close() error
		One(doc object.Interface) error
		List(docs object.Interfaces, opt object.ListOpt) (int, error)
		ListParent(parent string, docs object.Interfaces) error
		Insert(doc object.Interface) error
		Update(doc object.Interface) error
		Remove(doc object.Interface) error
		FindBy(name string, value interface{}, obj object.Interface) error
	}
)

// New creates storage handler from config.Storage and prepare it for use
// returns error if something went wrong during the preparations
func New(config *api.Config) (Interface, error) {
	// create handler based on the storage config
	var h Interface
	switch handlerName(config.Storage.Handler) {
	case mongodbHandler:
		h = mongo.NewHandler(config)
	default:
		return nil, errors.New("invalid storage handler `" + config.Storage.Handler + "`")
	}

	if err := h.Prepare(); err != nil {
		return nil, err
	}

	// prepare handler
	return h, nil
}
