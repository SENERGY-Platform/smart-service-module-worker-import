/*
 * Copyright (c) 2022 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package pkg

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/SENERGY-Platform/smart-service-module-worker-import/pkg/imports"
	lib "github.com/SENERGY-Platform/smart-service-module-worker-lib"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/auth"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/camunda"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/configuration"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/model"
	"github.com/SENERGY-Platform/smart-service-module-worker-lib/pkg/smartservicerepository"
)

func Start(ctx context.Context, wg *sync.WaitGroup, config imports.Config, libConfig configuration.Config) error {
	handlerFactory := func(auth *auth.Auth, smartServiceRepo *smartservicerepository.SmartServiceRepository) (camunda.Handler, error) {
		interval, err := time.ParseDuration(config.HealthCheckInterval)
		if err != nil {
			return nil, err
		}

		handler := imports.New(
			config,
			libConfig,
			auth,
			smartServiceRepo,
		)

		healthCheck := func(module model.SmartServiceModule) (health error, err error) {
			token, err := auth.ExchangeUserToken(module.UserId)
			if err != nil {
				return nil, err
			}
			id, err := getImportId(module.ModuleData)
			if err != nil {
				return nil, err
			}
			code, err := handler.CheckImport(token, id)
			if err != nil {
				return nil, err
			}
			if code >= 300 {
				return fmt.Errorf("import health check returned status-code %v", code), nil
			}
			return nil, nil
		}
		moduleQuery := model.ModulQuery{TypeFilter: &libConfig.CamundaWorkerTopic}
		smartServiceRepo.StartHealthCheck(ctx, interval, moduleQuery, healthCheck) //timer loop
		smartServiceRepo.RunHealthCheck(moduleQuery, healthCheck)                  //initial check

		return handler, nil
	}
	return lib.Start(ctx, wg, libConfig, handlerFactory)
}

func getImportId(moduleData map[string]interface{}) (string, error) {
	imp, ok := moduleData["import"]
	if !ok {
		return "", fmt.Errorf("missing import in module data")
	}
	impObj, ok := imp.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid import in module data")
	}
	id, ok := impObj["id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid import in module data (id is not string)")
	}
	return id, nil
}
