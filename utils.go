package qiniutils

import (
	"github.com/qiniu/api.v7/storage"
)

type Qn struct {
	bm     *storage.BucketManager
	bucket string
	prefix string
}

func clone(u *Qn) (r *Qn) {
	r = &Qn{}
	r.bm = u.bm
	r.bucket = u.bucket
	r.prefix = u.prefix
	return
}

func New() *Qn {
	return &Qn{}
}

func (q *Qn) Bucket(bm *storage.BucketManager, bucket string) (r *Qn) {
	r = clone(q)
	r = &Qn{bm, bucket, ""}
	return
}

func (u *Qn) Prefix(prefix string) (r *Qn) {
	u.prefix = prefix
	r = u
	return
}

func (u *Qn) ForEach(limit int, f func(entries []storage.ListItem, commonPrefixes []string) error) (err error) {
	var (
		entries        []storage.ListItem
		commonPrefixes []string
		nextMarker     string
		hasNext        bool = true
	)
	for hasNext {
		entries, commonPrefixes, nextMarker, hasNext, err = u.bm.ListFiles(u.bucket, u.prefix, "", nextMarker, 1000)
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
