package contract

type Config struct {
	Debug                   string `envconfig:"DEBUG"`
	Port                    int    `envconfig:"PORT"`
	ServerBasePath          string `envconfig:"SERVER_BASE_PATH"`
	ServerBaseUrl           string `envconfig:"SERVER_BASE_URL"`
	ServerSecure            string `envconfig:"SERVER_LISTEN_SECURE"`
	ServerTrustProxy        string `envconfig:"SERVER_TRUST_PROXY"`
	DatabaseDriver          string `envconfig:"DB_DRIVER"`
	DatabaseHost            string `envconfig:"DB_HOST"`
	DatabasePort            uint16 `envconfig:"DB_PORT"`
	DatabaseUsername        string `envconfig:"DB_USER"`
	DatabasePassword        string `envconfig:"DB_PASS"`
	DatabaseName            string `envconfig:"DB_NAME"`
	DatabaseMaxIdleConn     string `envconfig:"DB_POOL_MAX_IDLE_CONN"`
	DatabaseMaxOpenConn     string `envconfig:"DB_POOL_MAX_OPEN_CONN"`
	DatabaseMaxConnLifetime string `envconfig:"DB_POOL_MAX_CONN_LIFETIME"`
	DatabaseBootMigration   string `envconfig:"DB_BOOT_MIGRATION"`
}
