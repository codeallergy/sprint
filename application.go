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
	"flag"
	"github.com/codeallergy/glue"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/http"
	"reflect"
)

var (
	LogClass = reflect.TypeOf((*zap.Logger)(nil)) // *zap.Logger
	LumberjackClass = reflect.TypeOf((*lumberjack.Logger)(nil)) // *lumberjack.Logger

	FileSystemClass = reflect.TypeOf((*http.FileSystem)(nil)).Elem() // http.FileSystem
	FlagSetClass    = reflect.TypeOf((*flag.FlagSet)(nil))           // *flag.FlagSet
)

var CoreScannerClass = reflect.TypeOf((*CoreScanner)(nil)).Elem()

/**
	Core scanner is using to find beans to create core context when application is already created.
*/
type CoreScanner interface {

	/**
	Gets core beans that will extend application beans, parent context.
	 */
	CoreBeans() []interface{}

}


/**
Generic component class that has a name
*/
var ComponentClass = reflect.TypeOf((*Component)(nil)).Elem()

type Component interface {
	glue.NamedBean

	GetStats(cb func(name, value string) bool) error
}


var ApplicationClass = reflect.TypeOf((*Application)(nil)).Elem()

/**
	Application is the base entry point class for golang application.
 */
type Application interface {
	context.Context
	glue.InitializingBean
	glue.NamedBean
	Component

	/**
	Add beans to application context
	 */
	AppendBeans(beans ...interface{})

	/**
	Gets application name, represents lower case normalized executable name
	 */
	Name() string

	/**
	Gets application version at the time of compilation
	*/
	Version() string

	/**
	Gets application build at the time of compilation
	*/
	Build() string

	/**
	Gets application runtime profile, could be: dev, qa, prod and etc.
	 */
	Profile() string

	/**
	Checks if application running in dev mode
	 */
	IsDev() bool

	/**
	Gets application binary name, used on startup, could be different with application name
	 */
	Executable() string

	/**
	Gets home directory of the application, usually parent directory of binary folder where executable is running, not current directory
	 */
	ApplicationDir() string

	/**
	Run application with command line arguments
	 */
	Run(args []string) error

	/**
	Indicator if application is active and not in shutting down mode
	 */
	Active() bool

	/**
	Sets the flag that application is in shutting down mode then notify all go routines by ShutdownChannel then notify signal channel with interrupt signal

	Additionally sets the flag that application is going to be restarted after shutdown
	 */
	Shutdown(restart bool)

	/**
	Indicator if application needs to be restarted by autoupdate or remote command after shutdown
	*/
	Restarting() bool

}

var SystemEnvironmentPropertyResolverClass = reflect.TypeOf((*SystemEnvironmentPropertyResolver)(nil)).Elem()

type SystemEnvironmentPropertyResolver interface {

	PromptProperty(key string) (string, bool)

	Environ(withValues bool) []string

}

var ApplicationFlagsClass = reflect.TypeOf((*ApplicationFlags)(nil)).Elem()

type ApplicationFlags interface {
	glue.PropertyResolver

	Daemon() bool

	Verbose() bool

	Properties() map[string]string

}

var FlagSetRegistrarClass = reflect.TypeOf((*FlagSetRegistrar)(nil)).Elem()

type FlagSetRegistrar interface {
	RegisterFlags(fs *flag.FlagSet)

	RegisterServerArgs(args []string) []string
}

