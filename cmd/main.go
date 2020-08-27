package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"os"

	"bitbucket.org/exzeo-usa/docusign-lambda/pkg/connect"
	"github.com/golang/glog"
	"github.com/jfcote87/esign"
)

type ProgramFlags struct {
	AccountID     string
	AccountToken  string
	JWTConfigFile string
	JWTUser       string
}

type Config struct {
	Credential *esign.OAuth2Credential
}

type App struct {
	cfg *Config
	svc *connect.Service
}

func (app *App) CheckFailure() {
	resp, err := app.svc.EventsListFailures().Do(context.Background())
	if err != nil {
		glog.Fatalf("Errr %s", err)
	}
	glog.Infof("Huh %+v", resp)
}

func NewApp(cfg *Config) *App {

	app := &App{
		cfg: cfg,
		svc: connect.New(cfg.Credential),
	}
	return app
}

var Flags ProgramFlags
var DefaultConfig Config

func init() {
	flag.Set("logtostderr", "true")
	flag.StringVar(&Flags.AccountID, "docusign.id",
		getvar("DOCUSIGN_AccountID", "65744228"), "The docusign account id")
	flag.StringVar(&Flags.AccountToken, "docusign.token",
		getvar("DOCUSIGN_Token", ""), "The docusign token")
	flag.StringVar(&Flags.JWTConfigFile, "docusign.jwt_config_file",
		getvar("DOCUSIGN_JWT_CONFIG_FILE", "./jwt.json"), "The Docusign jwt config file")
	flag.StringVar(&Flags.JWTUser, "docusign.jwt_user",
		getvar("DOCUSIGTN_JWT_USER", ""), "The docusign JWT user")
}

func main() {
	flag.Parse()
	if Flags.AccountID != "" {
		initCredsFromToken()
	} else {
		initCredsFromJWT()
	}
	glog.Infof("Parsed flags to %+v: %+v", Flags, DefaultConfig)
	app := NewApp(&DefaultConfig)
	app.CheckFailure()

}

func getvar(varname string, def string) string {
	if maybe_env := os.Getenv(varname); maybe_env == "" {
		return def
	} else {
		return maybe_env
	}
}

func initCredsFromToken() error {
	if credential := esign.TokenCredential(Flags.AccountToken, true).WithAccountID(Flags.AccountID); credential == nil {
		return errors.New("Unable to create credential")
	} else {
		DefaultConfig.Credential = credential
	}

	return nil
}

func initCredsFromJWT() error {

	if buffer, err := ioutil.ReadFile(Flags.JWTConfigFile); err != nil {
		glog.Errorf("%s open: %v", Flags.JWTConfigFile, err)
		return err
	} else {
		var cfg *esign.JWTConfig
		if err = json.Unmarshal(buffer, &cfg); err != nil {
			glog.Errorf("Failed to decode json (%s): %s", buffer, err)
			return err
		}
		if credential, err := cfg.Credential(Flags.JWTUser, nil, nil); err != nil {
			return err
		} else {
			DefaultConfig.Credential = credential
		}
	}
	return nil
}
