// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bind

import (
	"github.com/qinchende/gofast/fst/cst"
	"net/http"
)

// Binding describes the interface which needs to be implemented for binding the
// data present in the request such as JSON request body, query parameters or
// the form POST.
type Binding interface {
	Name() string
	Bind(*http.Request, interface{}) error
}

// BindingBody adds BindBody method to Binding. BindBody is similar with Bind,
// but it reads the body from supplied bytes instead of req.Body.
type BindingBody interface {
	Binding
	BindBody([]byte, interface{}) error
}

// BindingUri adds BindUri method to Binding. BindUri is similar with Bind,
// but it read the Params.
type BindingUri interface {
	Name() string
	BindUri(map[string][]string, interface{}) error
}

// StructValidator is the minimal interface which needs to be implemented in
// order for it to be used as the validator engine for ensuring the correctness
// of the request. Gin provides a default implementation for this using
// https://github.com/go-playground/validator/tree/v8.18.2.
type StructValidator interface {
	// ValidateStruct can receive any kind of type and it should never panic, even if the configuration is not right.
	// If the received type is not a struct, any validation should be skipped and nil must be returned.
	// If the received type is a struct or pointer to a struct, the validation should be performed.
	// If the struct is not valid or the validation itself fails, a descriptive error should be returned.
	// Otherwise nil must be returned.
	ValidateStruct(interface{}) error

	// Engine returns the underlying validator engine which powers the
	// StructValidator implementation.
	Engine() interface{}
}

// Validator is the default validator which implements the StructValidator
// interface. It uses https://github.com/go-playground/validator/tree/v8.18.2
// under the hood.
var Validator StructValidator = &defaultValidator{}

// These implement the Binding interface and can be used to bind the data
// present in the request to struct instances.
var (
	Pms  = pmsBinding{} // 这是个特殊
	JSON = jsonBinding{}
	XML  = xmlBinding{}
	Form = formBinding{}
	//Query         = queryBinding{}
	//FormPost      = formPostBinding{}
	FormMultipart = formMultipartBinding{}
	ProtoBuf      = protobufBinding{}
	MsgPack       = msgpackBinding{}
	YAML          = yamlBinding{}
	//Uri           = uriBinding{}
	//Header        = headerBinding{}
)

// Default returns the appropriate Binding instance based on the HTTP method
// and the content type.
func Default(method, contentType string) Binding {
	if method == http.MethodGet {
		return Form
	}

	switch contentType {
	case cst.MIMEAppJson:
		return JSON
	case cst.MIMEAppXml, cst.MIMEXml:
		return XML
	case cst.MIMEProtoBuf:
		return ProtoBuf
	case cst.MIMEMsgPack, cst.MIMEXMsgPack:
		return MsgPack
	case cst.MIMEYaml:
		return YAML
	case cst.MIMEMultiPostForm:
		return FormMultipart
	default: // case MIMEPOSTForm:
		return Form
	}
}

func validate(obj interface{}) error {
	if Validator == nil {
		return nil
	}
	return Validator.ValidateStruct(obj)
}