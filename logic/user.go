package logic

import (
	"fmt"
	"line-bot-otp-back/db"
	"line-bot-otp-back/model"

	log "github.com/sirupsen/logrus"
)

type UserLigic struct {
	User *model.User
}

func (ul *UserLigic) Create() error {
	log.Debugln("Start crate user")
	user := ul.User

	query := fmt.Sprintf("insert into users (id, name, password, line_uid) values (?, ?, ?, ?)")
	log.Debugln("--- insert query ---")
	log.Debugln(query)
	log.Debugln("-------------------------")
	stmt, err := db.Db.Prepare(query)
	if err != nil {
		log.Errorln("Prepare error: ", err)
		return err
	}

	_, err = stmt.Exec(user.Id, user.Name, user.Password, user.LineUid)
	if err != nil {
		log.Errorln("Exec error: ", err)
		return err
	}

	query = fmt.Sprintf("insert into scores (user_id) values (?)")
	log.Debugln("--- insert query ---")
	log.Debugln(query)
	log.Debugln("-------------------------")
	stmt, err = db.Db.Prepare(query)
	if err != nil {
		log.Errorln("Prepare error: ", err)
		return err
	}

	_, err = stmt.Exec(user.Id)
	if err != nil {
		log.Errorln("Exec error: ", err)
		return err
	}

	return nil
}

func (ul *UserLigic) SelectById() (bool, error) {
	log.Debugln("Start select user")

	query := fmt.Sprintf("select id, name, password, line_uid from users where id = ?")
	log.Debugln("--- select user query ---")
	log.Debugln(query)
	log.Debugln("-------------------------")

	rows, err := db.Db.Query(query, ul.User.Id)
	if err != nil {
		log.Errorln("Exec error: ", err)
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		var lineUid *string
		err = rows.Scan(&ul.User.Id, &ul.User.Name, &ul.User.Password, &lineUid)
		if err != nil {
			log.Errorln("Exec error: ", err)
			return false, err
		}
		ul.User.LineUid = lineUid
	}

	return true, nil
}

func (ul *UserLigic) VaridatePassword(password string) bool {
	return ul.User.Password == password
}

func (ul *UserLigic) IdIsExists() (bool, error) {
	query := fmt.Sprintf("select id from users where id = ?")
	log.Debugln("--- select user query ---")
	log.Debugln(query)
	log.Debugln("-------------------------")

	rows, err := db.Db.Query(query, ul.User.Id)
	if err != nil {
		log.Errorln("Exec error: ", err)
		return false, err
	}
	defer rows.Close()
	count := 0
	var id string
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			log.Errorln("Exec error: ", err)
			return false, err
		}
		count++
	}
	if count == 0 {
		return false, nil
	} else {
		return true, nil
	}
}
