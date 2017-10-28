package qiniutils

import (
	"time"

	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

type Qiniu struct {
	bm            *storage.BucketManager
	bucket        string
	prefix        string
	storageConfig *storage.Config
	mac           *qbox.Mac
}

type URLMaker struct {
	qn     *Qiniu
	domain string
	key    string
	ttl    time.Duration
}

func (u *Qiniu) clone() (r *Qiniu) {
	r = &Qiniu{}
	r.bm = u.bm
	r.bucket = u.bucket
	r.prefix = u.prefix
	r.mac = u.mac
	r.storageConfig = u.storageConfig
	return
}

func New() *Qiniu {
	return &Qiniu{}
}

func (u *URLMaker) clone() (r *URLMaker) {
	r = &URLMaker{}
	r.qn = u.qn
	r.domain = u.domain
	r.key = u.key
	r.ttl = u.ttl
	return
}

func (q *Qiniu) NewURLMaker() (r *URLMaker) {
	r = &URLMaker{}
	r.qn = q
	return
}

func (u *URLMaker) Domain(domain string) (r *URLMaker) {
	r = u.clone()
	r.domain = domain
	return
}

func (u *URLMaker) TTL(ttl time.Duration) (r *URLMaker) {
	r = u.clone()
	r.ttl = ttl
	return
}

func (u *URLMaker) Key(key string) (r *URLMaker) {
	r = u.clone()
	r.key = key
	return
}

func (u *URLMaker) PublicURL() (r string) {
	storage.MakePublicURL(u.domain, u.key)
	return
}

func (u *URLMaker) PrivateURL() (r string) {
	deadline := time.Now().Add(u.ttl).Unix()
	return storage.MakePrivateURL(u.qn.mac, u.domain, u.key, deadline)
}

func (q *Qiniu) Mac(accessKey, secretKey string) (r *Qiniu) {
	r = q.clone()
	r.mac = qbox.NewMac(accessKey, secretKey)
	return
}

func (q *Qiniu) StorageConfig(cfg *storage.Config) (r *Qiniu) {
	r = q.clone()
	r.storageConfig = cfg
	return
}

func (q *Qiniu) Bucket(bucket string) (r *Qiniu) {
	r = q.clone()
	r.bucket = bucket
	return
}

func (q *Qiniu) Prefix(prefix string) (r *Qiniu) {
	r = q.clone()
	r.prefix = prefix
	return
}

func (u *Qiniu) BucketManager() *storage.BucketManager {
	if u.bm != nil {
		return u.bm
	}

	if u.mac == nil {
		panic("call .Mac(xx) first to set mac")
	}

	if u.storageConfig == nil {
		panic("call .StorageConfig(xx) to set storage config")
	}

	u.bm = storage.NewBucketManager(u.mac, u.storageConfig)

	return u.bm
}

func (u *Qiniu) ForEach(limit int, f func(entries []storage.ListItem, commonPrefixes []string) error) (err error) {
	var (
		entries        []storage.ListItem
		commonPrefixes []string
		nextMarker     string
		hasNext        bool = true
	)
	for hasNext {
		entries, commonPrefixes, nextMarker, hasNext, err = u.BucketManager().ListFiles(u.bucket, u.prefix, "", nextMarker, 1000)
		if err != nil {
			return
		}
		if len(entries) == 0 {
			return
		}
		err = f(entries, commonPrefixes)
		if err != nil {
			return
		}
	}
	return
}
