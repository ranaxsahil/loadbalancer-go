package health

import (
	"sync"

	"github.com/ranaxsahil/loadbalancer-go/app/loadbalancer"
)

type passiveServers struct {
	server  string
	health  heatlhStatus
	failure int
	success int
}

type passiveHealthChecks struct {
	servers            []passiveServers
	findServer         map[string]bool
	errorCodes         map[string]bool
	consicutiveFails   int
	consicutiveSuccess int
	mu                 sync.Mutex
}

func NewPassiveChecks() *passiveHealthChecks {
	return &passiveHealthChecks{findServer: make(map[string]bool), errorCodes: make(map[string]bool)}
}

// for NOW taking input as a variables Will Later make it as a Struct of CONSICUTIVEfails CONSICUTIVEsuccess and ERRORcodes
// data struct ??
func (p *passiveHealthChecks) Init(data *loadbalancer.Data) error {

	err := p.addServers(data.Servers)
	if err != nil {
		return err
	}

	p.addErrors(data.ErrorCodes)
	p.addConsicutiveFailsSuccess(data.ConsicutiveFails, data.ConsicutiveSuccess)
	return nil
}

func (p *passiveHealthChecks) addServers(servers []string) error {
	if len(servers) == 0 {
		return ErrEmptyInput
	}
	for _, server := range servers {
		if _, found := p.findServer[server]; found { // removing dublicates
			continue
		}
		currServer := passiveServers{server: server, failure: 0, success: 0, health: ACTIVE}

		// saving in the servers array
		p.servers = append(p.servers, currServer)
		p.findServer[server] = true
	}
	return nil
}

func (p *passiveHealthChecks) addErrors(errorCodes []string) {
	for val := range errorCodes {
		p.errorCodes[errorCodes[val]] = true
	}
}

func (p *passiveHealthChecks) addConsicutiveFailsSuccess(consicutiveFails, consicutiveSuccess int) {
	p.consicutiveFails = consicutiveFails
	p.consicutiveSuccess = consicutiveSuccess
}

func (p *passiveHealthChecks) Check() {

	// TODO Check if the request times out if it does make the server fail
	// Check the Response code of the request if the request have any value from errorCodes make the server fail
	// also use ZERO copy method Dont want to copy and make the code slow
}
