package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"gitlab.com/zenvia/service/go-agent-lambda.git/internal/awsutils"
	"gitlab.com/zenvia/service/go-agent-lambda.git/internal/types"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/denisenkom/go-mssqldb"

	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, event events.SQSEvent) (string, error) {

	log.Printf("[Info] Starting process....")

	for _, message := range event.Records {

		var receivedMessage types.ReceivedMessage
		json.Unmarshal([]byte(message.Body), &receivedMessage)

		var conn *sql.DB
		var err error

		if receivedMessage.Database == "PG" {
			conn, err = getDatabaseConnectionPG()
		} else {
			conn, err = getDatabaseConnectionSQL()
		}

		if err != nil {
			log.Printf("[Error] Failed to connect to database: %s", err.Error())
			return "[Error] Failed open database connection", nil
		}

		defer conn.Close()

		log.Printf("[Info] Starting process for agent: %s, tenant: %s", strconv.Itoa(receivedMessage.PersonId), strconv.Itoa(receivedMessage.TenantId))

		if receivedMessage.Database == "PG" {
			getTicketsPG(conn, receivedMessage)
		} else {
			getTickets(conn, receivedMessage)
		}
	}

	return "OK", nil
}

func getDatabaseConnectionPG() (conn *sql.DB, err error) {

	secretString, _, err := awsutils.GetSecret(os.Getenv("DB_SECRET_KEY_PG"))

	if err != nil {
		log.Printf("[Error] Failed get connection string: %s", err.Error())
		return
	}

	var secret types.SecretDBPG
	if err = json.Unmarshal([]byte(secretString), &secret); err != nil {
		log.Printf("[Error] Failed to parse connection string: %s", err.Error())
		return
	}

	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		secret.Host,
		secret.Port,
		secret.UserName,
		secret.Password,
		secret.Database)

	conn, err = sql.Open("postgres", connString)

	if err != nil {
		log.Printf("[Error] Opening connection failed: %s", err.Error())
		return
	}

	err = conn.Ping()

	if err != nil {
		log.Printf("[Error] Failed to ping database: %s", err.Error())
		return
	}

	return
}

func getDatabaseConnectionSQL() (conn *sql.DB, err error) {

	secretString, _, err := awsutils.GetSecret(os.Getenv("DB_SECRET_KEY"))

	if err != nil {
		log.Printf("[Error] Failed get connection string: %s", err.Error())
		return
	}

	var secret types.SecretDB
	if err = json.Unmarshal([]byte(secretString), &secret); err != nil {
		log.Printf("[Error] Failed to parse connection string: %s", err.Error())
		return
	}

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s",
		secret.Host,
		secret.UserName,
		secret.Password,
		secret.Port,
		secret.DBName)

	conn, err = sql.Open("mssql", connString)

	if err != nil {
		log.Printf("[Error] Opening connection failed: %s", err.Error())
		return
	}

	err = conn.Ping()

	if err != nil {
		log.Printf("[Error] Failed to ping database: %s", err.Error())
		return
	}

	return
}

func getTickets(conn *sql.DB, receivedMessage types.ReceivedMessage) {

	command := `SELECT Id
				FROM Ticket
				WHERE OwnerId = ?
				AND TenantId = ?
				AND IsDeleted = 0
				AND IsDraft = 0 	
				AND SystemStatus != 3
				AND SystemStatus != 4
				AND SystemStatus != 5`

	stmt, err := conn.Prepare(command)

	if err != nil {
		log.Printf("[Error] Prepare to get ticket failed: %s", err.Error())
	}

	defer stmt.Close()

	rows, err := stmt.Query(receivedMessage.PersonId, receivedMessage.TenantId)

	if err != nil {
		log.Printf("[Error] %s", err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var (
			id int
		)
		if err := rows.Scan(&id); err != nil {
			log.Printf("[Error] %s", err.Error())
		}

		postTicketsToQueue(id, receivedMessage)
	}
}

func getTicketsPG(conn *sql.DB, receivedMessage types.ReceivedMessage) {

	command := `SELECT id FROM dbo.ticket WHERE ownerid = ` + strconv.Itoa(receivedMessage.PersonId) +
		` AND tenantid = ` + strconv.Itoa(receivedMessage.TenantId) +
		` AND isdeleted = false
		    AND isdraft = false
		    AND systemstatus != 3
			AND systemstatus != 4
			AND systemstatus != 5`

	stmt, err := conn.Prepare(command)

	if err != nil {
		log.Printf("[Error] Prepare to get ticket failed: %s", err.Error())
	}

	defer stmt.Close()

	rows, err := stmt.Query()

	if err != nil {
		log.Printf("[Error] %s", err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var (
			id int
		)
		if err := rows.Scan(&id); err != nil {
			log.Printf("[Error] %s", err.Error())
		}

		postTicketsToQueue(id, receivedMessage)
	}
}

func postTicketsToQueue(ticketId int, receivedMessage types.ReceivedMessage) {

	secretString, _, err := awsutils.GetSecret(os.Getenv("RABBITMQ_SECRET_KEY"))

	if err != nil {
		log.Printf("[Error] Failed get RabbitMQ connection string: %s", err.Error())
		return
	}

	var secret types.SecretRabbit
	if err = json.Unmarshal([]byte(secretString), &secret); err != nil {
		log.Printf("[Error] Failed to parse RabbitMQ connection string: %s", err.Error())
		return
	}

	connString := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		secret.UserName,
		secret.Password,
		secret.Host,
		secret.Port)

	connection, err := amqp.Dial(connString)

	if err != nil {
		log.Printf("[Error] Could not establish connection with RabbitMQ: %s", err.Error())
	}

	channel, err := connection.Channel()

	if err != nil {
		log.Printf("[Error] Could not open RabbitMQ channel: %s", err.Error())
	}

	err = channel.ExchangeDeclare("movidesk.ticket.distribuitor", "fanout", true, false, false, false, nil)

	if err != nil {
		log.Printf("[Error] %s", err.Error())
	}

	m := types.RabbitMessage{TicketId: ticketId, OwnerId: receivedMessage.PersonId, TenantId: receivedMessage.TenantId}
	b, err := json.Marshal(m)

	if err != nil {
		log.Printf("[Error] Could not convert Rabbit message: %s", err.Error())
	}

	message := amqp.Publishing{
		Body: b,
	}

	err = channel.Publish("movidesk.ticket.distribuitor", "random-key", false, false, message)

	if err != nil {
		log.Printf("[Error] Error publishing a message to the queue: %s", err.Error())
	}

	_, err = channel.QueueDeclare("movidesk.ticket.distribuitor.person", true, false, false, false, nil)

	if err != nil {
		log.Printf("[Error] Error declaring the queue: %s", err.Error())
	}

	err = channel.QueueBind("movidesk.ticket.distribuitor.person", "#", "movidesk.ticket.distribuitor.person", false, nil)

	if err != nil {
		log.Printf("[Error] Error binding to the queue: %s", err.Error())
	}

	log.Printf("[Info] Sent %d owned by %d \n", ticketId, receivedMessage.PersonId)
}
