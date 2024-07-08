module github.com/zitadel/zitadel

go 1.22.2

require (
	cloud.google.com/go/storage v1.43.0
	github.com/BurntSushi/toml v1.4.0
	github.com/DATA-DOG/go-sqlmock v1.5.2
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v1.24.0
	github.com/Masterminds/squirrel v1.5.4
	github.com/ajstarks/svgo v0.0.0-20211024235047-1546f124cd8b
	github.com/benbjohnson/clock v1.3.5
	github.com/boombuler/barcode v1.0.2
	github.com/brianvoe/gofakeit/v6 v6.28.0
	github.com/cockroachdb/cockroach-go/v2 v2.3.8
	github.com/common-nighthawk/go-figure v0.0.0-20210622060536-734e95fb86be
	github.com/crewjam/saml v0.4.14
	github.com/descope/virtualwebauthn v1.0.2
	github.com/dop251/goja v0.0.0-20240627195025-eb1f15ee67d2
	github.com/dop251/goja_nodejs v0.0.0-20240418154818-2aae10d4cbcf
	github.com/drone/envsubst v1.0.3
	github.com/envoyproxy/protoc-gen-validate v1.0.4
	github.com/fatih/color v1.17.0
	github.com/gabriel-vasile/mimetype v1.4.4
	github.com/go-jose/go-jose/v4 v4.0.2
	github.com/go-ldap/ldap/v3 v3.4.8
	github.com/go-webauthn/webauthn v0.10.2
	github.com/gorilla/csrf v1.7.2
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/schema v1.4.1
	github.com/gorilla/securecookie v1.1.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.20.0
	github.com/h2non/gock v1.2.0
	github.com/hashicorp/golang-lru/v2 v2.0.7
	github.com/improbable-eng/grpc-web v0.15.0
	github.com/jackc/pgx/v5 v5.6.0
	github.com/jarcoal/jpath v0.0.0-20140328210829-f76b8b2dbf52
	github.com/jinzhu/gorm v1.9.16
	github.com/k3a/html2text v1.2.1
	github.com/kevinburke/twilio-go v0.0.0-20240623211326-c7334b537077
	github.com/lucasb-eyer/go-colorful v1.2.0
	github.com/minio/minio-go/v7 v7.0.73
	github.com/mitchellh/mapstructure v1.5.0
	github.com/muesli/gamut v0.3.1
	github.com/muhlemmer/gu v0.3.1
	github.com/muhlemmer/httpforwarded v0.1.0
	github.com/nicksnyder/go-i18n/v2 v2.4.0
	github.com/pquerna/otp v1.4.0
	github.com/rakyll/statik v0.1.7
	github.com/rs/cors v1.11.0
	github.com/santhosh-tekuri/jsonschema/v5 v5.3.1
	github.com/sony/sonyflake v1.2.0
	github.com/spf13/cobra v1.8.1
	github.com/spf13/viper v1.19.0
	github.com/stretchr/testify v1.9.0
	github.com/superseriousbusiness/exifremove v0.0.0-20210330092427-6acd27eac203
	github.com/ttacon/libphonenumber v1.2.1
	github.com/zitadel/logging v0.6.0
	github.com/zitadel/oidc/v3 v3.25.1
	github.com/zitadel/passwap v0.6.0
	github.com/zitadel/saml v0.1.3
	github.com/zitadel/schema v1.3.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.53.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.53.0
	go.opentelemetry.io/otel v1.28.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.28.0
	go.opentelemetry.io/otel/exporters/prometheus v0.50.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.28.0
	go.opentelemetry.io/otel/metric v1.28.0
	go.opentelemetry.io/otel/sdk v1.28.0
	go.opentelemetry.io/otel/sdk/metric v1.28.0
	go.opentelemetry.io/otel/trace v1.28.0
	go.uber.org/mock v0.4.0
	golang.org/x/crypto v0.24.0
	golang.org/x/exp v0.0.0-20240613232115-7f521ea00fb8
	golang.org/x/net v0.26.0
	golang.org/x/oauth2 v0.21.0
	golang.org/x/sync v0.7.0
	golang.org/x/text v0.16.0
	google.golang.org/api v0.187.0
	google.golang.org/genproto/googleapis/api v0.0.0-20240701130421-f6361c86f094
	google.golang.org/grpc v1.65.0
	google.golang.org/protobuf v1.34.2
	sigs.k8s.io/yaml v1.4.0
)

