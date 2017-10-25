package qiniutils_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"github.com/sunfmin/qiniutils"
)

func TestForEach(t *testing.T) {
	var err error
	mac := qbox.NewMac(os.Getenv("QINIU_AccessID"), os.Getenv("QINIU_AccessKey"))
	storageCfg := &storage.Config{}
	storageCfg.Zone = &storage.ZoneHuadong
	bm := storage.NewBucketManager(mac, storageCfg)

	var buckets []string
	buckets, err = bm.Buckets(true)
	if err != nil {
		panic(err)
	}
	fmt.Println(buckets)
	i := 0
	for _, bucket := range buckets {
		fmt.Printf("bucket: %s\n", bucket)
		err = qiniutils.Bucket(bm, bucket).ForEach(1000, func(entries []storage.ListItem, commonPrefixes []string) error {
			for _, en := range entries {
				i++
				fmt.Println(i, en.Key)
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
	}
}
