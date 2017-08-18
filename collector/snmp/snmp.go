// +build linux

/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2016 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package snmp

import (
	"fmt"
	"strings"
	"time"

	"github.com/intelsdi-x/snap-plugin-collector-snmp/collector/configReader"
	"github.com/k-sone/snmpgo"
)

func NewHandler(agentConfig configReader.SnmpAgent) (*snmpgo.SNMP, error) {
	handler, err := snmpgo.NewSNMP(snmpgo.SNMPArguments{
		Version:          getSNMPVersion(agentConfig.SnmpVersion),
		Network:          agentConfig.Network,
		Address:          agentConfig.Address,
		Timeout:          time.Duration(agentConfig.Timeout) * time.Second,
		Retries:          agentConfig.Retries,
		Community:        agentConfig.Community,
		UserName:         agentConfig.UserName,
		SecurityLevel:    getSNMPSecurityLevel(agentConfig.SecurityLevel),
		AuthPassword:     agentConfig.AuthPassword,
		AuthProtocol:     getSNMPAuthProtocol(agentConfig.AuthProtocol),
		PrivPassword:     agentConfig.PrivPassword,
		PrivProtocol:     getPrivProtocol(agentConfig.PrivProtocol),
		SecurityEngineId: agentConfig.SecurityEngineId,
		ContextEngineId:  agentConfig.ContextEngineId,
		ContextName:      agentConfig.ContextName,
	})

	if err != nil {
		return nil, err
	}
	return handler, nil
}

func ReadElements(handler *snmpgo.SNMP, oid string, mode string) ([]*snmpgo.VarBind, error) {

	//results received through SNMP requests
	results := []*snmpgo.VarBind{}

	if err := handler.Open(); err != nil {
		// Failed to open connection
		return results, err
	}

	//get elements in node OID
	nodeOid := strings.Trim(oid, ".")
	oidParts := strings.Split(nodeOid, ".")

	//get length of node OID (used to stop reading in table and walk modes)
	nodeOIDLength := len(oidParts)

	//previous OID (used to stop reading in table and walk modes)
	var prevOid string

	//loop through one node of MIB
	for {
		oids, err := snmpgo.NewOids([]string{oid})
		if err != nil {
			// Failed to parse Oids
			return results, err
		}

		var pdu snmpgo.Pdu
		if mode == configReader.ModeSingle {
			pdu, err = handler.GetRequest(oids)
		} else {
			pdu, err = handler.GetNextRequest(oids)
		}
		if err != nil {
			// Failed to request
			return results, err
		}

		if pdu.ErrorStatus() != snmpgo.NoError {
			// Received an error from the agent
			return results, fmt.Errorf("Received an error from the SNMP agent: %v", pdu.ErrorStatus())
		}

		if len(pdu.VarBinds()) != 1 {
			return results, fmt.Errorf("Unaccepted number of results, received %v results", len(pdu.VarBinds()))
		}

		// select a VarBind
		result := pdu.VarBinds()[0]

		if mode == configReader.ModeSingle {
			results = append(results, result)
			break
		} else {
			oid = result.Oid.String()

			//get current elements in node OID
			currOidParts := strings.Split(strings.Trim(oid, "."), ".")

			// if length of new oid is lower then it is the another node
			if len(currOidParts) < nodeOIDLength {
				break
			}

			currNodeOid := strings.Join(currOidParts[:nodeOIDLength], ".")

			//check if there is a new element to read
			if nodeOid != currNodeOid || prevOid == oid ||
				(mode == configReader.ModeTable && (len(oidParts)+1) != len(currOidParts)) {
				break
			}
			prevOid = oid
			results = append(results, result)
		}
	}
	return results, nil
}

func getSNMPVersion(s string) snmpgo.SNMPVersion {
	var snmpVersion snmpgo.SNMPVersion
	switch s {
	case "v1":
		snmpVersion = snmpgo.V1
	case "v2c":
		snmpVersion = snmpgo.V2c
	case "v3":
		snmpVersion = snmpgo.V3
	}
	return snmpVersion
}

func getSNMPSecurityLevel(s string) snmpgo.SecurityLevel {
	var securitylevel snmpgo.SecurityLevel
	switch s {
	case "NoAuthNoPriv":
		securitylevel = snmpgo.NoAuthNoPriv
	case "AuthNoPriv":
		securitylevel = snmpgo.AuthNoPriv
	case "AuthPriv":
		securitylevel = snmpgo.AuthPriv
	}
	return securitylevel
}
func getSNMPAuthProtocol(s string) snmpgo.AuthProtocol {
	var authProtocol snmpgo.AuthProtocol
	switch s {
	case "MD5":
		authProtocol = snmpgo.Md5
	case "SHA":
		authProtocol = snmpgo.Sha
	}
	return authProtocol
}

func getPrivProtocol(s string) snmpgo.PrivProtocol {
	var privProtocol snmpgo.PrivProtocol
	switch s {
	case "DES":
		privProtocol = snmpgo.Des
	case "AES":
		privProtocol = snmpgo.Aes
	}
	return privProtocol
}
