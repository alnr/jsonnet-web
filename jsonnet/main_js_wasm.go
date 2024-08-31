package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"syscall/js"

	"github.com/google/go-jsonnet"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

var vm = jsonnet.MakeVM()

func main() {
	js.Global().Get("window").Set("sonnet", js.FuncOf(
		func(this js.Value, args []js.Value) interface{} {
			vm.ExtReset()
			vm.TLAReset()
			vm.ExtCode("claims", args[0].String())

			jsonStr, err := vm.EvaluateAnonymousSnippet("", args[1].String())
			if err != nil {
				return js.ValueOf(err.Error())
			}

			var x struct {
				Identity map[string]any `json:"identity"`
			}
			if err := json.Unmarshal([]byte(jsonStr), &x); err != nil {
				return js.ValueOf("JSON instance invalid: " + err.Error())
			}
			if len(x.Identity) == 0 {
				return js.ValueOf("JSON instance invalid: identity field not found")
			}
			var buf bytes.Buffer
			_ = json.NewEncoder(&buf).Encode(x.Identity)

			inst, err := jsonschema.UnmarshalJSON(&buf)
			if err != nil {
				return js.ValueOf("JSON instance invalid: " + err.Error())
			}

			schema, err := jsonschema.UnmarshalJSON(strings.NewReader(args[2].String()))
			if err != nil {
				return js.ValueOf("JSON schema invalid (1): " + err.Error())
			}
			compiler := jsonschema.NewCompiler()
			compiler.DefaultDraft(jsonschema.Draft7)
			if err := compiler.AddResource("schema", schema); err != nil {
				return js.ValueOf("JSON schema invalid (2): " + err.Error())
			}
			schemaC, err := compiler.Compile("schema")
			if err != nil {
				return js.ValueOf("JSON schema invalid (3): " + err.Error())
			}
			err = schemaC.Validate(inst)
			if err != nil {
				return js.ValueOf(fmt.Sprintf("Schema validation error: %s\n\n%s", err.Error(), jsonStr))
			}
			return js.ValueOf("Result satisfies Schema 🥰\n\n" + jsonStr)
		},
	))

	select {}
}
