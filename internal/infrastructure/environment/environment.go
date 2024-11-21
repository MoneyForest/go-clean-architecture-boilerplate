package environment

type Environment struct {
	Port        string `env:"PORT,required"`
	Environment string `env:"ENV,required"`
	DBEnvironment
	RedisEnvironment
	SQSEnvironment
}

type DBEnvironment struct {
	DBHost     string `env:"DB_HOST,required"`
	DBPort     string `env:"DB_PORT,required"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBDatabase string `env:"DB_DATABASE,required"`
}

type RedisEnvironment struct {
	RedisHost     string `env:"REDIS_HOST,required"`
	RedisPort     string `env:"REDIS_PORT,required"`
	RedisPassword string `env:"REDIS_PASSWORD,required"`
}

type SQSEnvironment struct {
	AWSRegion          string `env:"AWS_REGION,required"`
	AWSEndpoint        string `env:"AWS_ENDPOINT,required"`
	SQSQueueNameSample string `env:"SQS_QUEUE_NAME_SAMPLE,required"`
}
