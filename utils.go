package qiniutils

import (
	"github.com/qiniu/api.v7/storage"
)

type udata struct {
	bm     *storage.BucketManager
	bucket string
	prefix string
}

func Bucket(bm *storage.BucketManager, bucket string) (r *udata) {
	r = &udata{bm, bucket, ""}
	return
}

func (u *udata) Prefix(prefix string) (r *udata) {
	u.prefix = prefix
	r = u
	return
}

func (u *udata) ForEach(limit int, f func(entries []storage.ListItem, commonPrefixes []string) error) (err error) {
	var (
		entries        []storage.ListItem
		commonPrefixes []string
		nextMarker     string
		hasNext        bool = true
	)
	for hasNext {
		entries, commonPrefixes, nextMarker, hasNext, err = u.bm.ListFiles(u.bucket, "", "", nextMarker, 1000)
		if err != nil {
			return
		}
		err = f(entries, commonPrefixes)
		if err != nil {
			return
		}
	}
	return
}
