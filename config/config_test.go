package config_test

import (
	"fmt"
	"github.com/agorago/wego/config"
	"os"
)

func ExampleDevtest() {
	os.Setenv(config.ENVVAR, "dev")
	config.InitConfig("test-configs/env")

	fmt.Println(config.Value("config_test.property1"))
	fmt.Println(config.Value("config_test.property2"))
	fmt.Println(config.Value("config_test.property3"))
	// Output:
	// default-value1
	// dev-value2
	// dev-value3
}

func ExampleProdtest() {
	os.Setenv(config.ENVVAR, "prod")
	config.InitConfig("test-configs/env")
	fmt.Println(config.Value("config_test.property1"))
	fmt.Println(config.Value("config_test.property2"))
	fmt.Println(config.Value("config_test.property3"))
	// Output:
	// default-value1
	// prod-value2
	// prod-value3
}

// ExampleEnvOverride - tests an environment override to the same property
func ExampleEnvOverride() {
	os.Setenv(config.ENVVAR, "prod")
	os.Setenv("CONFIG_TEST.PROPERTY2", "env-value2")
	config.InitConfig("test-configs/env")
	fmt.Println(config.Value("config_test.property1"))
	fmt.Println(config.Value("config_test.property2"))
	fmt.Println(config.Value("config_test.property3"))
	os.Unsetenv("CONFIG_TEST.PROPERTY2")
	// Output:
	// default-value1
	// env-value2
	// prod-value3
}

// ExampleEnv__Override - tests an environment override with __ instead of .
func ExampleEnvOverrideWith__(){
	os.Setenv(config.ENVVAR, "prod")
	os.Setenv("CONFIG_TEST__PROPERTY2", "env--value2")
	config.InitConfig("test-configs/env")
	fmt.Println(config.Value("config_test.property1"))
	fmt.Println(config.Value("config_test.property2"))
	fmt.Println(config.Value("config_test.property3"))
	os.Unsetenv("CONFIG_TEST.PROPERTY2")
	// Output:
	// default-value1
	// env--value2
	// prod-value3
}

/*
func ExampleEtcdOverride(){
	os.Setenv(config.ENVVAR, "prod")
	config.InitConfig("test-configs/env")
	//configureEtcd()
	//log.Printf("I am here\n")
	// setPropertyInEtcd("config_test.property2","etcd-value2")

	fmt.Println(config.Value("config_test.property1"))
	fmt.Println(config.Value("config_test.property2"))
	fmt.Println(config.Value("config_test.property3"))
	fmt.Println(config.Value("config_test.property8"))
	// Output:
	// default-value1
	// prod-value2
	// prod-value3
	// etcd-value8
}
*/
