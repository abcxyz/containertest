// Copyright 2024 The Authors (see AUTHORS file)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mysqltest

// This file wraps the containertest docker implementation.
//
// This file is only intended to be used outside of Google. Inside of Google,
// this file should be replaced with the Google-internal version.

import (
	"fmt"
	"io"

	"github.com/abcxyz/containertest"
)

var nopCloser = io.NopCloser(nil)

func start(conf *config) (ConnInfo, io.Closer, error) {
	driver := &containertest.MySQL{Version: conf.mySQLVersion}
	translatedOpts := make([]containertest.Option, 0, 2)
	if conf.killAfterSec != 0 {
		translatedOpts = append(translatedOpts, containertest.WithKillAfterSeconds(conf.killAfterSec))
	}
	if conf.progressLogger != nil {
		translatedOpts = append(translatedOpts, containertest.WithLogger(LoggerBridge{conf.progressLogger}))
	}

	ci, err := containertest.Start(driver, translatedOpts...)
	if err != nil {
		return ConnInfo{}, nopCloser, fmt.Errorf("failed to start container: %w", err)
	}

	return ConnInfo{
		Username: driver.Username(),
		Password: driver.Password(),
		Hostname: ci.Host,
		Port:     ci.PortMapper(driver.StartupPorts()[0]),
	}, ci, nil
}
