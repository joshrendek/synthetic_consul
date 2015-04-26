package synthetic_consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

type SyntheticConsul struct {
	Host   string
	client *api.Client
}

type SyntheticNode struct {
	registration *api.CatalogRegistration
	Client       *api.Client
}

func NewSyntheticConsul(host string) (SyntheticConsul, error) {
	client, err := api.NewClient(&api.Config{Address: host})
	s := SyntheticConsul{host, client}
	return s, err
}

func (s *SyntheticConsul) CreateNode(name string, address string) (*SyntheticNode, error) {
	catalog := s.client.Catalog()
	registration := &api.CatalogRegistration{
		Node:    name,
		Address: address,
	}

	// datacenter and token acl if we need it later
	write_options := &api.WriteOptions{}
	_, err := catalog.Register(registration, write_options)
	node := &SyntheticNode{registration, s.client}
	return node, err
}

func (n *SyntheticNode) CreateService(id string, name string, port int, address string) error {
	catalog := n.Client.Catalog()
	service_registration := &api.AgentService{
		ID:      id,
		Service: name,
		Port:    port,
		Address: address,
	}

	service_check := &api.AgentCheck{
		Node:        n.registration.Node,
		Name:        fmt.Sprintf("Service '%s' check", name),
		Status:      "passing",
		ServiceID:   id,
		ServiceName: name,
	}

	n.registration.Check = service_check
	n.registration.Service = service_registration
	write_options := &api.WriteOptions{}
	_, err := catalog.Register(n.registration, write_options)
	return err
}

func (n *SyntheticNode) UpdateService(status string) error {
	catalog := n.Client.Catalog()
	n.registration.Check.Status = status
	write_options := &api.WriteOptions{}
	_, err := catalog.Register(n.registration, write_options)
	return err
}
