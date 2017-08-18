package helpers

import (
	"application/libraries/toml"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/coreos/etcd/clientv3"
)

func LoadFile(filename string) (filedata string, err error) {
	fileDataByte, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Read failed", err)
		return
	}

	filedata = string(fileDataByte)
	return
}

func EtcdDail(etcdconfig toml.DBConfig) (c *clientv3.Client, err error) {
	cfg := clientv3.Config{
		Endpoints: []string{etcdconfig.Host},
		//Transport:   client.DefaultTransport,
		DialTimeout: 5 * time.Second,
	}
	c, err = clientv3.New(cfg)

	return
}

func Sha1Str(s string) string {
	r := sha1.Sum([]byte(s))
	return hex.EncodeToString(r[:])
}
