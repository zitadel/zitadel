package execution

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	zhttp "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/oidc/sign"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/denylist"
	"github.com/zitadel/zitadel/internal/domain"
	target_domain "github.com/zitadel/zitadel/internal/execution/target"
	"github.com/zitadel/zitadel/internal/repository/execution"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/actions"
)

type ContextInfo interface {
	GetHTTPRequestBody() []byte
	GetContent() interface{}
	SetHTTPResponseBody([]byte) error
}

type GetActiveSigningWebKey func(ctx context.Context) (*jose.JSONWebKey, error)

// CallTargets call a list of targets in order with handling of error and responses
func CallTargets(
	ctx context.Context,
	targets []target_domain.Target,
	info ContextInfo,
	alg crypto.EncryptionAlgorithm,
	activeSigningKey GetActiveSigningWebKey,
	deniedIPList []denylist.AddressChecker,
) (_ any, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	// We make sure the signer and its key are only fetched once per CallTargets call.
	signerOnce := sign.GetSignerOnce(activeSigningKey)
	// Create a map to cache encrypters by key ID to avoid recreating them for each target.
	encrypters := &sync.Map{}

	for _, target := range targets {
		// call the type of target
		resp, err := CallTarget(ctx, target, info, alg, signerOnce, encrypters, deniedIPList)
		// handle error if interrupt is set
		logging.WithFields("instanceID", authz.GetInstance(ctx).InstanceID(), "target", target.GetTargetID()).OnError(err).Error("error calling target")
		if err != nil && target.IsInterruptOnError() {
			return nil, err
		}
		if len(resp) > 0 {
			// error in unmarshalling
			if err := info.SetHTTPResponseBody(resp); err != nil && target.IsInterruptOnError() {
				return nil, err
			}
		}
	}
	return info.GetContent(), nil
}

type ContextInfoRequest interface {
	GetHTTPRequestBody() []byte
}

// CallTarget call the desired type of target with handling of responses
func CallTarget(
	ctx context.Context,
	target target_domain.Target,
	info ContextInfoRequest,
	alg crypto.EncryptionAlgorithm,
	signerOnce sign.SignerFunc,
	encrypters *sync.Map,
	deniedIPList []denylist.AddressChecker,
) (res []byte, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	signingKey, err := target.GetSigningKey(alg)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "EXEC-thiiCh5b", "Errors.Internal")
	}

	if target.GetEndpoint() != "" {
		endpointURL, err := url.Parse(target.GetEndpoint())
		if err != nil {
			return nil, zerrors.ThrowInvalidArgument(err, "EXEC-N5lu09", "Errors.Endpoint.Invalid")
		}
		if err := denylist.IsHostBlocked(deniedIPList, endpointURL, net.LookupIP); err != nil {
			return nil, zerrors.ThrowInvalidArgument(err, "EXEC-N5lu09", "Errors.Endpoint.Denied")
		}
	}

	body, err := payload(ctx, info.GetHTTPRequestBody(), target, signerOnce, encrypters)
	if err != nil {
		return nil, err
	}

	switch target.GetTargetType() {
	// get request, ignore response and return request and error for handling in list of targets
	case target_domain.TargetTypeWebhook:
		return nil, webhook(ctx, target.GetEndpoint(), target.GetTimeout(), body, signingKey)
	// get request, return response and error
	case target_domain.TargetTypeCall:
		return Call(ctx, target.GetEndpoint(), target.GetTimeout(), body, signingKey)
	case target_domain.TargetTypeAsync:
		go func(ctx context.Context, target target_domain.Target, info []byte) {
			if _, err := Call(ctx, target.GetEndpoint(), target.GetTimeout(), info, signingKey); err != nil {
				logging.WithFields("target", target.GetTargetID()).OnError(err).Info(err)
			}
		}(context.WithoutCancel(ctx), target, body)
		return nil, nil
	default:
		return nil, zerrors.ThrowInternal(nil, "EXEC-auqnansr2m", "Errors.Execution.Unknown")
	}
}

func payload(ctx context.Context, payload []byte, target target_domain.Target, signerOnce sign.SignerFunc, encrypters *sync.Map) ([]byte, error) {
	switch target.GetPayloadType() {
	case target_domain.PayloadTypeUnspecified,
		target_domain.PayloadTypeJSON:
		return payload, nil
	case target_domain.PayloadTypeJWT:
		return payloadJWT(ctx, payload, signerOnce)
	case target_domain.PayloadTypeJWE:
		return payloadJWE(ctx, payload, target, signerOnce, encrypters)
	default:
		return payload, nil
	}
}

func payloadJWT(ctx context.Context, payload []byte, signerOnce sign.SignerFunc) ([]byte, error) {
	signer, _, err := signerOnce(ctx)
	if err != nil {
		return nil, err
	}
	sig, err := signer.Sign(payload)
	if err != nil {
		return nil, err
	}
	data, err := sig.CompactSerialize()
	if err != nil {
		return nil, err
	}
	return []byte(data), nil
}

