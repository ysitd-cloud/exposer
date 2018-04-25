package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/tonyhhyip/vodka"

	"code.ysitd.cloud/component/exposer/internal/manager"
)

func (s *Service) getBond(c *vodka.Context) {
	hostname := c.UserValue("hostname").(string)
	conn, err := s.Syncer.GetBonding(hostname)
	if err != nil {
		c.Error("Error when loading host", http.StatusBadGateway)
		return
	} else if conn == nil {
		c.NotFound()
		return
	}

	c.JSON(http.StatusOK, conn)
}

func (s *Service) listBonds(c *vodka.Context) {
	conn, err := s.Syncer.ListBonding()
	if err != nil {
		c.Error("Error when loading host", http.StatusBadGateway)
		return
	} else if conn == nil {
		c.NotFound()
		return
	}

	c.JSON(http.StatusOK, conn)
}

func (s *Service) createBond(c *vodka.Context) {
	defer c.Request.Body.Close()
	buffer, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Error("Fail in reading body", http.StatusBadRequest)
		return
	}

	var bond Bond
	if err := json.Unmarshal(buffer, &bond); err != nil {
		c.Error("Fail in parse JSON", http.StatusBadRequest)
		return
	}

	mutation := manager.Mutation{
		Action:      manager.Create,
		Hostname:    c.UserValue("hostname").(string),
		ServiceName: bond.Service,
		Port:        bond.Port,
	}

	s.Syncer.GetChannel() <- mutation

	c.Status(http.StatusCreated)
}

func (s *Service) updateBond(c *vodka.Context) {
	defer c.Request.Body.Close()
	buffer, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Error("Fail in reading body", http.StatusBadRequest)
		return
	}

	var bond Bond
	if err := json.Unmarshal(buffer, &bond); err != nil {
		c.Error("Fail in parse JSON", http.StatusBadRequest)
		return
	}

	mutation := manager.Mutation{
		Action:      manager.Update,
		Hostname:    c.UserValue("hostname").(string),
		ServiceName: bond.Service,
		Port:        bond.Port,
	}

	s.Syncer.GetChannel() <- mutation

	c.Status(http.StatusCreated)
}

func (s *Service) removeBond(c *vodka.Context) {

	mutation := manager.Mutation{
		Action:   manager.Delete,
		Hostname: c.UserValue("hostname").(string),
	}

	s.Syncer.GetChannel() <- mutation

	c.Status(http.StatusCreated)
}
