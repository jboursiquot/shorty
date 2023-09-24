package shorty

type Config struct {
	BaseURL string `env:"BASE_URL,default=http://localhost:8080"`
}
