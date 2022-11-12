package shared_test

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"plugin/shared"
	"runtime"
	"testing"
	"time"
	"web/configs"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"gopkg.in/yaml.v2"
)

func TestPlugin(t *testing.T) {
	logger := hclog.New(&hclog.LoggerOptions{
		Output: hclog.DefaultOutput,
		Level:  hclog.Error,
		Name:   "plugin",
	})
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		Logger:          logger,
		Cmd:             exec.Command("sh", "-c", "/Users/eleme/meteor/plugin/go/6304cf3d95182056b4794ad3/main"),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	})
	startT := time.Now()
	client.Start()
	defer client.Kill()
	rpcClient, err := client.Client()
	tc := time.Since(startT)
	fmt.Printf("InitPlugin time cost = %v\n", tc)

	if err != nil {
		fmt.Println(err.Error())
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(shared.DefaultGRPC)
	if err != nil {
		fmt.Println(err.Error())
	}
	plugin := raw.(shared.Plugin)
	r, err := plugin.Run(map[string]string{"a": "b"})
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(r)
}

func BenchmarkPlugin(b *testing.B) {
	runtime.GOMAXPROCS(1)
	b.SetParallelism(1)

	logger := hclog.New(&hclog.LoggerOptions{
		Output: hclog.DefaultOutput,
		Level:  hclog.Info,
		Name:   "plugin",
	})
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  shared.Handshake,
		Plugins:          shared.PluginMap,
		Logger:           logger,
		Cmd:              exec.Command("sh", "-c", "/Users/eleme/meteor/plugin/go/6304cf3d95182056b4794ad3/main"),
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})
	startT := time.Now()
	client.Start()
	defer client.Kill()
	rpcClient, err := client.Client()
	tc := time.Since(startT)
	fmt.Printf("InitPlugin time cost = %v\n", tc)
	if err != nil {
		fmt.Println(err.Error())
	}
	raw, err := rpcClient.Dispense(shared.DefaultGRPC)
	if err != nil {
		panic(err)
	}
	index := 0
	b.RunParallel(func(pb *testing.PB) {
		plugin := raw.(shared.Plugin)
		for pb.Next() {
			_, err := plugin.Run(map[string]string{"a": "b"})
			if err != nil {
				panic(err)
			}
			// fmt.Println(r)
		}
		index += 1
	})
	fmt.Println("****")
}

func TestPluginMan(t *testing.T) {
	path := "/Users/eleme/meteor/plugin/python/63120e906474dd15d4293631.py"
	cmd := exec.Command("sh", "-c", path)
	pc := shared.LoadPluginClient(cmd)
	startT := time.Now()
	rpcClient, err := pc.Client()
	if err != nil {
		fmt.Println(err.Error() + "xxxxxx")
	}
	tc := time.Since(startT)
	fmt.Println(rpcClient.Ping())
	fmt.Printf("InitPlugin time cost = %v\n", tc)
	// Request the plugin
	raw, err := rpcClient.Dispense(shared.DefaultGRPC)
	if err != nil {
		fmt.Println(err.Error())
	}
	plugin := raw.(shared.Plugin)
	r, err := plugin.Run(map[string]string{"a": "b"})
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(r)
}

func getConfig() *configs.Config {
	yamlFile, err := ioutil.ReadFile(*configs.ConfigPath)
	if err != nil {
		panic(err)
	}
	var c configs.Config
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		panic(err)
	}
	return &c
}

func TestName(t *testing.T) {

}
