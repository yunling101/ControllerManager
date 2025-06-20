package monitor

import (
	"encoding/base64"
	"fmt"
	"github.com/toolkits/file"
	"github.com/yunling101/ControllerManager/logf"
	"gopkg.in/yaml.v2"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"time"
)

const (
	alertManagerListen = "http://127.0.0.1:9093"
	defaultName        = "alertmanager.yml"
)

func NewAlertManager(store string) *Client {
	return &Client{
		storeDir: store,
		listen:   alertManagerListen,
		prefix:   defaultName,
		suffix:   "tmpl",
		reload:   true,
	}
}

func (c *Client) getFilename() string {
	return fmt.Sprintf("%s/%s", c.storeDir, c.prefix)
}

func (c *Client) LoadFile() (cfg AlertManagerConfig, err error) {
	var body []byte
	body, c.err = os.ReadFile(c.getFilename())
	if c.err != nil {
		err = c.err
		return
	}

	err = yaml.Unmarshal(body, &cfg)
	return
}

func (c *Client) ModifyFile(body []byte) (err error) {
	var oldCfg, newCfg AlertManagerConfig
	if err = yaml.Unmarshal(body, &newCfg); err != nil {
		return
	}

	if oldCfg, err = c.LoadFile(); err != nil {
		return
	}

	if reflect.DeepEqual(newCfg, oldCfg) {
		return
	}

	if err = c.backupFile(oldCfg); err != nil {
		return
	}

	var newBody []byte
	if newBody, err = yaml.Marshal(newCfg); err != nil {
		return
	}

	err = os.WriteFile(c.getFilename(), newBody, 0644)
	if err != nil {
		logf.Logger().ErrorIf(err)
		return
	}

	if c.reload {
		logf.Logger().ErrorIf(c.reloadConfig())
	}
	return
}

func (c *Client) backupFile(cfg AlertManagerConfig) (err error) {
	backupDir := c.storeDir + "/offset"
	if err = file.InsureDir(backupDir); err != nil {
		return
	}

	var b []byte
	if b, err = yaml.Marshal(cfg); err != nil {
		return
	}

	t := time.Unix(time.Now().Unix(), 0).Format("20060102_150405")
	filename := fmt.Sprintf("%s/%s%s", backupDir, t, path.Ext(c.prefix))
	err = os.WriteFile(filename, b, 0644)
	return
}

func (c *Client) Templates(filename, content string) (err error) {
	var decode []byte
	decode, err = base64.StdEncoding.DecodeString(content)
	if err != nil {
		return
	}

	fileDir := fmt.Sprintf("%s/template", c.storeDir)
	if err = file.InsureDir(fileDir); err != nil {
		return
	}

	name := filepath.Base(filename)
	err = os.WriteFile(fmt.Sprintf("%s/%s", fileDir, name), decode, 0644)
	if err != nil {
		logf.Logger().ErrorIf(err)
		return
	}

	if c.reload {
		logf.Logger().ErrorIf(c.reloadConfig())
	}
	return
}
