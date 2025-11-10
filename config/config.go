package config

// RaspiAgentConfig application config that maps env variables
type RaspiAgentConfig struct {
	// App
	ServerAddr string `env:"SERVER_ADDR" envDefault:"0.0.0.0:8080"`
	LogLevel   string `env:"LOG_LEVEL" envDefault:"info"`

	// Auth
	JWTKey string `env:"JWT_KEY" envDefault:""`

	// OAuth2: Google
	GoogleOAuth2ClientID     string `env:"GOOGLE_CLIENT_ID" envDefault:""`
	GoogleOAuth2ClientSecret string `env:"GOOGLE_CLIENT_SECRET" envDefault:""`
	GoogleOAuth2RedirectURL  string `env:"GOOGLE_CLIENT_REDIRECT" envDefault:""`

	// Open AI
	OpenAIAPIKey string `env:"OPENAI_API_KEY" envDefault:""`
	OpenAIAPIURL string `env:"OPENAI_API_URL" envDefault:"https://api.openai.com/v1"`

	// Step CA
	StepCAURL              string `env:"STEPCA_URL" envDefault:""`
	StepCAProvisionerName  string `env:"STEPCA_PROVISIONER_NAME" envDefault:""`
	StepCAProvisionerToken string `env:"STEPCA_PROVISIONER_TOKEN" envDefault:""`
	StepCAPEM              string `env:"STEPCA_PROVISIONER_PEM" envDefault:""`
	StepCAJWK              string `env:"STEPCA_PROVISIONER_JWK" envDefault:""`

	// Postgres
	PostgresHost     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	PostgresUser     string `env:"POSTGRES_USER" envDefault:"postgres"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" envDefault:""`
	PostgresDB       string `env:"POSTGRES_DB" envDefault:"postgres"`
	PostgresPort     string `env:"POSTGRES_PORT" envDefault:"5432"`
}