require (
	cloud.google.com/go/auth v0.6.1 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.2 // indirect
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping v0.48.0 // indirect
	github.com/bmatcuk/doublestar/v4 v4.6.1 // indirect
	github.com/crewjam/httperr v0.2.0 // indirect
	github.com/go-chi/chi/v5 v5.0.12 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/go-webauthn/x v0.1.9 // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/go-tpm v0.9.0 // indirect
	github.com/google/pprof v0.0.0-20240409012703-83162a5b38cd // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattermost/xml-roundtrip-validator v0.1.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/twilio/twilio-go v1.22.2 // indirect
	github.com/zenazn/goji v1.0.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	google.golang.org/genproto v0.0.0-20240624140628-dc46fd24d27d // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240701130421-f6361c86f094 // indirect
)

require (
	cloud.google.com/go v0.115.0 // indirect
	cloud.google.com/go/compute/metadata v0.3.0 // indirect
	cloud.google.com/go/iam v1.1.8 // indirect
	cloud.google.com/go/trace v1.10.7 // indirect
	github.com/Azure/go-ntlmssp v0.0.0-20221128193559-754e69321358 // indirect
	github.com/amdonov/xmlsig v0.1.0 // indirect
	github.com/beevik/etree v1.3.0
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/dlclark/regexp2 v1.11.0 // indirect
	github.com/dsoprea/go-exif v0.0.0-20230826092837-6579e82b732d // indirect
	github.com/dsoprea/go-exif/v2 v2.0.0-20230826092837-6579e82b732d // indirect
	github.com/dsoprea/go-iptc v0.0.0-20200610044640-bc9ca208b413 // indirect
	github.com/dsoprea/go-jpeg-image-structure v0.0.0-20221012074422-4f3f7e934102 // indirect
	github.com/dsoprea/go-logging v0.0.0-20200710184922-b02d349568dd // indirect
	github.com/dsoprea/go-photoshop-info-format v0.0.0-20200610045659-121dd752914d // indirect
	github.com/dsoprea/go-png-image-structure v0.0.0-20210512210324-29b889a6093d // indirect
	github.com/dsoprea/go-utility v0.0.0-20221003172846-a3e1774ef349 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/fxamacker/cbor/v2 v2.6.0 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.5 // indirect
	github.com/go-errors/errors v1.5.1 // indirect
	github.com/go-sourcemap/sourcemap v2.1.4+incompatible // indirect
	github.com/go-xmlfmt/xmlfmt v1.1.2 // indirect
	github.com/gofrs/flock v0.8.1 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/golang/geo v0.0.0-20230421003525-6adc56603217 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/uuid v1.6.0
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.5 // indirect
	github.com/gorilla/handlers v1.5.2 // indirect
	github.com/h2non/filetype v1.1.3 // indirect
	github.com/h2non/parth v0.0.0-20190131123155-b4df798d6542 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jonboulle/clockwork v0.4.0
	github.com/kevinburke/go-types v0.0.0-20210723172823-2deba1f80ba7 // indirect
	github.com/kevinburke/rest v0.0.0-20240617045629-3ed0ad3487f0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/muesli/clusters v0.0.0-20200529215643-2700303c1762 // indirect
	github.com/muesli/kmeans v0.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2
	github.com/prometheus/client_golang v1.19.1
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/rs/xid v1.5.0 // indirect
	github.com/russellhaering/goxmldsig v1.4.0 // indirect
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/ttacon/builder v0.0.0-20170518171403-c099f663e1c2 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/xrash/smetrics v0.0.0-20240312152122-5f08fbb34913 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.28.0 // indirect
	go.opentelemetry.io/proto/otlp v1.3.1 // indirect
	golang.org/x/sys v0.21.0
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	nhooyr.io/websocket v1.8.11 // indirect
)
