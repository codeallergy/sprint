/*
 * Copyright (c) 2022-2023 Zander Schwid & Co. LLC.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 */

package sprint

import (
	"github.com/codeallergy/glue"
	"reflect"
)

var WhoisServiceClass = reflect.TypeOf((*WhoisService)(nil)).Elem()

type Whois struct {
	Domain    string
	NServer   []string
	State     string
	Person    string
	Email     string
	Registrar string
	Created   string
	PaidTill  string
}

type WhoisService interface {

	Parse(whoisResp string) *Whois

	Whois(domain string) (string, error)

}

// DNSRecord DNS record representation.
type DNSRecord struct {
	ID        string `json:"id,omitempty"`
	Hostname  string `json:"hostname,omitempty"`
	TTL       int    `json:"ttl,omitempty"`
	Type      string `json:"type,omitempty"`
	Priority  int    `json:"priority,omitempty"`
	Value     string `json:"value,omitempty"`
}

var DNSProviderClientClass = reflect.TypeOf((*DNSProviderClient)(nil)).Elem()

type DNSProviderClient interface {

	GetPublicIP() (addr string, err error)

	GetRecords(zoneID string) ([]*DNSRecord, error)

	CreateRecord(zoneID string, record *DNSRecord) (*DNSRecord, error)

	RemoveRecord(zoneID, recordID string) error

}

var DNSProviderClass = reflect.TypeOf((*DNSProvider)(nil)).Elem()

type DNSProvider interface {
	glue.NamedBean

	Detect(whois *Whois) bool

	RegisterChallenge(legoClient interface{}, token string) error

	NewClient() (DNSProviderClient, error)
}

var DynDNSServiceClass = reflect.TypeOf((*DynDNSService)(nil)).Elem()

type DynDNSService interface {
	glue.NamedBean
	glue.InitializingBean

	EnsureAllPublic(subDomains ...string) error

	EnsureCustom(func(client DNSProviderClient, zone string, externalIP string) error) error

}

