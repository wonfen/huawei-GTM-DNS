package hwdns

// Zone represents a DNS zone from Huawei Cloud
type Zone struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ZoneType string `json:"zone_type"`
	Status   string `json:"status"`
}

type PageMetadata struct {
	TotalCount int `json:"total_count"`
}

type ZoneListResponse struct {
	Zones    []Zone       `json:"zones"`
	Metadata PageMetadata `json:"metadata"`
}

type RecordSet struct {
	ID        string   `json:"id"`
	ZoneID    string   `json:"zone_id"`
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	TTL       int      `json:"ttl"`
	Records   []string `json:"records"`
	Status    string   `json:"status"`
	Line      string   `json:"line"`
	Weight      *int     `json:"weight"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
	Description string   `json:"description,omitempty"`
}

type RecordSetListResponse struct {
	Recordsets []RecordSet  `json:"recordsets"`
	Metadata   PageMetadata `json:"metadata"`
}

type UpdateRecordSetRequest struct {
	Name    string   `json:"name,omitempty"`
	Type    string   `json:"type,omitempty"`
	TTL     int      `json:"ttl,omitempty"`
	Records []string `json:"records,omitempty"`
	Status      string  `json:"status,omitempty"`
	Weight      *int    `json:"weight,omitempty"`
	Description *string `json:"description,omitempty"`
}

type CreateRecordSetRequest struct {
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	TTL     int      `json:"ttl,omitempty"`
	Records []string `json:"records,omitempty"`
	Line    string   `json:"line,omitempty"`
	Weight      *int     `json:"weight,omitempty"`
	Status      string   `json:"status,omitempty"`
	Description string   `json:"description,omitempty"`
}

type CreateZoneRequest struct {
	Name        string `json:"name"`
	ZoneType    string `json:"zone_type"`
	Description string `json:"description,omitempty"`
	Email       string `json:"email,omitempty"`
}
