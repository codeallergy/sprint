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
	"crypto/tls"
	"github.com/codeallergy/glue"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"net"
	"net/http"
	"reflect"
)


var (
	TlsConfigClass = reflect.TypeOf((*tls.Config)(nil)) // *tls.Config

	GrpcServerClass    = reflect.TypeOf((*grpc.Server)(nil))   // *grpc.Server
	HealthCheckerClass = reflect.TypeOf((*health.Server)(nil)) // *health.Server
	HttpServerClass    = reflect.TypeOf((*http.Server)(nil))   // *http.Server
)

var ServerScannerClass = reflect.TypeOf((*ServerScanner)(nil)).Elem()

type ServerScanner interface {

	ServerBeans() []interface{}

}

var ServerClass = reflect.TypeOf((*Server)(nil)).Elem()

type Server interface {
	glue.InitializingBean
	glue.DisposableBean

	Bind() error

	Active() bool

	ListenAddress() net.Addr

	Serve() error

	Stop()
}

/**
Page interface for html pages
 */

var PageClass = reflect.TypeOf((*Page)(nil)).Elem()

type Page interface {
	http.Handler

	Pattern() string
}


