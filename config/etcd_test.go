package config_test

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/embed"
	"gitlab.intelligentb.com/devops/bplus/config"
	"log"
	"os"
	"time"
)

func configureEtcd(){
	os.Setenv(config.ETCD_ENDPOINTVAR,"localhost:2379")
	os.Setenv(config.ETCD_POLLINGDELAYVAR,"30")
	startEmbeddedEtcd()
	config.ConfigureEtcd()
}

func startEmbeddedEtcd(){
	cfg := embed.NewConfig()
	cfg.Dir = "default.etcd"
	e, err := embed.StartEtcd(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer e.Close()
	select {
	case <-e.Server.ReadyNotify():
		log.Printf("Server is ready!")
		setPropertyInEtcd(config.PROPERTY_KEY_IN_ETCD,`
[config_test]
property8 = "etcd-value8"
`)
		return
	case <-time.After(60 * time.Second):
		e.Server.Stop() // trigger a shutdown
		log.Printf("Server took too long to start!")
	}

	log.Fatal(<-e.Err())
}

func setPropertyInEtcd(name string,value string){
	log.Printf("I am trying to set the property of %s\n",name)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		log.Printf("Cannot set property %s. Err = %s\n", name,err)
	}

	defer cli.Close()
	x,err := cli.Put(context.TODO(),name,value)
	if err != nil {
		log.Printf("Cannot set the key to the value. error = %s\n",err.Error())
	}

	log.Printf("property %s set. response.Header = %#v. prevkv = %#v. Unrecognized = %#v\n",name,x.Header,x.PrevKv,x.XXX_unrecognized)
}