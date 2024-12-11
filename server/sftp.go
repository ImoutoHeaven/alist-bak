package server

import (
	"context"
	"github.com/KirCute/sftpd-alist"
	"github.com/alist-org/alist/v3/internal/conf"
	"github.com/alist-org/alist/v3/internal/model"
	"github.com/alist-org/alist/v3/internal/op"
	"github.com/alist-org/alist/v3/internal/setting"
	"github.com/alist-org/alist/v3/pkg/utils"
	"github.com/alist-org/alist/v3/server/ftp"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"net/http"
)

type SftpDriver struct {
	proxyHeader *http.Header
	config      *sftpd.Config
}

func NewSftpDriver() (*SftpDriver, error) {
	header := &http.Header{}
	header.Add("User-Agent", setting.GetStr(conf.FTPProxyUserAgent))
	return &SftpDriver{
		proxyHeader: header,
	}, nil
}

func (d *SftpDriver) GetConfig() *sftpd.Config {
	if d.config != nil {
		return d.config
	}
	serverConfig := ssh.ServerConfig{
		NoClientAuth:         true,
		NoClientAuthCallback: d.NoClientAuth,
		PasswordCallback:     d.PasswordAuth,
		AuthLogCallback:      d.AuthLogCallback,
		BannerCallback:       d.GetBanner,
	}
	for _, k := range conf.SSHSigners {
		serverConfig.AddHostKey(k)
	}
	d.config = &sftpd.Config{
		ServerConfig: serverConfig,
		HostPort:     conf.Conf.SFTP.Listen,
		LogFunc:      utils.Log.Info,
	}
	return d.config
}

func (d *SftpDriver) GetFileSystem(sc *ssh.ServerConn) (sftpd.FileSystem, error) {
	userObj, err := op.GetUserByName(sc.User())
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, "user", userObj)
	ctx = context.WithValue(ctx, "meta_pass", "")
	ctx = context.WithValue(ctx, "client_ip", sc.RemoteAddr().String())
	ctx = context.WithValue(ctx, "proxy_header", d.proxyHeader)
	return &ftp.SftpDriverAdapter{FtpDriver: ftp.NewAferoAdapter(ctx)}, nil
}

func (d *SftpDriver) Close() {
}

func (d *SftpDriver) NoClientAuth(conn ssh.ConnMetadata) (*ssh.Permissions, error) {
	if conn.User() != "guest" {
		return nil, errors.New("only guest is allowed to login without authorization")
	}
	guest, err := op.GetGuest()
	if err != nil {
		return nil, err
	}
	if guest.Disabled || !guest.CanFTPAccess() {
		return nil, errors.New("user is not allowed to access via SFTP")
	}
	return nil, nil
}

func (d *SftpDriver) PasswordAuth(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
	userObj, err := op.GetUserByName(conn.User())
	if err != nil {
		return nil, err
	}
	passHash := model.StaticHash(string(password))
	if err = userObj.ValidatePwdStaticHash(passHash); err != nil {
		return nil, err
	}
	if userObj.Disabled || !userObj.CanFTPAccess() {
		return nil, errors.New("user is not allowed to access via SFTP")
	}
	return nil, nil
}

func (d *SftpDriver) AuthLogCallback(conn ssh.ConnMetadata, method string, err error) {
	ip := conn.RemoteAddr().String()
	if err == nil {
		utils.Log.Infof("[SFTP] %s(%s) logged in via %s", conn.User(), ip, method)
	} else if method != "none" {
		utils.Log.Infof("[SFTP] %s(%s) tries logging in via %s but with error: %s", conn.User(), ip, method, err)
	}
}

func (d *SftpDriver) GetBanner(_ ssh.ConnMetadata) string {
	return setting.GetStr(conf.Announcement)
}
