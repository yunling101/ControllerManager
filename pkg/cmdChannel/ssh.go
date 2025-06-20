package cmdChannel

import (
	"fmt"
	"github.com/toolkits/file"
	"github.com/yunling101/ControllerManager/common"
	"golang.org/x/crypto/ssh"
	"os"
	"time"
)

type SSH struct {
	Host            string
	Port            int
	User            string
	Type            string
	Password        string
	KeyBody         string
	KeyFile         string
	KeyFilePassword string
}

func (s *SSH) NewClient() (client *ssh.Client, err error) {
	config := &ssh.ClientConfig{
		Timeout:         time.Second * 60,
		User:            s.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	var signer ssh.AuthMethod
	if s.Type == common.KEY {
		if s.KeyBody != "" {
			signer, err = s.keyBodyAuthFunc()
		} else {
			signer, err = s.keyPathAuthFunc()
		}
	} else {
		signer = ssh.Password(s.Password)
	}

	config.Auth = []ssh.AuthMethod{signer}
	client, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port), config)

	return
}

func (s *SSH) RunCommand(client *ssh.Client, command string) (r common.Response) {
	session, err := client.NewSession()
	if err != nil {
		r.Error = fmt.Errorf("NewSessionError: %s", err.Error())
		return
	}
	defer func() {
		_ = session.Close()
		_ = client.Close()
	}()

	r.Data, r.Error = session.CombinedOutput(command)
	return
}

func (s *SSH) keyPathAuthFunc() (signer ssh.AuthMethod, err error) {
	if !file.IsExist(s.KeyFile) {
		err = fmt.Errorf("%s key file does not exist", s.KeyFile)
		return
	}

	var (
		keyBody []byte
		parse   ssh.Signer
	)
	keyBody, err = os.ReadFile(s.KeyFile)
	if err != nil {
		return
	}

	if s.KeyFilePassword != "" {
		parse, err = ssh.ParsePrivateKeyWithPassphrase(keyBody, []byte(s.KeyFilePassword))
	} else {
		parse, err = ssh.ParsePrivateKey(keyBody)
	}

	signer = ssh.PublicKeys(parse)
	return
}

func (s *SSH) keyBodyAuthFunc() (signer ssh.AuthMethod, err error) {
	var parse ssh.Signer
	if s.KeyFilePassword != "" {
		parse, err = ssh.ParsePrivateKeyWithPassphrase([]byte(s.KeyBody), []byte(s.KeyFilePassword))
	} else {
		parse, err = ssh.ParsePrivateKey([]byte(s.KeyBody))
	}

	signer = ssh.PublicKeys(parse)
	return
}
