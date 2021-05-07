package config

import (
	"time"

	"gopkg.in/mgo.v2"
)

const (
	EnvPath = "local.env"

	MachineRawData     = "iii.dae.MachineRawData"
	MachineRawDataHist = "iii.dae.MachineRawDataHist"
	Statistic          = "iii.dae.Statistics"
	DailyStatistics    = "iii.dae.DailyStatistics"
	MonthlyStatistics  = "iii.dae.MonthlyStatistics"
	YearlyStatistics   = "iii.dae.YearlyStatistics"
	EventLatest        = "iii.dae.EventLatest"
	EventHist          = "iii.dae.EventHist"
	GroupTopo          = "iii.cfg.GroupTopology"
	TPCList            = "iii.cfg.TPCList"
	UserList           = "iii.cfg.UserList"
	EventConfig        = "iii.cfg.EventConfig"
)

var (
	MongodbURL        string
	MongodbDatabase   string
	MongodbUsername   string
	MongodbPassword   string
	TaipeiTimeZone, _ = time.LoadLocation("Asia/Taipei")
	UTCTimeZone, _    = time.LoadLocation("UTC")

	DB      *mgo.Database
	Session *mgo.Session
)
