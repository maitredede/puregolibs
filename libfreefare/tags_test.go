package libfreefare

import (
	"testing"

	"github.com/maitredede/puregolibs/libnfc"
)

func TestGetTagInfos(t *testing.T) {
	nv := libnfc.Version()
	fv := Version()

	t.Logf("libnfc: %s", nv)
	t.Logf("libfreefare: %s", fv)

	context, err := libnfc.InitContext()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		context.Close()
	})

	device, err := context.OpenDefault()
	if err != nil {
		t.Fatal(err)
	}

	tags, err := GetTagInfos(device)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("found %d tags", len(tags))
	for _, tag := range tags {
		t.Logf("- %+v", tag)
	}
}
