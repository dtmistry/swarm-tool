package action

import (
	"testing"
)

func TestGetArgs(t *testing.T) {
	flags := []string{
		"copiedfrom=it",
		"test-label=test",
	}
	args, err := GetArgs(flags)
	if err != nil {
		t.Errorf("args should have been created")
	}
	if args.Len() != 2 {
		t.Errorf("args should have been created with len=2")
	}
}

func TestGetArgsErr(t *testing.T) {
	flags := []string{
		"copiedfrom=it",
		"test-label",
	}
	_, err := GetArgs(flags)
	if err == nil {
		t.Errorf("Expected an err")
	}
}

func TestGetMapErr(t *testing.T) {
	flags := []string{
		"copiedfrom=it",
		"test-label",
	}
	_, err := GetMap(flags)
	if err == nil {
		t.Errorf("Expected an err")
	}
}

func TestGetMap(t *testing.T) {
	flags := []string{
		"copiedfrom=it",
		"test-label=test",
	}
	labels, err := GetMap(flags)
	if err != nil {
		t.Errorf("labels should have been created. Got %s", err)
	}
	if len(labels) != 2 {
		t.Errorf("labels should have been created with len=2")
	}
}