func loadEncrypter(target target_domain.Target, encrypters *sync.Map) (jose.Encrypter, error) {
	if encrypter, ok := encrypters.Load(target.GetEncryptionKeyID()); ok {
		return encrypter.(jose.Encrypter), nil
	}
	encryptionKey := target.GetEncryptionKey()
	if len(encryptionKey) == 0 {
		return nil, zerrors.ThrowInternal(nil, "EXEC-2n8fhs7g", "Errors.Execution.MissingEncryptionKey")
	}
	publicKey, algorithm, err := publicKeyFromBytes(encryptionKey)
	if err != nil {
		return nil, err
	}

	encrypter, err := jose.NewEncrypter(
		jose.A256GCM,
		jose.Recipient{
			Algorithm: algorithm,
			Key:       publicKey,
			KeyID:     target.GetEncryptionKeyID(),
		},
		(&jose.EncrypterOptions{}).
			WithType("JWT").
			WithContentType("JWT"),
	)
	if err != nil {
		return nil, err
	}
	encrypters.Store(target.GetEncryptionKeyID(), encrypter)
	return encrypter, nil
}

func payloadJWE(
	ctx context.Context,
	payload []byte,
	target target_domain.Target,
	signerOnce sign.SignerFunc,
	encrypters *sync.Map,
) ([]byte, error) {
	payload, err := payloadJWT(ctx, payload, signerOnce)
	if err != nil {
		return nil, err
	}
	encrypter, err := loadEncrypter(target, encrypters)
	if err != nil {
		return nil, err
	}

	enc, err := encrypter.Encrypt(payload)
	if err != nil {
		return nil, err
	}
	crypted, err := enc.CompactSerialize()
	if err != nil {
		return nil, err
	}
	return []byte(crypted), nil
}

func publicKeyFromBytes(data []byte) (any, jose.KeyAlgorithm, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, "", zerrors.ThrowInternal(nil, "EXEC-3n8fhs7g", "Errors.Execution.InvalidPublicKey")
	}
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, "", err
	}
	switch pk := publicKey.(type) {
	case *rsa.PublicKey:
		return pk, jose.RSA_OAEP_256, nil
	case *ecdsa.PublicKey:
		return pk, jose.ECDH_ES_A256KW, nil
	default:
		return nil, "", zerrors.ThrowInternal(nil, "EXEC-NKJe2", "Errors.Execution.InvalidPublicKey")
	}
}

// webhook call a webhook, ignore the response but return the errror
func webhook(ctx context.Context, url string, timeout time.Duration, body []byte, signingKey string) error {
	_, err := Call(ctx, url, timeout, body, signingKey)
	return err
}

// Call function to do a post HTTP request to a desired url with timeout
func Call(ctx context.Context, url string, timeout time.Duration, body []byte, signingKey string) (_ []byte, err error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	ctx, span := tracing.NewSpan(ctx)
	defer func() {
		cancel()
		span.EndWithError(err)
	}()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if signingKey != "" {
		req.Header.Set(actions.SigningHeader, actions.ComputeSignatureHeader(time.Now(), body, signingKey))
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return HandleResponse(resp)
}

func HandleResponse(resp *http.Response) ([]byte, error) {
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// Check for success between 200 and 299, redirect 300 to 399 is handled by the client, return error with statusCode >= 400
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		var errorBody ErrorBody
		if err := json.Unmarshal(data, &errorBody); err != nil {
			// if json unmarshal fails, body has no ErrorBody information, so will be taken as successful response
			return data, nil
		}
		if errorBody.ForwardedStatusCode != 0 || errorBody.ForwardedErrorMessage != "" {
			if errorBody.ForwardedStatusCode >= 400 && errorBody.ForwardedStatusCode < 500 {
				return nil, zhttp.HTTPStatusCodeToZitadelError(nil, errorBody.ForwardedStatusCode, "EXEC-reUaUZCzCp", errorBody.ForwardedErrorMessage)
			}
			return nil, zerrors.ThrowPreconditionFailed(nil, "EXEC-bmhNhpcqpF", errorBody.ForwardedErrorMessage)
		}
		// no ErrorBody filled in response, so will be taken as successful response
		return data, nil
	}

	return nil, zerrors.ThrowPreconditionFailed(nil, "EXEC-dra6yamk98", "Errors.Execution.Failed")
}

type ErrorBody struct {
	ForwardedStatusCode   int    `json:"forwardedStatusCode,omitempty"`
	ForwardedErrorMessage string `json:"forwardedErrorMessage,omitempty"`
}

func QueryExecutionTargetsForRequest(
	ctx context.Context,
	fullMethod string,
) []target_domain.Target {
	ctx, span := tracing.NewSpan(ctx)
	defer span.End()

	requestTargets, _ := authz.GetInstance(ctx).ExecutionRouter().GetEventBestMatch(execution.ID(domain.ExecutionTypeRequest, fullMethod))
	return requestTargets
}

func QueryExecutionTargetsForResponse(
	ctx context.Context,
	fullMethod string,
) []target_domain.Target {
	ctx, span := tracing.NewSpan(ctx)
	defer span.End()

	responseTargets, _ := authz.GetInstance(ctx).ExecutionRouter().GetEventBestMatch(execution.ID(domain.ExecutionTypeResponse, fullMethod))
	return responseTargets
}

func QueryExecutionTargetsForFunction(ctx context.Context, function string) []target_domain.Target {
	executionTargets, _ := authz.GetInstance(ctx).ExecutionRouter().GetEventBestMatch(function)
	return executionTargets
}
