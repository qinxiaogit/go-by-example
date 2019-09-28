package internal

import (
	"github.com/qinxiaogit/go-by-example/spider/config"
	"github.com/qinxiaogit/go-by-example/spider/internal/metadata/hotso"
	"gopkg.in/mgo.v2"
)

//MongoDB /
type MongoDB struct {
}

var session *mgo.Session

func init() {
	if session == nil {
		if s, err := mgo.Dial(config.GetConfig().Mongo.Host); err != nil {
			panic(err.Error())
		} else {
			session = s
		}
	}
}

//NewMongoDB ...
func NewMongoDB() *MongoDB {
	return &MongoDB{}
}

//OnWeiBoInsert ...
func (m *MongoDB) OnWeiBoInsert(data interface{}) {
	s := session.Copy()
	defer s.Close()
	if err := s.DB("hotso").C("weibo").Insert(data); err != nil {
		panic(err.Error())
	}
}

//OnWeiBoFindOne ...
func (m *MongoDB) OnWeiBoFindOne() *hotso.HotData {
	s := session.Copy()
	defer s.Close()
	var data hotso.HotData
	if err := s.DB("hotso").C("weibo").Find(nil).Sort("-intime").Limit(1).One(&data); err != nil {
		panic(err.Error())
	}
	return &data
}

//OnBaiDuInsert ...
func (m *MongoDB) OnBaiDuInsert(data interface{}) {
	s := session.Copy()
	defer s.Close()
	if err := s.DB("hotso").C("baidu").Insert(data); err != nil {
		panic(err.Error())
	}
}

//OnBaiDuFindOne ...
func (m *MongoDB) OnBaiDuFindOne() *hotso.HotData {
	s := session.Copy()
	defer s.Close()
	var data hotso.HotData
	if err := s.DB("hotso").C("baidu").Find(nil).Sort("-intime").Limit(1).One(&data); err != nil {
		panic(err.Error())
	}
	return &data
}

//OnZhiHuInsert ...
func (m *MongoDB) OnZhiHuInsert(data interface{}) {
	s := session.Copy()
	defer s.Close()
	if err := s.DB("hotso").C("zhihu").Insert(data); err != nil {
		panic(err.Error())
	}
}

//OnZhiHuFindOne ...
func (m *MongoDB) OnZhiHuFindOne() *hotso.HotData {
	s := session.Copy()
	defer s.Close()
	var data hotso.HotData
	if err := s.DB("hotso").C("baidu").Find(nil).Sort("-intime").Limit(1).One(&data); err != nil {
		panic(err)
	}
	return &data

}
