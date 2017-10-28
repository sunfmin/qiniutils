package qiniutils_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/qiniu/api.v7/storage"
	"github.com/sunfmin/qiniutils"
)

func TestForEach(t *testing.T) {
	var err error

	qn := qiniutils.New()
	storageCfg := &storage.Config{}
	storageCfg.Zone = &storage.ZoneHuadong
	qn = qn.Mac(os.Getenv("QINIU_AccessID"), os.Getenv("QINIU_AccessKey")).
		StorageConfig(storageCfg)

	var buckets []string
	buckets, err = qn.GetBuckets(true)
	if err != nil {
		panic(err)
	}
	fmt.Println(buckets)
	i := 0
	for _, bucket := range buckets {
		fmt.Printf("bucket: %s\n", bucket)
		err = qn.Bucket(bucket).ForEach(1000, func(entries []storage.ListItem, commonPrefixes []string) error {
			for _, en := range entries {
				i++
				fmt.Println(i, en.Key, qn.URL().Domain("sunfmin.com").TTL(time.Hour).Key(en.Key).PrivateURL())
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
	}
}
