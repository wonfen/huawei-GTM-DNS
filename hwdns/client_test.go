package hwdns_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gtm-dns/internal/service/dnsprovider"
	"github.com/wonfen/huawei-GTM-DNS/hwdns"
)

func TestListZones(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v2/zones", r.URL.Path)
		assert.Equal(t, "500", r.URL.Query().Get("limit"))
		assert.NotEmpty(t, r.Header.Get("Authorization"))
		json.NewEncoder(w).Encode(hwdns.ZoneListResponse{
			Zones:    []hwdns.Zone{{ID: "z1", Name: "example.com.", ZoneType: "public", Status: "ACTIVE"}},
			Metadata: hwdns.PageMetadata{TotalCount: 1},
		})
	}))
	defer srv.Close()

	client := hwdns.NewClient("test-ak", "test-sk", srv.URL)
	zones, err := client.ListZones(context.Background())
	require.NoError(t, err)
	require.Len(t, zones, 1)
	assert.IsType(t, dnsprovider.Zone{}, zones[0])
	assert.Equal(t, "z1", zones[0].ID)
	assert.Equal(t, "example.com.", zones[0].Name)
}

func TestListRecordSets(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v2.1/recordsets", r.URL.Path)
		assert.Equal(t, "z1", r.URL.Query().Get("zone_id"))
		assert.Equal(t, "500", r.URL.Query().Get("limit"))
		json.NewEncoder(w).Encode(hwdns.RecordSetListResponse{
			Recordsets: []hwdns.RecordSet{
				{ID: "r1", ZoneID: "z1", Name: "www.example.com.", Type: "A", TTL: 300, Records: []string{"1.2.3.4"}, Status: "ACTIVE", Line: "default_view"},
			},
			Metadata: hwdns.PageMetadata{TotalCount: 1},
		})
	}))
	defer srv.Close()

	client := hwdns.NewClient("test-ak", "test-sk", srv.URL)
	records, err := client.ListRecordSets(context.Background(), "z1")
	require.NoError(t, err)
	require.Len(t, records, 1)
	assert.IsType(t, dnsprovider.RecordSet{}, records[0])
	assert.Equal(t, "r1", records[0].ID)
	assert.Equal(t, "default_view", records[0].Line)
}

func TestListRecordSetsPaginated(t *testing.T) {
	// Simulate a zone with 3 records returned across 2 pages (page size 2)
	// This tests the offset-based pagination loop.
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		assert.Equal(t, "/v2.1/recordsets", r.URL.Path)
		assert.Equal(t, "z1", r.URL.Query().Get("zone_id"))
		offset := r.URL.Query().Get("offset")
		if offset == "0" || offset == "" {
			json.NewEncoder(w).Encode(hwdns.RecordSetListResponse{
				Recordsets: []hwdns.RecordSet{
					{ID: "r1", ZoneID: "z1", Name: "a.example.com.", Type: "A", TTL: 300, Records: []string{"1.1.1.1"}, Status: "ACTIVE", Line: "default_view"},
					{ID: "r2", ZoneID: "z1", Name: "b.example.com.", Type: "A", TTL: 300, Records: []string{"2.2.2.2"}, Status: "ACTIVE", Line: "default_view"},
				},
				Metadata: hwdns.PageMetadata{TotalCount: 3},
			})
		} else {
			json.NewEncoder(w).Encode(hwdns.RecordSetListResponse{
				Recordsets: []hwdns.RecordSet{
					{ID: "r3", ZoneID: "z1", Name: "*.example.com.", Type: "A", TTL: 300, Records: []string{"3.3.3.3"}, Status: "ACTIVE", Line: "Abroad"},
				},
				Metadata: hwdns.PageMetadata{TotalCount: 3},
			})
		}
	}))
	defer srv.Close()

	client := hwdns.NewClient("test-ak", "test-sk", srv.URL)
	records, err := client.ListRecordSets(context.Background(), "z1")
	require.NoError(t, err)
	assert.Equal(t, 2, calls, "should have made exactly 2 API calls")
	require.Len(t, records, 3)
	assert.Equal(t, "r1", records[0].ID)
	assert.Equal(t, "r2", records[1].ID)
	assert.Equal(t, "r3", records[2].ID)
	assert.Equal(t, "*.example.com.", records[2].Name, "wildcard record should be fetched")
}

func TestSetRecordSetStatus(t *testing.T) {
	for _, tc := range []struct{ action, wantPath string }{
		{"ENABLE", "/v2.1/recordsets/r1/statuses/set"},
		{"DISABLE", "/v2.1/recordsets/r1/statuses/set"},
	} {
		t.Run(tc.action, func(t *testing.T) {
			var capturedStatus string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodPut, r.Method)
				assert.Equal(t, tc.wantPath, r.URL.Path)
				var body struct {
					Status string `json:"status"`
				}
				json.NewDecoder(r.Body).Decode(&body)
				capturedStatus = body.Status
				w.WriteHeader(http.StatusAccepted)
			}))
			defer srv.Close()

			client := hwdns.NewClient("test-ak", "test-sk", srv.URL)
			err := client.SetRecordSetStatus(context.Background(), "r1", tc.action)
			require.NoError(t, err)
			assert.Equal(t, tc.action, capturedStatus)
		})
	}
}

func TestDeleteRecordSet(t *testing.T) {
	deleted := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/v2.1/zones/z1/recordsets/r1", r.URL.Path)
		deleted = true
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	client := hwdns.NewClient("test-ak", "test-sk", srv.URL)
	err := client.DeleteRecordSet(context.Background(), "z1", "r1")
	require.NoError(t, err)
	assert.True(t, deleted)
}

func TestUpdateRecordSet(t *testing.T) {
	var capturedBody hwdns.UpdateRecordSetRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/v2.1/zones/z1/recordsets/r1", r.URL.Path)
		assert.NotEmpty(t, r.Header.Get("Authorization"))
		json.NewDecoder(r.Body).Decode(&capturedBody)
		json.NewEncoder(w).Encode(hwdns.RecordSet{
			ID: "r1", ZoneID: "z1", Name: capturedBody.Name,
			Type: capturedBody.Type, TTL: capturedBody.TTL,
			Records: capturedBody.Records, Status: "ACTIVE",
		})
	}))
	defer srv.Close()

	client := hwdns.NewClient("test-ak", "test-sk", srv.URL)
	updated, err := client.UpdateRecordSet(context.Background(), "z1", "r1", dnsprovider.UpdateRecordSetRequest{
		Name:    "api.example.com.",
		Type:    "A",
		TTL:     600,
		Records: []string{"5.6.7.8"},
	})
	require.NoError(t, err)
	assert.IsType(t, dnsprovider.RecordSet{}, updated)
	assert.Equal(t, []string{"5.6.7.8"}, updated.Records)
	assert.Equal(t, "api.example.com.", capturedBody.Name)
	assert.Equal(t, 600, capturedBody.TTL)
	assert.Equal(t, "api.example.com.", updated.Name)
}
