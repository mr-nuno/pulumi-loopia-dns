package provider

import (
	"context"
	"fmt"

	"github.com/pulumi/pulumi-go-provider/infer"
)

type DnsRecordArgs struct {
	Zone  string `pulumi:"zone" validate:"required"`
	Name  string `pulumi:"name" validate:"required"`
	Type  string `pulumi:"type" validate:"required"`
	Value string `pulumi:"value" validate:"required"`
	TTL   int    `pulumi:"ttl,omitempty"`
}

type DnsRecordOutputs struct {
	DnsRecordArgs
	RecordId string `pulumi:"recordId"`
}

type DnsRecord struct {
	getClient ClientFactory
}

func (r *DnsRecord) Create(ctx context.Context, req infer.CreateRequest[DnsRecordArgs]) (infer.CreateResponse[DnsRecordOutputs], error) {
	inputs := req.Inputs
	cfgVal := infer.GetConfig[Config](ctx)
	client, err := r.getClient(ctx, cfgVal)
	if err != nil {
		return infer.CreateResponse[DnsRecordOutputs]{}, fmt.Errorf("failed to create client: %w", err)
	}

	records, fetchErr := client.GetZoneRecords(inputs.Zone, inputs.Name)
	if fetchErr != nil {
		return infer.CreateResponse[DnsRecordOutputs]{}, fmt.Errorf("failed to get DNS records via Loopia API: %w", fetchErr)
	}
	for _, rec := range records {
		typeStr, _ := rec["type"].(string)
		rdataStr, _ := rec["rdata"].(string)
		ttlInt, _ := rec["ttl"].(int)
		if typeStr == inputs.Type && rdataStr == inputs.Value && ttlInt == inputs.TTL {
			idInt, _ := rec["record_id"].(int)
			id := fmt.Sprintf("%s:%s:%s:%d", inputs.Zone, inputs.Name, inputs.Type, idInt)
			output := DnsRecordOutputs{
				DnsRecordArgs: inputs,
				RecordId:      id,
			}
			return infer.CreateResponse[DnsRecordOutputs]{
				ID:     output.RecordId,
				Output: output,
			}, nil
		}
	}

	recordObj := map[string]interface{}{
		"type":     inputs.Type,
		"ttl":      inputs.TTL,
		"priority": 0,
		"rdata":    inputs.Value,
	}
	err = client.AddZoneRecord(inputs.Zone, inputs.Name, recordObj)
	if err != nil {
		return infer.CreateResponse[DnsRecordOutputs]{}, fmt.Errorf("failed to create DNS record via Loopia API: %w", err)
	}
	id := fmt.Sprintf("%s:%s:%s:%s", inputs.Zone, inputs.Name, inputs.Type, inputs.Value)
	output := DnsRecordOutputs{
		DnsRecordArgs: inputs,
		RecordId:      id,
	}
	return infer.CreateResponse[DnsRecordOutputs]{
		ID:     output.RecordId,
		Output: output,
	}, nil
}

func (r *DnsRecord) Read(ctx context.Context, req infer.ReadRequest[DnsRecordArgs, DnsRecordOutputs]) (infer.ReadResponse[DnsRecordArgs, DnsRecordOutputs], error) {
	cfgVal := infer.GetConfig[Config](ctx)
	client, err := r.getClient(ctx, cfgVal)
	if err != nil {
		return infer.ReadResponse[DnsRecordArgs, DnsRecordOutputs]{}, fmt.Errorf("failed to create client: %w", err)
	}
	zone := req.Inputs.Zone
	subdomain := req.Inputs.Name
	records, err := client.GetZoneRecords(zone, subdomain)
	if err != nil {
		return infer.ReadResponse[DnsRecordArgs, DnsRecordOutputs]{}, fmt.Errorf("failed to get DNS records via Loopia API: %w", err)
	}
	var found map[string]interface{}
	for _, rec := range records {
		typeStr, _ := rec["type"].(string)
		rdataStr, _ := rec["rdata"].(string)
		if typeStr == req.Inputs.Type && rdataStr == req.Inputs.Value {
			found = rec
			break
		}
	}
	if found == nil {
		return infer.ReadResponse[DnsRecordArgs, DnsRecordOutputs]{ID: ""}, nil
	}
	currentInputs := DnsRecordArgs{
		Zone:  zone,
		Name:  subdomain,
		Type:  req.Inputs.Type,
		Value: req.Inputs.Value,
		TTL:   found["ttl"].(int),
	}
	return infer.ReadResponse[DnsRecordArgs, DnsRecordOutputs]{
		ID:     req.ID,
		Inputs: currentInputs,
	}, nil
}

func (r *DnsRecord) Update(ctx context.Context, req infer.UpdateRequest[DnsRecordArgs, DnsRecordOutputs]) (infer.UpdateResponse[DnsRecordOutputs], error) {
	inputs := req.Inputs
	old := req.Inputs
	cfgVal := infer.GetConfig[Config](ctx)
	client, err := r.getClient(ctx, cfgVal)
	if err != nil {
		return infer.UpdateResponse[DnsRecordOutputs]{}, fmt.Errorf("failed to create client: %w", err)
	}
	zone := inputs.Zone
	subdomain := inputs.Name
	records, err := client.GetZoneRecords(zone, subdomain)
	if err != nil {
		return infer.UpdateResponse[DnsRecordOutputs]{}, fmt.Errorf("failed to get DNS records via Loopia API: %w", err)
	}
	for _, rec := range records {
		typeStr, _ := rec["type"].(string)
		rdataStr, _ := rec["rdata"].(string)
		idInt, _ := rec["record_id"].(int)
		if typeStr == old.Type && rdataStr == old.Value {
			_ = client.RemoveZoneRecord(zone, subdomain, idInt)
		}
	}
	recordObj := map[string]interface{}{
		"type":     inputs.Type,
		"ttl":      inputs.TTL,
		"priority": 0,
		"rdata":    inputs.Value,
	}
	err = client.AddZoneRecord(zone, subdomain, recordObj)
	if err != nil {
		return infer.UpdateResponse[DnsRecordOutputs]{}, fmt.Errorf("failed to update DNS record via Loopia API: %w", err)
	}
	id := fmt.Sprintf("%s:%s:%s:%s", inputs.Zone, inputs.Name, inputs.Type, inputs.Value)
	output := DnsRecordOutputs{
		DnsRecordArgs: inputs,
		RecordId:      id,
	}
	return infer.UpdateResponse[DnsRecordOutputs]{
		Output: output,
	}, nil
}

func (r *DnsRecord) Delete(ctx context.Context, req infer.DeleteRequest[DnsRecordOutputs]) error {
	old := req.State
	cfgVal := infer.GetConfig[Config](ctx)
	client, err := r.getClient(ctx, cfgVal)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	zone := old.Zone
	subdomain := old.Name
	records, err := client.GetZoneRecords(zone, subdomain)
	if err != nil {
		return fmt.Errorf("failed to get DNS records via Loopia API: %w", err)
	}
	for _, rec := range records {
		typeStr, _ := rec["type"].(string)
		rdataStr, _ := rec["rdata"].(string)
		idInt, _ := rec["record_id"].(int)
		if typeStr == old.Type && rdataStr == old.Value {
			_ = client.RemoveZoneRecord(zone, subdomain, idInt)
		}
	}
	return nil
}
