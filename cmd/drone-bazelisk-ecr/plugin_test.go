package main

import (
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
)

type mockECRClient struct {
	ecriface.ECRAPI
}

var testFailure string

func (m *mockECRClient) GetAuthorizationToken(input *ecr.GetAuthorizationTokenInput) (*ecr.GetAuthorizationTokenOutput, error) {
	if testFailure == "GetAuthorizationToken" {
		return nil, errors.New("GetAuthorizationToken")
	}

	data := &ecr.GetAuthorizationTokenOutput{
		AuthorizationData: []*ecr.AuthorizationData{
			{ProxyEndpoint: aws.String("https://0123456789.dkr.ecr.us-east-1.amazonaws.com")},
		},
	}

	return data, nil
}

func (m *mockECRClient) CreateRepository(input *ecr.CreateRepositoryInput) (*ecr.CreateRepositoryOutput, error) {
	if testFailure == "CreateRepository" {
		return nil, errors.New("CreateRepository")
	}

	if testFailure == "CreateRepositoryRepoExists" {
		return nil, awserr.New(ecr.ErrCodeRepositoryAlreadyExistsException, "", errors.New("CreateRepositoryRepoExists"))
	}

	return &ecr.CreateRepositoryOutput{}, nil
}

func TestGetArgs(t *testing.T) {
	tests := []struct {
		plugin plugin
		want   []string
	}{
		{
			plugin: plugin{Target: "test"},
			want:   []string{"run", "test"},
		},
		{
			plugin: plugin{Target: "test", Bazelrc: ".bazelrc.custom"},
			want:   []string{"--bazelrc=.bazelrc.custom", "run", "test"},
		},
		{
			plugin: plugin{Target: "test", Bazelrc: ".bazelrc.custom", Command: "test"},
			want:   []string{"--bazelrc=.bazelrc.custom", "test", "test"},
		},
		{
			plugin: plugin{Target: "test", Bazelrc: ".bazelrc.custom", CommandArgs: "--config=test"},
			want:   []string{"--bazelrc=.bazelrc.custom", "run", "--config=test", "test"},
		},
		{
			plugin: plugin{Target: "test", Bazelrc: ".bazelrc.custom", TargetArgs: "--var"},
			want:   []string{"--bazelrc=.bazelrc.custom", "run", "test", "--", "--var"},
		},
	}

	for _, test := range tests {
		got := test.plugin.getArgs()
		if !reflect.DeepEqual(test.want, got) {
			t.Errorf("%v is not equal to %v", test.want, got)
		}
	}
}

func TestSetenv(t *testing.T) {
	tests := []struct {
		env  map[string]string
		want plugin
		fail bool
	}{
		// test setting struct fields
		{
			env: map[string]string{
				"PLUGIN_TAG":               "tag",
				"PLUGIN_TARGET":            "target",
				"PLUGIN_REGISTRY":          "registry",
				"PLUGIN_REPOSITORY":        "repository",
				"PLUGIN_ACCESS_KEY":        "access",
				"PLUGIN_SECRET_KEY":        "secret",
				"PLUGIN_BAZELRC":           ".bazelrc.custom",
				"PLUGIN_CREATE_REPOSITORY": "true",
			},
			want: plugin{
				Tag:              "tag",
				Target:           "target",
				Registry:         "registry",
				Repository:       "repository",
				AccessKey:        "access",
				SecretKey:        "secret",
				Bazelrc:          ".bazelrc.custom",
				CreateRepository: true,
			},
			fail: false,
		},
		// test empty environment
		{
			env:  map[string]string{},
			want: plugin{},
			fail: true,
		},
	}

	for _, test := range tests {
		setEnvMap(test.env)

		got := newPlugin()
		err := got.setenv()
		if err != nil && !test.fail {
			t.Errorf(err.Error())
		}

		if test.want != got {
			t.Errorf("%v is not equal to %v", test.want, got)
		}

		unsetEnvMap(test.env)
	}
}

func TestRegion(t *testing.T) {
	tests := []struct {
		p    plugin
		want string
		fail bool
	}{
		{
			p:    plugin{Registry: "0123456789.dkr.ecr.us-east-1.amazonaws.com"},
			want: "us-east-1",
			fail: false,
		},
		{
			p:    plugin{},
			fail: true,
		},
	}

	for _, test := range tests {
		got, err := test.p.region()

		// check for invalid registries
		if err != nil && !test.fail {
			t.Errorf("invalid registry should have failed")
		}

		// check if regions match
		if test.want != got {
			t.Errorf("%v is not equal to %v", test.want, got)
		}
	}
}

func TestCreateRepository(t *testing.T) {
	tests := []struct {
		p       plugin
		failure string
	}{
		// test successful create
		{
			p: plugin{Registry: "0123456789.dkr.ecr.us-east-1.amazonaws.com", Repository: "repository"},
		},
		// test that repo exists error is ignored
		{
			p:       plugin{Registry: "0123456789.dkr.ecr.us-east-1.amazonaws.com", Repository: "repository"},
			failure: "CreateRepositoryRepoExists",
		},
		// test empty repository failure
		{
			p:       plugin{},
			failure: "must specify a repository",
		},
		// test auth token failure
		{
			p:       plugin{Repository: "repository"},
			failure: "GetAuthorizationToken",
		},
		// test for mismatched registry failure
		{
			p:       plugin{Registry: "thisdoesnotmatch.dkr.ecr.us-east-1.amazonaws.com", Repository: "repository"},
			failure: "provided credentials are not for the specified registry",
		},
		// test for generic repo creation failure
		{
			p:       plugin{Registry: "0123456789.dkr.ecr.us-east-1.amazonaws.com", Repository: "repository"},
			failure: "CreateRepository",
		},
	}

	for _, test := range tests {
		// set global failure variable to inform mocks
		testFailure = test.failure

		err := test.p.createRepository(&mockECRClient{})
		if err != nil && !strings.HasPrefix(err.Error(), testFailure) {
			t.Errorf(err.Error())
		}
	}
}

// set environment
func setEnvMap(env map[string]string) {
	for key, val := range env {
		os.Setenv(key, val)
	}
}

// unset environment
func unsetEnvMap(env map[string]string) {
	for key := range env {
		os.Unsetenv(key)
	}
}
