package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/tonyhhyip/vodka"

	"code.ysitd.cloud/component/exposer/internal/manager"
)

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
