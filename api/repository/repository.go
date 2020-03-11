package repository

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type TZConversion struct {
	TimeZone       string `bson:"timeZone" json:"timeZone"`
	TimeDifference string `bson:"timeDifference" json:"timeDifference"`
}

type Repository struct {
	dbSession    *mgo.Session
	dbServer     string
	dbDatabase   string
	dbCollection string
}

func NewRepository(dbServer string, dbDatabase string, dbCollection string) *Repository {
	repo := new(Repository)
	repo.dbServer = dbServer
	repo.dbDatabase = dbDatabase
	repo.dbCollection = dbCollection

	dbSession, err := mgo.Dial(repo.dbServer)
	if err != nil {
		log.Fatal(err)
	}
	repo.dbSession = dbSession
	return repo
}

func (repo *Repository) Close() {
	repo.dbSession.Clone()
}
func (repo *Repository) newSession() *mgo.Session {
	return repo.dbSession.Clone()
}

func (repo *Repository) FindAll() ([]TZConversion, error) {
	dbSession := repo.dbSession
	defer dbSession.Clone()
	coll := dbSession.DB(repo.dbDatabase).C(repo.dbCollection)

	var tzcs []TZConversion
	err := coll.Find(bson.M{}).All(&tzcs)
	return tzcs, err
}
func (repo *Repository) FindByTimeZone(tz string) (TZConversion, error) {
	dbSession := repo.newSession()
	defer dbSession.Close()

	coll := dbSession.DB(repo.dbDatabase).C(repo.dbCollection)

	var tzc TZConversion
	err := coll.Find(bson.M{"timeZone": tz}).One(&tzc)

	return tzc, err
}

func (repo *Repository) Insert(tzc TZConversion) error {
	dbSession := repo.newSession()
	defer dbSession.Close()
	coll := dbSession.DB(repo.dbDatabase).C(repo.dbCollection)

	err := coll.Insert(&tzc)
	return err
}
func (repo *Repository) Delete(tzc TZConversion) error {
	dbSession := repo.newSession()
	defer dbSession.Close()
	coll := dbSession.DB(repo.dbDatabase).C(repo.dbCollection)

	err := coll.Remove(bson.M{"timeZone": tzc.TimeZone})
	return err
}

func (repo *Repository) Update(tz string, tzc TZConversion) error {
	dbSession := repo.newSession()
	defer dbSession.Close()

	coll := dbSession.DB(repo.dbDatabase).C(repo.dbCollection)

	err := coll.Update(bson.M{"timeZone": tz}, &tzc)
	return err
}
