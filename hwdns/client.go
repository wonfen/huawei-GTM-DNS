package hwdns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gtm-dns/internal/service/dnsprovider"
)

type Client struct {
	endpoint   string
	signer     *Signer
	httpClient *http.Client
}

func NewClient(ak, sk, endpoint string) *Client {
	return &Client{
		endpoint:   endpoint,
		signer:     NewSigner(ak, sk),
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

const pageSize = 500

func (c *Client) ListZones(ctx context.Context) ([]dnsprovider.Zone, error) {
	var all []Zone
	offset := 0
	for {
		var resp ZoneListResponse
		path := fmt.Sprintf("/v2/zones?limit=%d&offset=%d", pageSize, offset)
		if err := c.get(ctx, path, &resp); err != nil {
			return nil, err
		}
		all = append(all, resp.Zones...)
		if len(all) >= resp.Metadata.TotalCount || len(resp.Zones) == 0 {
			break
		}
		offset += len(resp.Zones)
	}
	return toProviderZones(all), nil
}

func (c *Client) CreateZone(ctx context.Context, req dnsprovider.CreateZoneRequest) (dnsprovider.Zone, error) {
	var resp Zone
	if err := c.post(ctx, "/v2/zones", fromProviderCreateZoneReq(req), &resp); err != nil {
		return dnsprovider.Zone{}, err
	}
	return toProviderZone(resp), nil
}

func (c *Client) DeleteZone(ctx context.Context, zoneID string) error {
	path := fmt.Sprintf("/v2/zones/%s", zoneID)
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

func (c *Client) ListRecordSets(ctx context.Context, zoneID string) ([]dnsprovider.RecordSet, error) {
	var all []RecordSet
	offset := 0
	for {
		var resp RecordSetListResponse
		path := fmt.Sprintf("/v2.1/recordsets?zone_id=%s&limit=%d&offset=%d", zoneID, pageSize, offset)
		if err := c.get(ctx, path, &resp); err != nil {
			return nil, err
		}
		all = append(all, resp.Recordsets...)
		if len(all) >= resp.Metadata.TotalCount || len(resp.Recordsets) == 0 {
			break
		}
		offset += len(resp.Recordsets)
	}
	return toProviderRecordSets(all), nil
}

// SetRecordSetStatus enables or disables a record set via the dedicated Huawei
// status endpoint. status must be "ENABLE" or "DISABLE" — the general update
// endpoint rejects "ACTIVE" with DNS.0315 when transitioning from DISABLE.
func (c *Client) SetRecordSetStatus(ctx context.Context, recordSetID, status string) error {
	path := fmt.Sprintf("/v2.1/recordsets/%s/statuses/set", recordSetID)
	body := struct {
		Status string `json:"status"`
	}{Status: status}
	return c.put(ctx, path, body, nil)
}

func (c *Client) DeleteRecordSet(ctx context.Context, zoneID, recordSetID string) error {
	path := fmt.Sprintf("/v2.1/zones/%s/recordsets/%s", zoneID, recordSetID)
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

func (c *Client) UpdateRecordSet(ctx context.Context, zoneID, id string, req dnsprovider.UpdateRecordSetRequest) (dnsprovider.RecordSet, error) {
	var resp RecordSet
	path := fmt.Sprintf("/v2.1/zones/%s/recordsets/%s", zoneID, id)
	if err := c.put(ctx, path, fromProviderUpdateReq(req), &resp); err != nil {
		return dnsprovider.RecordSet{}, err
	}
	return toProviderRecordSet(resp), nil
}

func (c *Client) CreateRecordSet(ctx context.Context, zoneID string, req dnsprovider.CreateRecordSetRequest) (dnsprovider.RecordSet, error) {
	var resp RecordSet
	path := fmt.Sprintf("/v2.1/zones/%s/recordsets", zoneID)
	if err := c.post(ctx, path, fromProviderCreateReq(req), &resp); err != nil {
		return dnsprovider.RecordSet{}, err
	}
	return toProviderRecordSet(resp), nil
}

// SyncsRecordDescription reports that Huawei's ListRecordSets returns the
// recordset description, so sync treats it as authoritative (overwrites local).
func (c *Client) SyncsRecordDescription() bool { return true }

func (c *Client) get(ctx context.Context, path string, out any) error {
	return c.do(ctx, http.MethodGet, path, nil, out)
}

func (c *Client) put(ctx context.Context, path string, in, out any) error {
	return c.do(ctx, http.MethodPut, path, in, out)
}

func (c *Client) post(ctx context.Context, path string, in, out any) error {
	return c.do(ctx, http.MethodPost, path, in, out)
}

func (c *Client) do(ctx context.Context, method, path string, in, out any) error {
	var bodyBytes []byte
	if in != nil {
		var err error
		bodyBytes, err = json.Marshal(in)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, c.endpoint+path, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if err := c.signer.Sign(req); err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return &dnsprovider.APIError{StatusCode: resp.StatusCode, Body: string(body)}
	}

	if out != nil && resp.StatusCode != http.StatusNoContent {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}
