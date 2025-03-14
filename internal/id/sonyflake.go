package id

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/drone/envsubst"
	"github.com/jarcoal/jpath"
	"github.com/sony/sonyflake"
	"github.com/zitadel/logging"
)

type sonyflakeGenerator struct {
	*sonyflake.Sonyflake
}

func (s *sonyflakeGenerator) Next() (string, error) {
	id, err := s.NextID()
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(id, 10), nil
}

var (
	GeneratorConfig    *Config   = nil
	sonyFlakeGenerator Generator = nil
)

// SonyFlakeGenerator creates a new id generator
// the function panics if the generator cannot be created
func SonyFlakeGenerator() Generator {
	if sonyFlakeGenerator == nil {
		sfg := Generator(&sonyflakeGenerator{
			sonyflake.NewSonyflake(sonyflake.Settings{
				MachineID: machineID,
				StartTime: time.Date(2019, 4, 29, 0, 0, 0, 0, time.UTC),
			}),
		})

		sonyFlakeGenerator = sfg
	}

	return sonyFlakeGenerator
}

// the following is a copy of sonyflake (https://github.com/sony/sonyflake/blob/master/sonyflake.go)
// with the change of using the "POD-IP" if no private ip is found
func privateIPv4() (net.IP, error) {
	as, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, a := range as {
		ipnet, ok := a.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}

		ip := ipnet.IP.To4()
		if isPrivateIPv4(ip) {
			return ip, nil
		}
	}

	//change: use "POD_IP"
	ip := net.ParseIP(os.Getenv("POD_IP"))
	if ip == nil {
		return nil, errors.New("no private ip address")
	}
	if ipV4 := ip.To4(); ipV4 != nil {
		return ipV4, nil
	}
	return nil, errors.New("no pod ipv4 address")
}

func isPrivateIPv4(ip net.IP) bool {
	return ip != nil &&
		(ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168)
}

func MachineIdentificationMethod() string {
	if GeneratorConfig.Identification.PrivateIp.Enabled {
		return "Private Ip"
	}

	if GeneratorConfig.Identification.Hostname.Enabled {
		return "Hostname"
	}

	if GeneratorConfig.Identification.Webhook.Enabled {
		return "Webhook"
	}

	return "No machine identification method has been enabled"
}

func machineID() (uint16, error) {
	if GeneratorConfig == nil {
		logging.Panic("cannot create a unique id for the machine, generator has not been configured")
	}

	errors := []string{}
	if GeneratorConfig.Identification.PrivateIp.Enabled {
		ip, err := lower16BitPrivateIP()
		if err == nil {
			return ip, nil
		}
		errors = append(errors, fmt.Sprintf("failed to get Private IP address %s", err))
	}

	if GeneratorConfig.Identification.Hostname.Enabled {
		hn, err := hostname()
		if err == nil {
			return hn, nil
		}
		errors = append(errors, fmt.Sprintf("failed to get Hostname %s", err))
	}

	if GeneratorConfig.Identification.Webhook.Enabled {
		cid, err := metadataWebhookID()
		if err == nil {
			return cid, nil
		}
		errors = append(errors, fmt.Sprintf("failed to query metadata webhook %s", err))
	}

	if len(errors) == 0 {
		errors = append(errors, "No machine identification method enabled.")
	}

	logging.WithFields("errors", strings.Join(errors, ", ")).Panic("none of the enabled methods for identifying the machine succeeded")
	//this return will never happen because of panic one line before
	return 0, nil
}

func lower16BitPrivateIP() (uint16, error) {
	ip, err := privateIPv4()
	if err != nil {
		return 0, err
	}

	return uint16(ip[2])<<8 + uint16(ip[3]), nil
}

func hostname() (uint16, error) {
	host, err := os.Hostname()
	if err != nil {
		return 0, err
	}

	h := fnv.New32()
	_, hashErr := h.Write([]byte(host))
	if hashErr != nil {
		return 0, hashErr
	}

	return uint16(h.Sum32()), nil
}

func metadataWebhookID() (uint16, error) {
	webhook := GeneratorConfig.Identification.Webhook
	url, err := envsubst.EvalEnv(webhook.Url)
	if err != nil {
		url = webhook.Url
	}

	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return 0, err
	}

	if webhook.Headers != nil {
		for key, value := range *webhook.Headers {
			req.Header.Set(key, value)
		}
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 && resp.StatusCode < 600 {
		return 0, fmt.Errorf("metadata endpoint returned an unsuccessful status code %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	data, err := extractMetadataResponse(webhook.JPath, body)
	if err != nil {
		return 0, err
	}

	h := fnv.New32()
	if _, err = h.Write(data); err != nil {
		return 0, err
	}
	return uint16(h.Sum32()), nil
}

func extractMetadataResponse(path *string, data []byte) ([]byte, error) {
	if path != nil {
		jp, err := jpath.NewFromBytes(data)
		if err != nil {
			return nil, err
		}

		results := jp.Query(*path)
		if len(results) == 0 {
			return nil, fmt.Errorf("metadata endpoint response was successful, but JSONPath provided didn't match anything in the response: %s", string(data[:]))
		}

		return json.Marshal(results)
	}

	return data, nil
}
