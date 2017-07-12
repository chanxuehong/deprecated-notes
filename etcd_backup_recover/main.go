package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	etcd "github.com/coreos/etcd/clientv3"
)

const (
	etcdConfigPrefix = "root/config/"
	etcdDumpFileName = "etcd_config.txt"
)

var etcdEndpoints []string

func main() {
	recv := false
	etcdEndpointsStr := ""
	flag.BoolVar(&recv, "recover", false, "覆盖 etcd 配置")
	flag.StringVar(&etcdEndpointsStr, "etcd_endpoints", "", "etcd endpoints 地址, 逗号分隔, 127.0.0.1:2379,127.0.0.1:12379")
	flag.Parse()

	if etcdEndpointsStr == "" {
		flag.Usage()
		return
	}
	etcdEndpoints = strings.Split(etcdEndpointsStr, ",")

	if recv {
		fmt.Println(save())
	} else {
		fmt.Println(fetch())
	}
}

func fetch() (err error) {
	clt, err := newEtcdClient(etcdEndpoints)
	if err != nil {
		return err
	}
	defer clt.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	resp, err := clt.KV.Get(ctx, etcdConfigPrefix, etcd.WithPrefix())
	if err != nil {
		return err
	}

	kvs := make(map[string]string, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		kvs[string(kv.Key)] = string(kv.Value)
	}

	data, err := json.Marshal(kvs)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(etcdDumpFileName, data, 0666)
}

func save() (err error) {
	data, err := ioutil.ReadFile(etcdDumpFileName)
	if err != nil {
		return err
	}
	kvs := make(map[string]string, 1024)
	if err = json.Unmarshal(data, &kvs); err != nil {
		return err
	}

	clt, err := newEtcdClient(etcdEndpoints)
	if err != nil {
		return err
	}
	defer clt.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	for k, v := range kvs {
		if _, err = clt.KV.Put(ctx, k, v); err != nil {
			return err
		}
	}
	return nil
}

func newEtcdClient(endpoints []string) (*etcd.Client, error) {
	return etcd.New(etcd.Config{
		Endpoints:        endpoints,
		AutoSyncInterval: time.Hour,
		DialTimeout:      time.Second * 5,
	})
}
