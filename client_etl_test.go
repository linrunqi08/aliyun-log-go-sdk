package sls

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

func TestETLJobV2(t *testing.T) {
	suite.Run(t, new(ETLJobTestV2Suite))
}

type ETLJobTestV2Suite struct {
	suite.Suite
	endpoint           string
	projectName        string
	logstoreName       string
	accessKeyID        string
	accessKeySecret    string
	targetLogstoreName string
	client             *Client
}

func (s *ETLJobTestV2Suite) SetupTest() {
	s.endpoint = os.Getenv("LOG_TEST_ENDPOINT")
	s.projectName = os.Getenv("LOG_TEST_PROJECT")
	s.logstoreName = os.Getenv("LOG_TEST_LOGSTORE")
	s.targetLogstoreName = os.Getenv("LOG_TEST_TARGET_LOGSTORE")
	s.accessKeyID = os.Getenv("LOG_TEST_ACCESS_KEY_ID")
	s.accessKeySecret = os.Getenv("LOG_TEST_ACCESS_KEY_SECRET")
	s.client = &Client{
		AccessKeyID:     s.accessKeyID,
		AccessKeySecret: s.accessKeySecret,
		Endpoint:        s.endpoint,
	}
}

func (s *ETLJobTestV2Suite) createETLJobV2(etlName string) error {
	sink := ETLSink{
		AccessKeyId:     s.accessKeyID,
		AccessKeySecret: s.accessKeySecret,
		Endpoint:        s.endpoint,
		Logstore:        s.targetLogstoreName,
		Name:            "aliyun-etl-test",
		Project:         s.projectName,
	}
	config := ETLConfiguration{
		AccessKeyId:     s.accessKeyID,
		AccessKeySecret: s.accessKeySecret,
		FromTime:        time.Now().Unix(),
		Script:          "e_set('aliyun','new')",
		Version:         2,
		Logstore:        s.logstoreName,
		ETLSinks:        []ETLSink{sink},
		Parameters:      map[string]string{},
	}
	schedule := ETLSchedule{
		Type: "Resident",
	}
	etljob := ETL{
		Configuration: config,
		DisplayName:   "displayName",
		Description:   "go sdk case",
		Name:          etlName,
		Schedule:      schedule,
		Type:          "ETL",
	}
	return s.client.CreateETL(s.projectName, etljob)
}

func (s *ETLJobTestV2Suite) TestClient_UpdateETLJobV2() {
	etlName := "test_update_etl"
	err := s.createETLJobV2(etlName)
	s.Require().Nil(err)
	etljob, err := s.client.GetETL(s.projectName, etlName)
	s.Require().Nil(err)
	etljob.DisplayName = "update"
	etljob.Description = "update description"
	etljob.Configuration.Script = "e_set('update','update')"
	err = s.client.UpdateETL(s.projectName, *etljob)
	s.Require().Nil(err)
	etljob, err = s.client.GetETL(s.projectName, etlName)
	s.Require().Nil(err)
	s.Require().Equal("update", etljob.DisplayName)
	s.Require().Equal("update description", etljob.Description)
	err = s.client.DeleteETL(s.projectName, etlName)
	s.Require().Nil(err)
}

func (s *ETLJobTestV2Suite) TestClient_DeleteETLJobV2() {
	etlName := "test_delete_etl"
	err := s.createETLJobV2(etlName)
	s.Require().Nil(err)
	_, err = s.client.GetETL(s.projectName, etlName)
	s.Require().Nil(err)
	err = s.client.DeleteETL(s.projectName, etlName)
	s.Require().Nil(err)
	time.Sleep(time.Second * 100)
	_, err = s.client.GetETL(s.projectName, etlName)
	s.Require().NotNil(err)
}

func (s *ETLJobTestV2Suite) TestClient_ListETLJobV2() {
	etlName := "test_list_etl"
	err := s.createETLJobV2(etlName)
	s.Require().Nil(err)
	etljobList, err := s.client.ListETL(s.projectName, 0, 100)
	s.Require().Nil(err)
	s.Require().Equal(1, etljobList.Total)
	s.Require().Equal(1, etljobList.Count)
	err = s.client.DeleteETL(s.projectName, etlName)
	s.Require().Nil(err)
}

func (s *ETLJobTestV2Suite) TestClient_StartStopETLJobV2() {
	etlName := "test_start_stop_etl"
	err := s.createETLJobV2(etlName)
	s.Require().Nil(err)
	for {
		etljob, err := s.client.GetETL(s.projectName, etlName)
		s.Require().Nil(err)
		time.Sleep(10 * time.Second)
		if etljob.Status == "RUNNING" {
			break
		}
	}

	err = s.client.StopETL(s.projectName, etlName)
	for {
		etljob, err := s.client.GetETL(s.projectName, etlName)
		s.Require().Nil(err)
		time.Sleep(10 * time.Second)
		if etljob.Status == "STOPPED" {
			break
		}
	}
	err = s.client.StartETL(s.projectName, etlName)
	for {
		etljob, err := s.client.GetETL(s.projectName, etlName)
		s.Require().Nil(err)
		time.Sleep(10 * time.Second)
		if etljob.Status == "RUNNING" {
			break
		}
	}
	err = s.client.DeleteETL(s.projectName, etlName)
	s.Require().Nil(err)
}

func (s *ETLJobTestV2Suite) TestClient_RestartETLJobV2() {
	etlName := "test_restart_etl"
	err := s.createETLJobV2(etlName)
	s.Require().Nil(err)
	for {
		etljob, err := s.client.GetETL(s.projectName, etlName)
		s.Require().Nil(err)
		time.Sleep(10 * time.Second)
		if etljob.Status == "RUNNING" {
			break
		}
	}

	etljob, err := s.client.GetETL(s.projectName, etlName)
	s.Require().Nil(err)
	etljob.DisplayName = "update"
	etljob.Description = "update description"
	etljob.Configuration.Script = "e_set('update','update')"

	err = s.client.RestartETL(s.projectName, *etljob)
	s.Require().Nil(err)

	for {
		time.Sleep(10 * time.Second)
		etljob, err := s.client.GetETL(s.projectName, etlName)
		s.Require().Nil(err)
		if etljob.Status == "RUNNING" {
			break
		}
	}

	etljob, err = s.client.GetETL(s.projectName, etlName)
	s.Require().Nil(err)
	s.Require().Equal("update", etljob.DisplayName)
	s.Require().Equal("update description", etljob.Description)

	err = s.client.DeleteETL(s.projectName, etlName)
	s.Require().Nil(err)
}
