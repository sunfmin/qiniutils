package qiniutils_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/qiniu/api.v7/storage"
	"github.com/sunfmin/qiniutils"
)

func getQn() (qn *qiniutils.Qiniu) {
	qn = qiniutils.New()
	storageCfg := &storage.Config{}
	storageCfg.Zone = &storage.ZoneHuadong
	qn = qn.Mac(os.Getenv("QINIU_AccessID"), os.Getenv("QINIU_AccessKey")).
		StorageConfig(storageCfg)
	return
}

func TestForEach(t *testing.T) {
	var err error

	qn := getQn()
	var buckets []string
	buckets, err = qn.BucketManager().Buckets(true)
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
				fmt.Println(i, en.Key, qn.NewURLMaker().Domain("sunfmin.com").TTL(time.Hour).Key(en.Key).PrivateURL())
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
	}
}

func TestUpload(t *testing.T) {
	var err error
	key := "hello1.txt"
	qn := getQn().Bucket("sunfminpublic")
	err = qn.BucketManager().Delete("sunfminpublic", key)
	if err != nil {
		fmt.Println(err)
	}
	text := ""
	for i := 0; i < 10000; i++ {
		text += "Hello text"
	}
	err = qn.PutExtra(&storage.PutExtra{
		OnProgress: func(fsize, uploaded int64) {
			fmt.Println("fsize = ", fsize, "uploaded = ", uploaded)
		},
	}).Upload(key, strings.NewReader(text))
	if err != nil {
		panic(err)
	}
}
