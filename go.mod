module github.com/caos/zitadel

go 1.16

require (
	cloud.google.com/go/storage v1.16.0
	github.com/BurntSushi/toml v0.4.1
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v1.0.0-RC2
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/VictoriaMetrics/fastcache v1.6.0
	github.com/ajstarks/svgo v0.0.0-20210406150507-75cfd577ce75
	github.com/allegro/bigcache v1.2.1
	github.com/boombuler/barcode v1.0.1
	github.com/caos/logging v0.0.2
	github.com/caos/oidc v0.15.7
	github.com/caos/orbos v1.5.14-0.20210727080455-c90c315021f5
	github.com/cockroachdb/cockroach-go/v2 v2.1.1
	github.com/duo-labs/webauthn v0.0.0-20210727191636-9f1b88ef44cc
	github.com/envoyproxy/protoc-gen-validate v0.6.1
	github.com/getsentry/sentry-go v0.11.0
	github.com/ghodss/yaml v1.0.0
	github.com/golang/glog v0.0.0-20210429001901-424d2337a529
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/csrf v1.7.1
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/schema v1.2.0
	github.com/gorilla/securecookie v1.1.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.5.0
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/jinzhu/gorm v1.9.16
	github.com/kevinburke/go-types v0.0.0-20210723172823-2deba1f80ba7 // indirect
	github.com/kevinburke/rest v0.0.0-20210506044642-5611499aa33c // indirect
	github.com/kevinburke/twilio-go v0.0.0-20210327194925-1623146bcf73
	github.com/lib/pq v1.10.2
	github.com/lucasb-eyer/go-colorful v1.2.0
	github.com/manifoldco/promptui v0.8.0
	github.com/mattn/go-colorable v0.1.8 // indirect; indirect github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/minio/minio-go/v7 v7.0.12
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/muesli/gamut v0.2.0
	github.com/nicksnyder/go-i18n/v2 v2.1.2
	github.com/pkg/errors v0.9.1
	github.com/pquerna/otp v1.3.0
	github.com/rakyll/statik v0.1.7
	github.com/rs/cors v1.8.0
	github.com/sony/sonyflake v1.0.0
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	github.com/ttacon/builder v0.0.0-20170518171403-c099f663e1c2 // indirect
	github.com/ttacon/libphonenumber v1.2.1
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.22.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.22.0
	go.opentelemetry.io/otel v1.0.0-RC2
	go.opentelemetry.io/otel/exporters/metric/prometheus v0.21.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.0.0-RC2
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.0.0-RC2
	go.opentelemetry.io/otel/metric v0.22.0
	go.opentelemetry.io/otel/sdk v1.0.0-RC2
	go.opentelemetry.io/otel/sdk/export/metric v0.22.0
	go.opentelemetry.io/otel/sdk/metric v0.22.0
	go.opentelemetry.io/otel/trace v1.0.0-RC2
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97
	golang.org/x/oauth2 v0.0.0-20210805134026-6f1e6394065a
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/text v0.3.6
	golang.org/x/tools v0.1.5
	google.golang.org/api v0.52.0
	google.golang.org/genproto v0.0.0-20210805201207-89edb61ffb67
	google.golang.org/grpc v1.39.1
	google.golang.org/protobuf v1.27.1
	gopkg.in/square/go-jose.v2 v2.6.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	gotest.tools v2.2.0+incompatible
	k8s.io/api v0.22.0
	k8s.io/apiextensions-apiserver v0.22.0
	k8s.io/apimachinery v0.22.0
	k8s.io/client-go v0.22.0
	sigs.k8s.io/controller-runtime v0.9.5
)
