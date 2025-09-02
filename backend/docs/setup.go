package docs

import (
	"dario.cat/mergo"
	"github.com/swaggo/swag"
)

func SetupSwaggerParams(spec swag.Spec) {
	mergo.Merge(SwaggerInfo, spec, mergo.WithOverride)
}
