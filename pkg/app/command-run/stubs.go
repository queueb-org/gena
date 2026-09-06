package commandrun

import "regexp"

var (
	genRx         = regexp.MustCompile(`(?mi)^(?P<gen>deepcopy|conversion|defaulter|applyconfiguration|client|informer|lister|openapi|register|validation)-gen$`)
	allGenerators = []string{
		// helpers
		"deepcopy-gen", "defaulter-gen", "conversion-gen", "validation-gen", "register-gen",
		// clients
		// TODO: add
		// openapi
		// TODO: add
	}

	// genToAnnotation is a simple helper that keeps settings for each
	// generator. TODO: make it configurable on a user side.
	genToAnnotation = map[string]string{
		"":               "+custom-tag=true", // defaults
		"deepcopy-gen":   "+k8s:deepcopy-gen=package",
		"conversion-gen": "+k8s:conversion-gen=",
		"defaulter-gen":  "+k8s:defaulter-gen=TypeMeta",
		"validation-gen": "+k8s:validation-gen=TypeMeta",
		"register-gen":   "+groupName=",
		// OpenAPI
		"openapi-gen": "+k8s:openapi-gen=", // TODO: add better tag parsing.
		// Clients
		"applyconfiguration-gen": "+genclient",
		"client-gen":             "+genclient",
		"lister-gen":             "+genclient",
		"informer-gen":           "+genclient",
	}

	genToOptions = map[string]map[string]string{
		"": {},
		// helpers
		"deepcopy-gen":   {"--output-file": "zz_generated.deepcopy.go"},
		"register-gen":   {"--output-file": "zz_generated.register.go"},
		"defaulter-gen":  {"--output-file": "zz_generated.defaults.go"},
		"conversion-gen": {"--output-file": "zz_generated.conversion.go"},
		"validation-gen": {
			"--output-file":  "zz_generated.validations.go",
			"--readonly-pkg": "time,k8s.io/apimachinery/pkg/apis/meta/v1,k8s.io/apimachinery/pkg/types",
		},
		// clients
		"applyconfiguration-gen": {
			"--output-dir":                   "{{ .AppDir }}/pkg/generated/applyconfiguration",
			"--output-pkg":                   "{{ .App }}/pkg/generated/applyconfiguration",
			"--external-applyconfigurations": "",
			"--openapi-schema":               "",
		},
		"client-gen": {
			"--apply-configuration-package": "{{ .App }}/pkg/generated/applyconfiguration",
			"--output-pkg":                  "{{ .App }}/pkg/generated",
			"--output-dir":                  "{{ .AppDir }}/pkg/generated",
			"--input-base":                  "{{ .AppDir }}",
			"--plural-exceptions":           "",
			"--clientset-name":              "kubernetes",
			"--prefers-protobuf":            "false",
			"--fake-clientset":              "true", // defaults but explicitly set.
			"--input":                       `{{ join .Packages.Related "," }}`,
		},
		"lister-gen": {
			"--output-dir":        "{{ .AppDir }}/pkg/generated/listers",
			"--output-pkg":        "{{ .App }}/pkg/generated/listers",
			"--plural-exceptions": "",
		},
		"informer-gen": {
			"--output-dir":                  "{{ .AppDir }}/pkg/generated/informers",
			"--output-pkg":                  "{{ .App }}/pkg/generated/informers",
			"--listers-package":             "{{ .App }}/pkg/generated/listers",
			"--versioned-clientset-package": "{{ .App }}/pkg/generated/kubernetes",
			"--plural-exceptions":           "",
		},

		// openapi
		"openapi-gen": {
			"--report-filename":        "/tmp/update-openapi.sh.api_violations.b5Ak6T", // todo: Generate shell script name automatically
			"--output-dir":             "{{ .AppDir }}/pkg/generated/openapi",
			"--output-file":            "zz_generated.openapi.go",
			"--output-model-name-file": "zz_generated.openapi.go",
			// OpenAPI also require:
			// "k8s.io/apimachinery/pkg/apis/meta/v1" \
			// "k8s.io/apimachinery/pkg/runtime" \
			// "k8s.io/apimachinery/pkg/version" \
			// "k8s.io/apimachinery/pkg/api/resource" \
			// TODO: generate proper name.
			"--output-pkg": "{{ .App }}/pkg/generated/openapi",
		},
	}
)
