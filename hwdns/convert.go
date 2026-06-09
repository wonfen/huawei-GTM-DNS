package hwdns

import "gtm-dns/internal/service/dnsprovider"

// Wire structs (Zone, RecordSet, *Request in models.go) carry Huawei JSON tags
// and are used only on the wire. These helpers convert to/from the neutral DTOs
// at the package boundary.

func toProviderZone(z Zone) dnsprovider.Zone {
	return dnsprovider.Zone{ID: z.ID, Name: z.Name, ZoneType: z.ZoneType, Status: z.Status}
}

func toProviderZones(in []Zone) []dnsprovider.Zone {
	out := make([]dnsprovider.Zone, len(in))
	for i, z := range in {
		out[i] = toProviderZone(z)
	}
	return out
}

func toProviderRecordSet(r RecordSet) dnsprovider.RecordSet {
	return dnsprovider.RecordSet{
		ID: r.ID, ZoneID: r.ZoneID, Name: r.Name, Type: r.Type,
		TTL: r.TTL, Records: r.Records, Status: r.Status, Line: r.Line,
		Weight: r.Weight, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	}
}

func toProviderRecordSets(in []RecordSet) []dnsprovider.RecordSet {
	out := make([]dnsprovider.RecordSet, len(in))
	for i, r := range in {
		out[i] = toProviderRecordSet(r)
	}
	return out
}

func fromProviderUpdateReq(req dnsprovider.UpdateRecordSetRequest) UpdateRecordSetRequest {
	return UpdateRecordSetRequest{
		Name: req.Name, Type: req.Type, TTL: req.TTL,
		Records: req.Records, Status: req.Status, Weight: req.Weight,
	}
}

func fromProviderCreateReq(req dnsprovider.CreateRecordSetRequest) CreateRecordSetRequest {
	return CreateRecordSetRequest{
		Name: req.Name, Type: req.Type, TTL: req.TTL, Records: req.Records,
		Line: req.Line, Weight: req.Weight, Status: req.Status,
	}
}

func fromProviderCreateZoneReq(req dnsprovider.CreateZoneRequest) CreateZoneRequest {
	return CreateZoneRequest{
		Name: req.Name, ZoneType: req.ZoneType,
		Description: req.Description, Email: req.Email,
	}
}
