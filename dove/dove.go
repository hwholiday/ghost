package dove

import (
	"github.com/golang/protobuf/proto"
	api "github.com/hwholiday/ghost/dove/api/dove"
	"github.com/hwholiday/ghost/dove/network"
	"github.com/rs/zerolog/log"
	"sync"
)

type HandleFunc func(cli network.Conn, data *api.Dove)
type EventFunc func(protocol api.EventType, cli network.Conn)

type Dove interface {
	RegisterHandleFunc(id int32, fn HandleFunc)
	CanAccept(identity string) error
	Accept(opt ...network.Option) error
	Manage() *manage
	EventNotify(EventFunc)
}

type dove struct {
	rw            sync.RWMutex
	manage        *manage
	HandleFuncMap map[int32]HandleFunc
	eventFunc     EventFunc
}

func (h *dove) EventNotify(eventFunc EventFunc) {
	h.eventFunc = eventFunc
}

func NewDove() Dove {
	h := dove{
		manage:        newManage(),
		HandleFuncMap: make(map[int32]HandleFunc),
	}
	setup()
	return &h
}

func (h *dove) CanAccept(identity string) error {
	return h.manage.canAdd(identity)
}

func (h *dove) Manage() *manage {
	return h.manage
}

func (h *dove) RegisterHandleFunc(id int32, fn HandleFunc) {
	h.rw.Lock()
	defer h.rw.Unlock()
	if _, ok := h.HandleFuncMap[id]; ok {
		log.Warn().Int32("id", id).Msg("RegisterHandleFunc already register id")
		return
	}
	h.HandleFuncMap[id] = fn
}

func (h *dove) triggerHandle(client network.Conn, id int32, data *api.Dove) {
	fn, ok := h.HandleFuncMap[id]
	if !ok {
		log.Warn().Int32("id", id).Msg("[Dove] accept HandleFuncMap not register id")
		return
	}
	fn(client, data)
}

func (h *dove) Accept(opt ...network.Option) error {
	opts, err := network.NewOptions(opt...)
	if err != nil {
		return err
	}
	if err = h.manage.canAdd(opts.GetIdentity()); err != nil {
		return err
	}
	var client network.Conn
	if opts.HasConn() {
		client = network.NewConnWithOpt(opts)
	}
	if opts.HasWsConn() {
		client = network.NewWsConnWithOpt(opts)
	}
	h.manage.Add(client)
	if h.eventFunc != nil {
		h.eventFunc(api.EventType_ConnAccept, client)
	}
	go func() {
		for {
			byt, err := client.Read()
			if err != nil {
				h.manage.Del(client.Identity(), client.ConnID())
				if h.eventFunc != nil {
					h.eventFunc(api.EventType_ConnClose, client)
				}
				return
			}
			req, err := parseByt(byt)
			if err != nil {
				log.Error().Err(err).Msg("[Dove] accept parseByt failed")
				continue
			}
			h.triggerHandle(client, req.GetMetadata().GetCrcId(), req)
		}
	}()
	return nil
}

func parseByt(byt []byte) (*api.Dove, error) {
	var req api.Dove
	if err := proto.Unmarshal(byt, &req); err != nil {
		return nil, err
	}
	return &req, nil
}
