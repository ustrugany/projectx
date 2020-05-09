package cassandra

import (
	"fmt"
	"strings"

	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/ustrugany/projectx/api"
)

const (
	MessageTTL = 60 * 5
)

type MessageRepository struct {
	session *gocql.Session
}

func CreateMessageRepository(session *gocql.Session) MessageRepository {
	return MessageRepository{session: session}
}

func (mr MessageRepository) Create(title, content, email string, magicNumber int) (api.Message, error) {
	u := uuid.NewV4()
	message := api.NewMessage(u.String(), title, content, email, magicNumber)
	query := `insert into projectx.message(uuid, title, email, content, magic_number, meta_created_tstamp) values (?, ?, ?, ?, ?, toTimeStamp(toDate(now()))) using ttl ?;`
	if err := mr.session.Query(
		query,
		message.UUID,
		message.Title,
		message.Email,
		message.Content,
		message.MagicNumber,
		MessageTTL,
	).Exec(); err != nil {
		return message, fmt.Errorf("failed to create message: %w", err)
	}

	return message, nil
}

func (mr MessageRepository) GetByUUID(uuid string) (api.Message, error) {
	var message api.Message
	query := `select uuid, title, email, content, magic_number from projectx.message where uuid = ?`
	iterable := mr.session.Query(query, uuid).Iter()
	if iterable.NumRows() > 1 {
		return message, errors.Errorf("data error, expected one row but returned %d", iterable.NumRows())
	}
	iterable.Scan(&message.UUID, &message.Title, &message.Email, &message.Content, &message.MagicNumber)

	return message, nil
}

func (mr MessageRepository) FindByEmail(email string, pageSize, pageToken int) ([]api.Message, error) {
	var (
		sc       gocql.Scanner
		messages []api.Message
	)

	if pageToken == 0 {
		query := `
select uuid, title, email, content, magic_number 
	from projectx.message 
where 
	email = ? 
limit ?
allow filtering
`
		iterable := mr.session.Query(query, email, pageSize).Iter()
		sc = iterable.Scanner()
	} else {
		query := `
select uuid, title, email, content, magic_number 
	from projectx.message 
where 
	email = ?
	and token(magic_number, email) > TOKEN(?, ?)
limit ?
allow filtering
`
		iterable := mr.session.Query(query, email, pageToken, email, pageSize).Iter()
		sc = iterable.Scanner()
	}

	for sc.Next() {
		m := api.Message{}
		err := sc.Scan(&m.UUID, &m.Title, &m.Email, &m.Content, &m.MagicNumber)
		if err != nil {
			continue
		}
		messages = append(messages, m)
	}

	return messages, nil
}

func (mr MessageRepository) FindByMagicNumber(magicNumber int) ([]api.Message, error) {
	var messages []api.Message
	query := "select uuid, title, email, content, magic_number from projectx.message where magic_number = ?"
	iterable := mr.session.Query(query, magicNumber).Iter()
	scanner := iterable.Scanner()
	for scanner.Next() {
		m := api.Message{}
		err := scanner.Scan(&m.UUID, &m.Title, &m.Email, &m.Content, &m.MagicNumber)
		if err != nil {
			continue
		}
		messages = append(messages, m)
	}

	return messages, nil
}

func (mr MessageRepository) Update(uuid, title, content, email string, magicNumber int) (api.Message, error) {
	return api.Message{}, errors.New("not implemented")
}

func (mr MessageRepository) DeleteByUUID(uuid string) (int, error) {
	query := "delete from projectx.m where uuid = ?"
	err := mr.session.Query(query, uuid).Exec()
	if err != nil {
		return 0, fmt.Errorf("failed to delete message: %w", err)
	}

	return 1, nil
}

func (mr MessageRepository) DeleteByUUIDs(uuids []string) (int, error) {
	if len(uuids) == 0 {
		return 0, nil
	}

	query := "delete from projectx.m where uuid in (?)"
	iterable := mr.session.Query(query, strings.Join(uuids, ", ")).Iter()

	return iterable.NumRows(), nil
}
