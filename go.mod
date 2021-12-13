module github.com/caos/zitadel

go 1.16

require (
	cloud.google.com/go/storage v1.18.2
	cloud.google.com/go/trace v1.0.0 // indirect
	github.com/BurntSushi/toml v0.4.1
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v1.0.0
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/Masterminds/squirrel v1.5.2
	github.com/VictoriaMetrics/fastcache v1.8.0
	github.com/ajstarks/svgo v0.0.0-20211024235047-1546f124cd8b
	github.com/allegro/bigcache v1.2.1
	github.com/aws/aws-sdk-go-v2 v1.11.2
	github.com/aws/aws-sdk-go-v2/config v1.11.0
	github.com/aws/aws-sdk-go-v2/credentials v1.6.4
	github.com/aws/aws-sdk-go-v2/service/s3 v1.21.0
	github.com/boombuler/barcode v1.0.1
	github.com/caos/logging v0.0.2
	github.com/caos/oidc v1.0.0
	github.com/caos/orbos v1.5.14-0.20211102124704-34db02bceed2
	github.com/cockroachdb/cockroach-go/v2 v2.2.4
	github.com/dop251/goja v0.0.0-20211129110639-4739a1d10a51
	github.com/dop251/goja_nodejs v0.0.0-20211022123610-8dd9abb0616d
	github.com/duo-labs/webauthn v0.0.0-20210727191636-9f1b88ef44cc
	github.com/envoyproxy/protoc-gen-validate v0.6.2
	github.com/getsentry/sentry-go v0.11.0
	github.com/golang/glog v1.0.0
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/csrf v1.7.1
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/schema v1.2.0
	github.com/gorilla/securecookie v1.1.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.7.1
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/jinzhu/gorm v1.9.16
	github.com/kevinburke/go-types v0.0.0-20210723172823-2deba1f80ba7 // indirect
	github.com/kevinburke/rest v0.0.0-20210506044642-5611499aa33c // indirect
	github.com/kevinburke/twilio-go v0.0.0-20210327194925-1623146bcf73
	github.com/lib/pq v1.10.4
	github.com/lucasb-eyer/go-colorful v1.2.0
	github.com/manifoldco/promptui v0.9.0
	github.com/mattn/go-colorable v0.1.12 // indirect; indirect github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/minio/minio-go/v7 v7.0.16
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/muesli/gamut v0.2.0
	github.com/nicksnyder/go-i18n/v2 v2.1.2
	github.com/nsf/jsondiff v0.0.0-20210926074059-1e845ec5d249
	github.com/pkg/errors v0.9.1
	github.com/pquerna/otp v1.3.0
	github.com/rakyll/statik v0.1.7
	github.com/rs/cors v1.8.0
	github.com/sony/sonyflake v1.0.0
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	github.com/ttacon/builder v0.0.0-20170518171403-c099f663e1c2 // indirect
	github.com/ttacon/libphonenumber v1.2.1
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.27.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.28.0
	go.opentelemetry.io/otel v1.3.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.2.0
	go.opentelemetry.io/otel/exporters/prometheus v0.25.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.2.0
	go.opentelemetry.io/otel/metric v0.26.0
	go.opentelemetry.io/otel/sdk v1.2.0
	go.opentelemetry.io/otel/sdk/export/metric v0.25.0
	go.opentelemetry.io/otel/sdk/metric v0.25.0
	go.opentelemetry.io/otel/trace v1.3.0
	golang.org/x/crypto v0.0.0-20211117183948-ae814b36b871
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20211205182925-97ca703d548d // indirect
	golang.org/x/text v0.3.7
	golang.org/x/tools v0.1.8
	google.golang.org/api v0.61.0
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/square/go-jose.v2 v2.6.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	gotest.tools v2.2.0+incompatible
	k8s.io/api v0.22.4
	k8s.io/apiextensions-apiserver v0.22.2
	k8s.io/apimachinery v0.22.4
	k8s.io/client-go v0.22.4
	sigs.k8s.io/controller-runtime v0.10.3
	sigs.k8s.io/yaml v1.3.0
)

replace github.com/gin-gonic/gin => github.com/gin-gonic/gin v1.7.4
