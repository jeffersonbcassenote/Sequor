package types

type SecretDB struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Engine   string `json:"engine"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	DBName   string `json:"dbname"`
}

type SecretDBPG struct {
	UserName            string `json:"username"`
	Password            string `json:"password"`
	Engine              string `json:"engine"`
	Host                string `json:"host"`
	Port                int    `json:"port"`
	DbClusterIdentifier string `json:"dbClusterIdentifier"`
	Database            string `json:"database"`
}

type SecretRabbit struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

type RabbitMessage struct {
	TicketId int
	OwnerId  int
	TenantId int
}

type ReceivedMessage struct {
	PersonId int
	TenantId int
	Database string
}
