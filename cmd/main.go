package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"os"

	"github.com/golang/glog"
	"github.com/jfcote87/esign"
	"github.com/jfcote87/esign/v2.1/users"
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
	svc *users.Service
}

func (app *App) List() {
	tList, err := app.svc.List().Do(context.Background())
	if err != nil {
		glog.Fatalf("Errr %s", err)
	}
	glog.Infof("Huh %+v", tList)
}

func NewApp(cfg *Config) *App {
	app := &App{
		cfg: cfg,
		svc: users.New(cfg.Credential),
	}
	return app
}

var Flags ProgramFlags
var DefaultConfig Config

func init() {
	flag.Set("logtostderr", "true")
	flag.StringVar(&Flags.AccountID, "docusign.id",
		getvar("DOCUSIGN_AccountID", ""), "The docusign account id")
	flag.StringVar(&Flags.AccountToken, "docusign.token",
		getvar("DOCUSIGN_Token", ""), "The docusign token")
	flag.StringVar(&Flags.JWTConfigFile, "docusign.jwt_config_file",
		getvar("DOCUSIGN_JWT_CONFIG_FILE", "./jwt.json"), "The Docusign jwt config file")
	flag.StringVar(&Flags.JWTUser, "docusign.jwt_user",
		getvar("DOCUSIGTN_JWT_USER", "209fcf56-3fb4-4b18-b07b-1f950b30990f"), "The docusign JWT user")
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
	app.List()

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
