package id

import (
	"errors"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	"strconv"
	"time"

	"github.com/sony/sonyflake"
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
	SonyFlakeGenerator = Generator(&sonyflakeGenerator{
		sonyflake.NewSonyflake(sonyflake.Settings{
			MachineID: machineID,
			StartTime: time.Date(2019, 4, 29, 0, 0, 0, 0, time.UTC),
		}),
	})
)

// the following is a copy of sonyflake (https://github.com/sony/sonyflake/blob/master/sonyflake.go)
//with the change of using the "POD-IP" if no private ip is found
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

func machineID() (uint16, error) {
	ip, ipErr := lower16BitPrivateIP()
	if ipErr == nil {
		return ip, nil
	}

	cid, cidErr := cloudRunContainerID()
	if cidErr != nil {
		return 0, fmt.Errorf("neighter found a private ip nor a cloud run container instance id: private ip err: %w, cloud run ip err: %s", ipErr, cidErr.Error())
	}
	return cid, nil
}

func lower16BitPrivateIP() (uint16, error) {
	ip, err := privateIPv4()
	if err != nil {
		return 0, err
	}

	return uint16(ip[2])<<8 + uint16(ip[3]), nil
}

func cloudRunContainerID() (uint16, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		"http://metadata.google.internal/computeMetadata/v1/instance/id",
		nil,
	)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Metadata-Flavor", "Google")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 && resp.StatusCode < 600 {
		return 0, fmt.Errorf("cloud metadata returned an unsuccessful status code %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	h := fnv.New32()
	if _, err = h.Write(body); err != nil {
		return 0, err
	}
	return uint16(h.Sum32()), nil
}
