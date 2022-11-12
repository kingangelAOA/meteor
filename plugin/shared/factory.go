package shared

import (
	"common"
	"errors"
	"fmt"
	"os"
	"os/exec"
	SPlugin "plugin"
	"web/configs"

	"github.com/hashicorp/go-hclog"
	GPlugin "github.com/hashicorp/go-plugin"
)

type PluginFactory struct {
	ID       string
	Code     string
	Language string
	RootPath string
}

func NewPluginFactory(id, code, language, RootPath string) *PluginFactory {
	return &PluginFactory{
		ID:       id,
		Code:     code,
		Language: language,
		RootPath: RootPath,
	}
}

func (pf *PluginFactory) GetExecutionFilePath() (string, error) {
	switch pf.Language {
	case common.Python:
		return buildPythonPlugin(pf.ID, pf.RootPath, []byte(pf.Code))
	case common.Go:
		return buildGoPlugin(pf.ID, pf.RootPath, []byte(pf.Code))
	default:
		return "", errors.New("supported python and go")
	}
}

func (pf *PluginFactory) GetPluginStream() ([]byte, error) {
	switch pf.Language {
	case common.Python:
		return []byte(pf.Code), nil
	case common.Go:
		path, err := buildGoPlugin(pf.ID, pf.RootPath, []byte(pf.Code))
		if err != nil {
			return nil, err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, errors.New("read plugin binary error")
		}
		return content, nil
	default:
		return nil, errors.New("supported python and go")
	}
}

func (pf *PluginFactory) CheckFile() error {
	switch pf.Language {
	case common.Python:
		return nil
	case common.Go:
		_, err := buildGoPlugin(pf.ID, pf.RootPath, []byte(pf.Code))
		return err
	default:
		return errors.New("supported python and go")
	}
}

func (pf *PluginFactory) GetPluginClientFromDB(downloadBinary func(string, string) error) (PluginClient, error) {
	binaryPath, err := GetBinaryPath(pf.ID, pf.Language, *configs.PluginRootPath)
	if err != nil {
		return nil, err
	}
	err = downloadBinary(pf.ID, binaryPath)
	if err != nil {
		return nil, err
	}
	return loadPlugin(binaryPath, pf.Language)
}

func (pf *PluginFactory) GetPluginClient() (PluginClient, error) {
	path, err := pf.GetExecutionFilePath()
	if err != nil {
		return nil, err
	}
	return loadPlugin(path, pf.Language)
}

func GetBinaryPath(key, language, pluginRootPath string) (string, error) {
	switch language {
	case common.Python:
		return getPythonbinaryPath(key, pluginRootPath)
	case common.Go:
		return getGoBinaryPath(key, pluginRootPath)
	default:
		return "", errors.New("supported python and go")
	}
}

func getGoBinaryPath(key, pluginRootPath string) (string, error) {
	goPluginPath, err := getGoPluginPath(key, pluginRootPath)
	if err != nil {
		return "", err
	}
	if common.GoPluginPerformanceSwitch {
		return fmt.Sprintf("%s/%s", goPluginPath, common.GoPluginBinaryNameSO), nil
	} else {
		return fmt.Sprintf("%s/%s", goPluginPath, common.GoPluginBinaryName), nil
	}
}

func getGoMainPath(key, pluginRootPath string) (string, error) {
	goPluginPath, err := getGoPluginPath(key, pluginRootPath)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/main.go", goPluginPath), nil
}

func getGoPluginPath(key, pluginRootPath string) (string, error) {
	pluginPath := fmt.Sprintf("%s/%s", pluginRootPath, "go")
	err := initDir(pluginPath)
	if err != nil {
		err = os.Mkdir(pluginPath, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("mkdir go plugin fold error")
		}
	}
	goPluginPath := fmt.Sprintf("%s/%s", pluginPath, key)
	err = initDir(goPluginPath)
	if err != nil {
		err = os.Mkdir(goPluginPath, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("mkdir go plugin fold error")
		}
	}
	return goPluginPath, nil
}

