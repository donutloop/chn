package mongo

import (
	"fmt"
	"github.com/donutloop/chn/storiesservice/internal/config"
	"github.com/donutloop/chn/storiesservice/internal/storage/object"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/satori/go.uuid"
	"sync"
	"time"
)

type handler struct {
	client   *mgo.Session
	dailInfo *mgo.DialInfo
	database string
}

// NewHandler create new mongo handler
func NewHandler(config *config.Config) *handler {
	h := &handler{
		dailInfo: &mgo.DialInfo{
			PoolLimit: 4096,
			Timeout:   time.Second,
			FailFast:  true,
			Username:  config.Storage.Username,
			Password:  config.Storage.Password,
			Addrs:     []string{config.Storage.Address},
			Database:  config.Storage.Database,
		},
		database: config.Storage.Database,
	}
	return h
}

// Prepare
func (h *handler) Prepare() error {
	mgoSession, err := mgo.DialWithInfo(h.dailInfo)
	if err != nil {
		return err
	}

	// Switch the session to a monotonic behavior.
	mgoSession.SetMode(mgo.Monotonic, true)
	defer mgoSession.Close()

	if err := mgoSession.Ping(); err != nil {
		return err
	}

	s := mgoSession.Clone()
	h.client = s

	return nil
}

// Close
func (h *handler) Close() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Close err %s", r)
		}
	}()
	h.client.Close()
	return
}

func (h *handler) ListParent(parent string, obj object.Interfaces) (err error) {
	s := h.client.Clone()
	defer s.Close()
	return s.DB(h.database).C(obj.GetNamespace()).Find(bson.M{"parent": parent}).All(obj)
}

func (h *handler) List(obj object.Interfaces, opt object.ListOpt) (n int, err error) {

	wg := sync.WaitGroup{}
	wg.Add(2)

	// list
	go func() {
		defer wg.Done()
		s := h.client.Clone()
		defer s.Close()
		var msort string
		switch opt.Sort {
		case object.SortNatural:
			msort = "$natural"
		case object.SortCreatedDesc:
			msort = "-created"
		case object.SortCreatedAsc:
			msort = "+created"
		case object.SortUpdatedDesc:
			msort = "-updated"
		case object.SortUpdatedAsc:
			msort = "+updated"
		}
		err = s.DB(h.database).C(obj.GetNamespace()).Find(nil).Skip(int(opt.Page * opt.Limit)).Limit(int(opt.Limit)).Sort(msort).All(obj)
	}()

	// count
	go func() {
		defer wg.Done()
		s := h.client.Clone()
		defer s.Close()
		n, err = s.DB(h.database).C(obj.GetNamespace()).Count()
	}()

	wg.Wait()
	return
}

func (h *handler) One(obj object.Interface) error {

	s := h.client.New()

	defer s.Close()

	if err := s.DB(h.database).C(obj.GetNamespace()).Find(bson.D{bson.DocElem{Name: "_id", Value: obj.GetId()}}).One(obj); err != nil {
		return fmt.Errorf("`%s::%s::%s`", obj.GetNamespace(), obj.GetId(), err.Error())
	}

	return nil

}

func (h *handler) FindBy(name string, value interface{}, obj object.Interface) error {

	s := h.client.New()

	defer s.Close()

	if err := s.DB(h.database).C(obj.GetNamespace()).Find(bson.M{name: value}).One(obj); err == mgo.ErrNotFound {
		return mgo.ErrNotFound
	} else if err != nil {
		return fmt.Errorf("`%s::%s::%s`", obj.GetNamespace(), obj.GetId(), err.Error())
	}

	return nil
}

// Insert
func (h *handler) Insert(obj object.Interface) error {

	s := h.client.New()

	defer s.Close()

	if v, ok := obj.(object.IdSetter); ok {
		v.SetId(uuid.NewV4().String())
	}

	if v, ok := obj.(object.TimeTracker); ok {
		v.SetCreated(time.Now().Unix())
		v.SetUpdated(time.Now().Unix())
	}

	if err := s.DB(h.database).C(obj.GetNamespace()).Insert(obj); err != nil {
		return err
	}

	return nil

}

func (h *handler) Update(obj object.Interface) error {

	s := h.client.Clone()

	defer s.Close()

	if v, ok := obj.(object.TimeTracker); ok {
		v.SetUpdated(time.Now().Unix())
	}

	if err := s.DB(h.database).C(obj.GetNamespace()).UpdateId(obj.GetId(), obj); err != nil {
		return err
	}

	return nil

}

func (h *handler) Remove(obj object.Interface) error {

	s := h.client.Clone()

	defer s.Close()

	if err := s.DB(h.database).C(obj.GetNamespace()).Remove(bson.D{bson.DocElem{Name: "_id", Value: obj.GetId()}}); err != nil {
		return err
	}

	return nil

}
