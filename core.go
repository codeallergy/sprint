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
	"context"
	"crypto"
	"crypto/x509"
	"encoding/json"
	"github.com/codeallergy/glue"
	"github.com/codeallergy/sprintpb"
	"github.com/codeallergy/store"
	"github.com/codeallergy/uuid"
	"golang.org/x/crypto/acme/autocert"
	htmlTemplate "html/template"
	"io"
	"net"
	"reflect"
	textTemplate "text/template"
	"time"
)


var ResourceServiceClass = reflect.TypeOf((*ResourceService)(nil)).Elem()

type ResourceService interface {

	/*
	Gets resource by name
	 */
	GetResource(name string) ([]byte, error)

	/*
	Gets text template resource by name
	 */
	TextTemplate(name string) (*textTemplate.Template, error)

	/*
	Gets html template resource by name
	 */
	HtmlTemplate(name string) (*htmlTemplate.Template, error)

	/*
	Gets using licences of imported modules
	 */
	GetLicenses(name string) (string, error)

	/*
	Gets open api swagger JSON files for resource source
	 */
	GetOpenAPI(source string) string
}


var ConfigRepositoryClass = reflect.TypeOf((*ConfigRepository)(nil)).Elem()

type ConfigRepository interface {
	glue.DisposableBean
	glue.PropertyResolver

	/**
	Gets property value the property name (key) if found or default value.

	If property not found in storage then function will return empty string with no error.

	In case of issue function will return error.
	*/

	Get(key string) (string, error)

	/**
	Enumerates all properties that start with provided prefix.
	Prefix could be an empty string, in this case all properties will be enumerated.

	On each call callback function should return true to continue enumeration.

	In case of issue function will return error.
	*/

	EnumerateAll(prefix string, cb func(key, value string) bool) error

	/**
	Sets specific string property with key.
	If value is empty string, then the property would be removed from config storage.
	All properties are stored in string values on backend.

	In case of issue function will return error.
	*/

	Set(key, value string) error

	/**
	Watch updates with prefix on backend system during specific active context.

	On each call callback function should return true to continue watching on changes.

	If property deleted, then watching value will be empty for the specific key.

	In case of issue function will return error.

	Use Application as context.
	*/

	Watch(context context.Context, prefix string, cb func(key, value string) bool) (context.CancelFunc, error)

	/**
	Gets backend using for storing properties
	*/

	Backend() store.DataStore

	/**
	Sets backend using for storing properties
	*/
	SetBackend(storage store.DataStore)

}

var CertificateRepositoryClass = reflect.TypeOf((*CertificateRepository)(nil)).Elem()

type CertificateRepository interface {
	glue.DisposableBean

	/**
	Self Signer CRUID
	*/
	SaveSelfSigner(self *sprintpb.SelfSigner) error

	FindSelfSigner(name string) (*sprintpb.SelfSigner, error)

	ListSelfSigners(prefix string, cb func(*sprintpb.SelfSigner) bool) error

	DeleteSelfSigner(name string) error

	/**
	Acme Account CRUID
	 */
	SaveAccount(account *sprintpb.AcmeAccount) error

	FindAccount(email string) (*sprintpb.AcmeAccount, error)

	ListAccounts(prefix string, cb func(*sprintpb.AcmeAccount) bool) error

	DeleteAccount(email string) error

	/**
	Domain zone CRUID
	 */

	SaveZone(zone *sprintpb.Zone) error

	FindZone(zone string) (*sprintpb.Zone, error)

	ListZones(prefix string, cb func(*sprintpb.Zone) bool) error

	DeleteZone(zone string) error

	/**
	Watch zone changes
	 */

	Watch(ctx context.Context, cb func(zone, event string) bool) (cancel context.CancelFunc, err error)

	/**
	Gets backend using for storing certificates
	*/

	Backend() store.DataStore

	/**
	Sets backend using for storing certificates
	*/
	SetBackend(storage store.DataStore)
}

var CertificateServiceClass = reflect.TypeOf((*CertificateService)(nil)).Elem()

type AcmeAccount struct {
	Status string `json:"status,omitempty"`
	Contact []string `json:"contact,omitempty"`
	TermsOfServiceAgreed bool `json:"termsOfServiceAgreed,omitempty"`
	Orders string `json:"orders,omitempty"`
	OnlyReturnExisting bool `json:"onlyReturnExisting,omitempty"`
	ExternalAccountBinding json.RawMessage `json:"externalAccountBinding,omitempty"`
}

type AcmeResource struct {
	Body  AcmeAccount `json:"body,omitempty"`
	URI   string       `json:"uri,omitempty"`
}

type AcmeUser struct {
	Email         string
	Registration  *AcmeResource
	PrivateKey    crypto.PrivateKey
}

type CertificateService interface {
	glue.InitializingBean

	CreateAcmeAccount(email string) error

	GetOrCreateAcmeUser(email string) (user *AcmeUser, err error)

	CreateSelfSigner(cn string, withInter bool) error

	RenewCertificate(zone string) error

	ExecuteCommand(cmd string, args []string) (string, error)

	IssueAcmeCertificate(entry *sprintpb.Zone) (string, error)

	IssueSelfSignedCertificate(entry *sprintpb.Zone) error

}

