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
	"google.golang.org/grpc"
	"io"
	"reflect"
)

var  GrpcClientClass    = reflect.TypeOf((*grpc.ClientConn)(nil))   // *grpc.ClientConn

var ClientScannerClass = reflect.TypeOf((*ClientScanner)(nil)).Elem()

type ClientScanner interface {

	ClientBeans() []interface{}

}

var CommandClass = reflect.TypeOf((*Command)(nil)).Elem()

type Command interface {
	glue.NamedBean

	Run(args []string) error

	Desc() string
}

var ControlClientClass = reflect.TypeOf((*ControlClient)(nil)).Elem()

type ControlClient interface {
	glue.InitializingBean
	glue.DisposableBean

	Status() (string, error)

	Shutdown(restart bool) (string, error)

	ConfigCommand(command string, args []string) (string, error)

	CertificateCommand(command string, args []string) (string, error)

	JobCommand(command string, args []string) (string, error)

	StorageCommand(command string, args []string) (string, error)

	StorageConsole(writer io.StringWriter, errWriter io.StringWriter) error
}
