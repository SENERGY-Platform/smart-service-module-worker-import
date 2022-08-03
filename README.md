## Outputs

### Import-Id

- Desc: id of created import
- Variable-Name: import_id

## Camunda-Input-Variables

### Module-Data

- Desc: sets fields for Module.ModuleData
- Variable-Name-Template: `{{config.WorkerParamPrefix}}.module_data`
- Variable-Name-Example: `import.module_data`
- Value: `json.Marshal(map[string]interface{})`

### Request

- Desc: request forwarded to import-deploy to create import
- Variable-Name-Template: `{{config.WorkerParamPrefix}}.request`
- Variable-Name-Example: `import.request`
- Value: json.Marshal(Instance{})

### Config Overwrite

- Desc: overwrites config value with json value. useful to set camunda placeholder as value of integer config.
- Variable-Name-Template: `{{config.WorkerParamPrefix}}.config.json_overwrite.{{request.config.name}}`
- Variable-Name-Example: `import.config.json_overwrite.foo`
- Value: json.Marshal(interface{})