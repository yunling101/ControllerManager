package monitor

import (
	"encoding/json"
	"fmt"
	"github.com/yunling101/ControllerManager/logf"
	"github.com/yunling101/ControllerManager/pkg/plugin/request"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
)

const (
	prometheusListen = "http://127.0.0.1:9090"
	groupNameSuffix  = "StatsAlert"
	defaultSuffix    = "rules"
	defaultPrefix    = "yone"
)

type Client struct {
	listen             string
	username, password string
	prefix, suffix     string
	storeDir           string
	group              Group
	fileList           []string
	reload             bool
	err                error
}

func NewPrometheus(store string) *Client {
	return &Client{
		storeDir: store,
		suffix:   defaultSuffix,
		prefix:   defaultPrefix,
		listen:   prometheusListen,
		reload:   true,
	}
}

func (c *Client) SetBasicAuth(username, password string) *Client {
	c.username = username
	c.password = password
	return c
}

func (c *Client) SetPrefix(prefix string) *Client {
	c.prefix = prefix
	return c
}

func (c *Client) SetSuffix(suffix string) *Client {
	c.suffix = suffix
	return c
}

func (c *Client) Listen(address string) *Client {
	c.listen = address
	return c
}

func (c *Client) Reload(reload bool) *Client {
	c.reload = reload
	return c
}

func (c *Client) getGroupFileName(name string) string {
	return c.prefix + "-" + name + "." + c.suffix
}

func (c *Client) getGroupName(name string, stats bool) string {
	if !stats {
		return strings.TrimSuffix(name, groupNameSuffix)
	}
	return name + groupNameSuffix
}

func (c *Client) groupNameConversion(group Group) Group {
	for idx, n := range group.Groups {
		if !strings.Contains(n.Name, groupNameSuffix) {
			group.Groups[idx].Name = c.getGroupName(n.Name, true)
		}
	}
	return group
}

func (c *Client) groupHandler(rules string, handler func(Group) Group) (err error) {
	if err = json.Unmarshal([]byte(rules), &c.group); err != nil {
		return
	}

	if err = c.ergodicDir(c.storeDir, "."+c.suffix); err != nil {
		return
	}

	c.group = handler(c.group)
	if len(c.group.Groups) != 0 {
		logf.Logger().ErrorIf(c.writeFile(c.groupNameConversion(c.group)))
	}
	if c.reload {
		logf.Logger().ErrorIf(c.reloadConfig())
	}
	return
}

func (c *Client) AddRules(rules string) (err error) {
	err = c.groupHandler(rules, func(group Group) Group {
		for _, rule := range group.Groups {
			newName := c.getGroupFileName(rule.Name)
			if c.groupFileName(newName) {
				var fileContent Group
				fileContent, c.err = c.readFile(fmt.Sprintf("%s/%s", c.storeDir, newName))
				if c.err == nil {
					for idx, n := range fileContent.Groups {
						alert, index, ok := c.getOldRuleExits(n.Rules, rule.Rules)
						fileContent.Groups[idx].Name = c.getGroupName(rule.Name, true)
						if !ok {
							n.Rules = append(n.Rules, rule.Rules...)
							fileContent.Groups[idx].Rules = n.Rules
						} else {
							fileContent.Groups[idx].Rules[index] = alert
						}
					}
					group = fileContent
				}
			}
		}
		return group
	})
	return
}

func (c *Client) DeleteRule(rules string) (err error) {
	err = c.groupHandler(rules, func(group Group) Group {
		for _, rule := range group.Groups {
			newName := c.getGroupFileName(rule.Name)
			newRules := make([]Rule, 0)
			if c.groupFileName(newName) {
				var fileContent Group
				fileContent, c.err = c.readFile(fmt.Sprintf("%s/%s", c.storeDir, newName))
				if c.err == nil {
					for _, n := range fileContent.Groups {
						newRule := Rule{Name: n.Name, Rules: []Alert{}}
						for _, v := range n.Rules {
							if !c.getRuleKey(rule.Rules, v.Alert) {
								newRule.Rules = append(newRule.Rules, v)
							}
						}
						newRules = append(newRules, newRule)
					}
					if len(newRules) == 1 && len(newRules[0].Rules) == 0 {
						filename := fmt.Sprintf("%s/%s", c.storeDir, c.getGroupFileName(rule.Name))
						if err = os.RemoveAll(filename); err != nil {
							logf.Logger().Warn(err)
						} else {
							newRules = make([]Rule, 0)
						}
					}
				}
			}
			group.Groups = newRules
		}
		return group
	})
	return
}

func (c *Client) groupFileName(newName string) bool {
	for _, rule := range c.fileList {
		oldName := filepath.Base(rule)
		if oldName == newName {
			return true
		}
	}
	return false
}

func (c *Client) getOldRuleExits(oldRule, newRule []Alert) (old Alert, index int, ok bool) {
	for index, old = range oldRule {
		for _, n := range newRule {
			if old.Alert == n.Alert {
				old = n
				ok = true
				return
			}
		}
	}
	return
}

func (c *Client) getRuleKey(rules []Alert, key string) bool {
	for _, k := range rules {
		if k.Alert == key {
			return true
		}
	}
	return false
}

func (c *Client) ergodicDir(path, suffix string) error {
	root := string(os.PathSeparator)
	if string(path[len(path)-1]) != root {
		path = path + root
	}

	dirEntry, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range dirEntry {
		if entry.IsDir() {
			_ = c.ergodicDir(path+entry.Name()+root, suffix)
		} else {
			if ok := strings.HasSuffix(entry.Name(), suffix); ok {
				c.fileList = append(c.fileList, path+entry.Name())
			}
		}
	}

	return nil
}

func (c *Client) reloadConfig() error {
	url := fmt.Sprintf("%s/%s", c.listen, "-/reload")
	response := request.New(url).SetMethod("POST").SetBasicAuth(c.username, c.password).Do()
	return response.Error
}

func (c *Client) readFile(filename string) (rule Group, err error) {
	var readBody []byte
	readBody, c.err = os.ReadFile(filename)
	if c.err != nil {
		err = c.err
		return
	}

	err = yaml.Unmarshal(readBody, &rule)
	return
}

func (c *Client) writeFile(rule Group) (err error) {
	var readBody []byte
	readBody, c.err = yaml.Marshal(rule)
	if c.err != nil {
		err = c.err
		return
	}

	for _, n := range rule.Groups {
		filename := fmt.Sprintf("%s/%v", c.storeDir,
			c.getGroupFileName(c.getGroupName(n.Name, false)))
		logf.Logger().WarnIf(os.WriteFile(filename, readBody, 0644))
	}
	return
}
