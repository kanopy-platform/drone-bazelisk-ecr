package main

import (
	"os"
	"reflect"
	"testing"
)

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
				"PLUGIN_TAG":        "tag",
				"PLUGIN_TARGET":     "target",
				"PLUGIN_REGISTRY":   "registry",
				"PLUGIN_REPOSITORY": "repository",
				"PLUGIN_ACCESS_KEY": "access",
				"PLUGIN_SECRET_KEY": "secret",
				"PLUGIN_BAZELRC":    ".bazelrc.custom",
			},
			want: plugin{
				Tag:        "tag",
				Target:     "target",
				Registry:   "registry",
				Repository: "repository",
				AccessKey:  "access",
				SecretKey:  "secret",
				Bazelrc:    ".bazelrc.custom",
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
