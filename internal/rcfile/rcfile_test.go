package rcfile

import (
	"fmt"
	"testing"
)

func TestCheckMirror(t *testing.T) {
	got := checkMirror("mail.com")
	want := false
	if got != want {
		t.Errorf("got %t; want %t", got, want)
	}
}
func TestCheckMirror2(t *testing.T) {
	got := checkMirror("asdfasf.com")
	want := false
	if got != want {
		t.Errorf("got %t; want %t", got, want)
	}
}
func TestTerraformrcFile(t *testing.T) {
	got := string(templateFile("mail.ru"))
	want := testTemplate
	if got != want {
		fmt.Println(len(got), len(want))
		t.Errorf("got %s; want %s", got, want)
	}
}

var testTemplate = `provider_installation {
  network_mirror {
    url = "https://mail.ru/"
    include = ["registry.terraform.io/*/*"]
  }
  direct {
    exclude = ["registry.terraform.io/*/*"]
  }
}`

func TestCreateTerraformrc(t *testing.T) {
	t.Setenv("HOME", "/s/sda")
	mirrors := []string{"terraform-mirror.yandexcloud.net", "registry.comcloud.xyz"}
	got := CreateFile(mirrors)
	if got == nil {
		t.Errorf("got %t; want Error", got)
	}
}

func TestCreateTerraformrc2(t *testing.T) {
	t.Setenv("HOME", "/tmp")
	mirrors := []string{"terraform-mirror.yandexcloud.net", "registry.comcloud.xyz"}
	got := CreateFile(mirrors)
	if got != nil {
		t.Errorf("got %t; want nil", got)
	}
}