var AutocertStorageClass = reflect.TypeOf((*AutocertStorage)(nil)).Elem()

type AutocertStorage interface {

	Cache(serverName string) autocert.Cache

}

var AutoupdateServiceClass = reflect.TypeOf((*AutoupdateService)(nil)).Elem()

type AutoupdateService interface {
	glue.InitializingBean
	glue.DisposableBean

	Freeze(jobName string) int64

	Unfreeze(handle int64)

	FreezeJobs() map[int64]string

}

var NodeServiceClass = reflect.TypeOf((*NodeService)(nil)).Elem()

type NodeService interface {
	glue.InitializingBean
	Component

	NodeId() uint64

	NodeIdHex() string

	Issue() uuid.UUID

	Parse(uuid.UUID) (timestampMillis int64, nodeId int64, clock int)
}

type StorageConsoleStream interface {
	Send(*sprintpb.StorageConsoleResponse) error

	Recv() (*sprintpb.StorageConsoleRequest, error)
}

type Record struct {
	Key   []byte
	Value []byte
}

var StorageServiceClass = reflect.TypeOf((*StorageService)(nil)).Elem()

type StorageService interface {
	glue.InitializingBean

	Execute(name, query string, cb func(string) bool) error

	ExecuteCommand(cmd string, args []string) (string, error)

	Console(stream StorageConsoleStream) error

	LocalConsole(writer io.StringWriter, errWriter io.StringWriter) error

}

var JobServiceClass = reflect.TypeOf((*JobService)(nil)).Elem()

type JobInfo struct {
	Name         string
	Schedule     string
	ExecutionFn  func(context.Context) error
}

type JobService interface {

	ListJobs() ([]string, error)

	AddJob(*JobInfo) error

	CancelJob(name string) error

	RunJob(ctx context.Context, name string) error

	ExecuteCommand(cmd string, args []string) (string, error)

}

var CertificateIssuerClass = reflect.TypeOf((*CertificateIssuer)(nil)).Elem()
var CertificateIssuerServiceClass = reflect.TypeOf((*CertificateIssueService)(nil)).Elem()

type CertificateDesc struct {
	Organization string
	Country      string
	Province     string
	City         string
	Street       string
	Zip          string
}

type IssuedCertificate interface {

	KeyFileContents() []byte

	CertFileContents() []byte

	PrivateKey() crypto.Signer

	Certificate() *x509.Certificate

}

type CertificateIssuer interface {

	Parent() (CertificateIssuer, bool)

	Certificate() IssuedCertificate

	IssueInterCert(cn string) (CertificateIssuer, error)

	IssueClientCert(cn string, password string) (cert IssuedCertificate, pfxData []byte, err error)

	IssueServerCert(cn string, domains []string, ipAddresses []net.IP) (IssuedCertificate, error)

}

type CertificateIssueService interface {

	LoadCertificateDesc() (*CertificateDesc, error)

	CreateIssuer(cn string, info *CertificateDesc) (CertificateIssuer, error)

	LoadIssuer(*sprintpb.SelfSigner) (CertificateIssuer, error)

	LocalIPAddresses(addLocalhost bool) ([]net.IP, error)

}

var AuthenticationServiceClass = reflect.TypeOf((*AuthorizationMiddleware)(nil)).Elem()

type AuthorizedUser struct {
	Username   string
	Roles      map[string]bool
	Context    map[string]string
	ExpiresAt  int64
	Token      string
}

type AuthorizationMiddleware interface {
	glue.InitializingBean

	Authenticate(ctx context.Context) (context.Context, error)

	GetUser(ctx context.Context) (*AuthorizedUser, bool)

	HasUserRole(ctx context.Context, role string) bool

	UserContext(ctx context.Context, name string) (string, bool)

	GenerateToken(user *AuthorizedUser) (string, error)

	ParseToken(token string) (*AuthorizedUser, error)

	InvalidateToken(token string)

}

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

var NatServiceClass = reflect.TypeOf((*NatService)(nil)).Elem()

type NatService interface {

	AllowMapping() bool

	AddMapping(protocol string, extport, intport int, name string, lifetime time.Duration) error

	DeleteMapping(protocol string, extport, intport int) error

	ExternalIP() (net.IP, error)

	ServiceName() string
}


var DynDNSServiceClass = reflect.TypeOf((*DynDNSService)(nil)).Elem()

type DynDNSService interface {
	glue.NamedBean
	glue.InitializingBean

	EnsureAllPublic(subDomains ...string) error

	EnsureCustom(func(client DNSProviderClient, zone string, externalIP string) error) error

}


var MailServiceClass = reflect.TypeOf((*MailService)(nil)).Elem()

type Mail struct {
	Sender  string
	Recipients []string
	Subject string
	TextTemplate string
	HtmlTemplate string   // optional
	Data interface{}
	Attachments []string  // optional
}

type MailService interface {
	glue.NamedBean

	SendMail(mail *Mail, timeout time.Duration, async bool) error

}

