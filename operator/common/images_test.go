package common

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestDockerHubReference(t *testing.T) {

	imgs := []dockerhubImage{CockroachImage, PostgresImage, FlywayImage, AlpineImage}

	type args struct {
		customImageRegistry string
	}
	tests := []struct {
		name string
		args args
		test func(result string) error
	}{{
		name: "Image should be pulled from docker hub by default",
		args: args{
			customImageRegistry: "",
		},
		test: func(result string) error {
			return expectRegistry(result, "")
		},
	}, {
		name: "Given a custom image registry, the registry should be prepended",
		args: args{
			customImageRegistry: "myreg.io",
		},
		test: func(result string) error {
			return expectRegistry(result, "myreg.io")
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := range imgs {
				got := imgs[i].Reference(tt.args.customImageRegistry)
				if err := tt.test(got); err != nil {
					t.Error(fmt.Errorf("DockerHubReference(%s): %w", imgs[i], err))
				}
			}
		})
	}
}

func TestZITADELReference(t *testing.T) {

	imgs := []zitadelImage{ZITADELImage, BackupImage}
	dummyVersion := "v99.99.99"

	type args struct {
		customImageRegistry string
	}
	tests := []struct {
		name string
		args args
		test func(result string) error
	}{{
		name: "Image should be pulled from GHCR by default",
		args: args{
			customImageRegistry: "",
		},
		test: func(result string) error {
			return expectRegistry(result, "ghcr.io/")
		},
	}, {
		name: "Given a random docker hub image and a custom image registry, the registry should be prepended",
		args: args{
			customImageRegistry: "myreg.io",
		},
		test: func(result string) error {
			return expectRegistry(result, "myreg.io")
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := range imgs {
				got := imgs[i].Reference(tt.args.customImageRegistry, dummyVersion)
				if err := tt.test(got); err != nil {
					t.Error(fmt.Errorf("ZITADELReference(%s): %w", imgs[i], err))
				}
			}
		})
	}
}

func expectRegistry(result, expect string) error {
	if !strings.HasPrefix(result, expect) {
		return fmt.Errorf("image is not prefixed by the registry %s", expect)
	}
	points := strings.Count(result[:strings.Index(result, ":")], ".")
	if expect == "" && points > 1 {
		return errors.New("doesn't look like a docker image")
	}

	if points > 1 {
		return errors.New("too many points")
	}
	return nil
}
