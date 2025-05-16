package main

type Config struct {
	Port        string `envconfig:"PORT" default:"8080"`
	JWTKey      string `envconfig:"JWT_KEY" required:"true"`
	DeepSeekAPI string `envconfig:"DEEPSEEK_API" required:"true"`
	Path        struct {
		Patent string `envconfig:"PATENT" required:"true"`
		Report string `envconfig:"REPORT" required:"true"`
	} `envconfig:"PATH" required:"true"`
	DB struct {
		User string `envconfig:"USER" required:"true"`
		Pass string `envconfig:"PASS" required:"true"`
		Name string `envconfig:"NAME" required:"true"`
		Host string `envconfig:"HOST" default:"127.0.0.1"`
		Port string `envconfig:"PORT" default:"3306"`
	} `envconfig:"DB" required:"true"`
}
