package nonce

import (
	"github.com/bradfitz/gomemcache/memcache"

	"github.com/kilfu0701/echo_app/core"
)

const Prefix = "nonce."

type NonceData struct {
	Nonce string
}

type Nonce struct {
	//Ctx        context.Context
	CacheClient *memcache.Client
}

func New(ac core.AppContext) *Nonce {
	return &Nonce{
		CacheClient: ac.AppMemcached,
	}
}

func (n *Nonce) FindByKey(key string) (*NonceData, error) {
	item, _ := n.CacheClient.Get(key)

	var nd NonceData

	if err := core.GobUnmarshal(item.Value, &nd); err != nil {
		return nil, err
	}

	return &nd, nil
}

func (n *Nonce) SetByKey(key string, i interface{}) error {
	v, err := core.GobMarshal(i)
	if err != nil {
		return err
	}

	return n.CacheClient.Set(&memcache.Item{Key: key, Value: v})
}
