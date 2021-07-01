package util

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
)

const (
	// vsc 文件夹，~/.vsc
	vscDirectory = "/.vsc"

	// vsc 配置文件, ~/.vsc/vsc-config.json
	VscConfigFilePath = "/vsc-config.json"
)

// config 文件设置
type VscConfigJSON struct {
	Golangci string `json:"golangci,omitempty"`
	Eslint   struct {
		TS string `json:"typescript,omitempty"`
		JS string `json:"javascript,omitempty"`
	} `json:"eslint,omitempty"`
}

// 全局 vsc 配置文件地址 ~/.vsc
func GetVscConfigDir() (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		return "", errors.New("$HOME is not exist, please set $HOME env")
	}

	return home + vscDirectory, nil
}

func (vs *VscConfigJSON) ReadFromDir(vscDir string) error {
	// read vsc config file
	f, err := os.Open(vscDir + VscConfigFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// ~/.vsc/vsc-config 文件存在, 读取文件
	err = vs.readJSON(f)
	if err != nil {
		return err
	}

	return nil
}

func (vs *VscConfigJSON) readJSON(reader io.Reader) error {
	de := json.NewDecoder(reader)
	return de.Decode(vs)
}

func (vs *VscConfigJSON) JSONIndentFormat() ([]byte, error) {
	return json.MarshalIndent(vs, "", "  ")
}

// 读取 .vscode/settings.json 文件, 获取想要的值
func ReadSettingJSON(v interface{}) error {
	// 读取 .vscode/settings.json
	settingsPath, err := filepath.Abs(SettingsJSONPath)
	if err != nil {
		return err
	}

	sf, err := os.Open(settingsPath)
	if err != nil {
		return err
	}
	defer sf.Close()

	// json 反序列化 settings.json
	jsonc, err := io.ReadAll(sf)
	if err != nil {
		return err
	}

	js, err := JSONCToJSON(jsonc)
	if err != nil {
		return err
	}

	err = json.Unmarshal(js, v)
	if err != nil {
		return err
	}

	return nil
}