func buildGoPlugin(key, pluginRootPath string, code []byte) (string, error) {
	mainPath, err := getGoMainPath(key, pluginRootPath)
	if err != nil {
		return "", err
	}
	binaryPath, err := getGoBinaryPath(key, pluginRootPath)
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(mainPath, code, 0755); err != nil {
		return "", fmt.Errorf("write plugin file error, %s", err.Error())
	}
	if err := goCmd(binaryPath, mainPath); err != nil {
		return "", err
	}
	return binaryPath, nil
}

func goCmd(binaryPath, mainPath string) error {
	var cmd *exec.Cmd
	if common.GoPluginPerformanceSwitch {
		cmd = exec.Command("go", "build", "-buildmode=plugin", "-o", binaryPath, mainPath)
	} else {
		cmd = exec.Command("go", "build", "-o", binaryPath, mainPath)
	}
	res, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(res))
	}
	return nil
}

func buildPythonPlugin(key, pluginRootPath string, code []byte) (string, error) {
	pyPluginPath, err := getPythonbinaryPath(key, pluginRootPath)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(pyPluginPath, code, 0755)
	if err != nil {
		return "", fmt.Errorf("write plugin file error, %s", err.Error())
	}
	return pyPluginPath, nil
}

func getPythonbinaryPath(key, pluginRootPath string) (string, error) {
	pluginPath := fmt.Sprintf("%s/%s", pluginRootPath, "python")
	err := initDir(pluginPath)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s.py", pluginPath, key), nil
}

func loadPlugin(path, language string) (PluginClient, error) {
	switch language {
	case common.Python:
		return newGRPCPluginClientWrap(loadGRPCPluginClient(path))
	case common.Go:
		if common.GoPluginPerformanceSwitch {
			p, err := loadPluginClient(path)
			if err != nil {
				return nil, err
			}
			return newPluginClientWrap(p), nil
		} else {
			return newGRPCPluginClientWrap(loadGRPCPluginClient(path))
		}
	default:
		return nil, errors.New("supported python and go")
	}
}

func loadPluginClient(path string) (Plugin, error) {
	p, err := SPlugin.Open(path)
	if err != nil {
		return nil, err
	}
	symbol, err := p.Lookup(common.SysPluginName)
	if err != nil {
		return nil, err
	}
	return symbol.(Plugin), nil
}

func loadGRPCPluginClient(path string) *GPlugin.Client {
	cmd := getCmd(path)
	logger := hclog.New(&hclog.LoggerOptions{
		Output: hclog.DefaultOutput,
		Level:  hclog.Error,
		Name:   "plugin",
	})

	return GPlugin.NewClient(&GPlugin.ClientConfig{
		HandshakeConfig: Handshake,
		Plugins:         PluginMap,
		Logger:          logger,
		Cmd:             cmd,
		AllowedProtocols: []GPlugin.Protocol{
			GPlugin.ProtocolNetRPC, GPlugin.ProtocolGRPC},
	})
}

func initDir(path string) error {
	_, erByStat := os.Stat(path)
	if erByStat != nil {
		return errors.New("need to deploy plugin folder")
	}
	return nil
}

func getCmd(path string) *exec.Cmd {
	return exec.Command("sh", "-c", path)
}

//func GetPluginFilePath(key, pluginRootPath, language string) (string, error) {
//	switch language {
//	case common.Python:
//		pluginPath := fmt.Sprintf("%s/%s", pluginRootPath, "python")
//		err := initDir(pluginPath)
//		if err != nil {
//			err = os.Mkdir(pluginPath, os.ModePerm)
//			if err != nil {
//				return "", fmt.Errorf("mkdir go plugin fold error")
//			}
//		}
//		return fmt.Sprintf("%s/%s.py", pluginPath, key), nil
//	case common.Go:
//		return getGoMainPath(key, pluginRootPath)
//	default:
//		return "", errors.New("supported python and go")
//	}
//}
