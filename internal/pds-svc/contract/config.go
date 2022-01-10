package contract

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"net/http"
	"os"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nsql"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nval"
)

type Config struct {
	Server      nhttp.ServerConfig
	Client      ClientConfig
	DataSources DataSourcesConfig
	CORS        nhttp.CORSConfig
	SMTP        SMTPConfig
	Firebase    FirebaseConfig
}

func LoadConfig() *Config {
	c := new(Config)

	// Set config server
	port := nval.ParseIntFallback(os.Getenv("PORT"), 3000)
	c.Server.ListenPort = port

	if c.Server.BasePath == "" {
		c.Server.BasePath = os.Getenv("SERVER_BASE_PATH")
	}

	c.Server.Secure = nval.ParseBooleanFallback(os.Getenv("SERVER_LISTEN_SECURE"), true)
	c.Server.TrustProxy = nval.ParseBooleanFallback(os.Getenv("SERVER_TRUST_PROXY"), true)
	c.Server.Debug = nval.ParseBooleanFallback(os.Getenv("DEBUG"), true)

	// Set config client
	c.Client.ClientID = nval.ParseStringFallback(os.Getenv("CLIENT_ID"), "")
	c.Client.ClientSecret = nval.ParseStringFallback(os.Getenv("CLIENT SECRET"), "")

	// Set config data resource
	c.DataSources.Postgres = nsql.Config{
		Driver:          os.Getenv("DB_DRIVER"),
		Host:            os.Getenv("DB_HOST"),
		Port:            os.Getenv("DB_PORT"),
		Username:        os.Getenv("DB_USER"),
		Password:        os.Getenv("DB_PASS"),
		Database:        os.Getenv("DB_NAME"),
		MaxIdleConn:     nsql.NewInt(10),
		MaxOpenConn:     nsql.NewInt(10),
		MaxConnLifetime: nsql.NewInt(1),
	}

	// Load cors
	corsEnabled := nval.ParseBooleanFallback(os.Getenv("CORS_ENABLED"), false)
	if corsEnabled {
		c.CORS = nhttp.CORSConfig{
			Enabled:        true,
			Origins:        nval.ParseStringArrayFallback(os.Getenv("CORS_ORIGINS"), []string{"*"}),
			AllowedHeaders: nval.ParseStringArrayFallback(os.Getenv("CORS_ALLOWED_HEADERS"), []string{"*"}),
			AllowedMethods: nval.ParseStringArrayFallback(os.Getenv("CORS_ALLOWED_METHODS"), []string{http.MethodGet,
				http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodOptions}),
		}
	}

	// Load smtp config
	c.SMTP = SMTPConfig{
		Host:     nval.ParseStringFallback(os.Getenv("SMTP_HOST"), ""),
		Port:     nval.ParseStringFallback(os.Getenv("SMTP_PORT"), "465"),
		Username: nval.ParseStringFallback(os.Getenv("SMTP_USERNAME"), ""),
		Password: nval.ParseStringFallback(os.Getenv("SMTP_PASSWORD"), ""),
	}

	// Load firebase
	c.Firebase.ServiceAccountCredential = nval.ParseStringFallback(os.Getenv("FIREBASE_SERVICE_ACCOUNT_CRED"), "")

	return c
}

func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Server),
		validation.Field(&c.DataSources),
	)
}

type DataSourcesConfig struct {
	Postgres nsql.Config
}

func (c DataSourcesConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Postgres, validation.Required),
	)
}

type ClientConfig struct {
	ClientID     string
	ClientSecret string
}

func (c ClientConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ClientID, validation.Required),
		validation.Field(&c.ClientSecret, validation.Required),
	)
}

type SMTPConfig struct {
	Host     string `json:"SMTP_HOST"`
	Port     string `json:"SMTP_PORT"`
	Username string `json:"SMTP_USERNAME"`
	Password string `json:"SMTP_PASSWORD"`
}

func (c SMTPConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Host, validation.Required),
		validation.Field(&c.Port, validation.Required),
		validation.Field(&c.Username, validation.Required),
		validation.Field(&c.Password, validation.Required),
	)
}

type FirebaseConfig struct {
	ServiceAccountCredential string
}

func (c FirebaseConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ServiceAccountCredential, validation.Required, is.JSON),
	)
}
