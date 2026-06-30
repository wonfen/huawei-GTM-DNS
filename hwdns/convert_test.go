package hwdns

import (
	"encoding/json"
	"strings"
	"testing"

	"gtm-dns/internal/service/dnsprovider"
)

func TestFromProviderUpdateReq_DescriptionThreeState(t *testing.T) {
	// nil = preserve → omitted from wire body
	got := fromProviderUpdateReq(dnsprovider.UpdateRecordSetRequest{Name: "a", Type: "A"})
	b, _ := json.Marshal(got)
	if strings.Contains(string(b), "description") {
		t.Fatalf("nil description must be omitted, got %s", b)
	}

	// &"" = clear → emit "description":""
	empty := ""
	got = fromProviderUpdateReq(dnsprovider.UpdateRecordSetRequest{Name: "a", Type: "A", Description: &empty})
	b, _ = json.Marshal(got)
	if !strings.Contains(string(b), `"description":""`) {
		t.Fatalf("empty-string pointer must emit description:\"\", got %s", b)
	}

	// &"x" = set
	x := "note"
	got = fromProviderUpdateReq(dnsprovider.UpdateRecordSetRequest{Name: "a", Type: "A", Description: &x})
	b, _ = json.Marshal(got)
	if !strings.Contains(string(b), `"description":"note"`) {
		t.Fatalf("set must emit description:\"note\", got %s", b)
	}
}

func TestToProviderRecordSet_CopiesDescription(t *testing.T) {
	r := RecordSet{ID: "1", Name: "a", Type: "A", Description: "hello"}
	if toProviderRecordSet(r).Description != "hello" {
		t.Fatal("toProviderRecordSet must copy Description")
	}
}
