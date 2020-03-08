package config

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	DEFAULT_ENV = "default"
	PROPERTY_KEY_IN_ETCD = "keys"
)
var viperConfig = viper.New()
func init() {
	viperConfig.SetConfigType("toml")
	viperConfig.AutomaticEnv()
	configPath := GetConfigPath() + "/env"
	InitConfig(configPath)
	// ConfigureEtcd()
}

// InitConfig - allow end users to set alternate config paths.
// By default we use config path defined in an environment variable
// This enables testing
func InitConfig(cpath string) {
	walkTree(cpath + "/" + DEFAULT_ENV)
	walkTree(cpath + "/" + GetEnv() )
}

func ConfigureEtcd(){
	log.Printf("Adding remote provider %s\n",GetEtcdEndPoint())
	viperConfig.AddRemoteProvider("etcd", GetEtcdEndPoint(), PROPERTY_KEY_IN_ETCD)

	// read from remote config the first time.
	err := viperConfig.ReadRemoteConfig()
	if err != nil {
		// log.Errorf("unable to read remote config: %v", err)
		fmt.Fprintf(os.Stderr,"Cannot read remote config: %v\n",err)
	}
	log.Printf("Successfully read the messages from etcd\n")

	// open a goroutine to watch remote changes forever
	go func(){
		for {
			time.Sleep(time.Second * time.Duration(GetEtcdPollingDelay())) // delay after each request

			err := viperConfig.WatchRemoteConfig()
			if err != nil {
				// log.Errorf("unable to read remote config: %v", err)
				fmt.Fprintf(os.Stderr,"Cannot read remote config: %v\n",err)
				return
			}
		}
	}()
}


func walkTree(cpath string){
	filepath.Walk(cpath, func(s string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		fmt.Fprintf(os.Stderr,"Reading file %s\n",s)
		buf, err1 := ioutil.ReadFile(s)
		if err1 != nil {
			return err
		}
		viperConfig.MergeConfig(bytes.NewBuffer(buf))
		return nil
	})
}

func Value(propname string)string{
	return viperConfig.GetString(propname)
}

func IntValue(propname string)int{
	return viperConfig.GetInt(propname)
}

func BoolValue(propname string)bool{
	return viperConfig.GetBool(propname)
}
