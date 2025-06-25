package provider

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

// Client defines the interface for Loopia DNS API operations.
type Client interface {
	AddZoneRecord(domain, subdomain string, recordObj map[string]interface{}) error // Add a DNS record
	RemoveZoneRecord(domain, subdomain string, recordID int) error                  // Remove a DNS record
	GetZoneRecords(domain, subdomain string) ([]map[string]interface{}, error)      // List DNS records
}

// LoopiaClient implements the Client interface for the Loopia XML-RPC API.
type LoopiaClient struct {
	Username   string       // Loopia API username
	Password   string       // Loopia API password
	Endpoint   string       // Loopia API endpoint
	HTTPClient *http.Client // HTTP client for requests
}

// NewLoopiaClient creates a new LoopiaClient instance.
func NewLoopiaClient(username, password, endpoint string) *LoopiaClient {
	return &LoopiaClient{
		Username:   username,
		Password:   password,
		Endpoint:   endpoint,
		HTTPClient: &http.Client{},
	}
}

// xmlrpcParam, xmlrpcValue, xmlrpcMethodCall are helpers for XML-RPC requests.
type xmlrpcParam struct {
	XMLName xml.Name    `xml:"param"`
	Value   xmlrpcValue `xml:"value"`
}
type xmlrpcValue struct {
	XMLName xml.Name    `xml:"value"`
	Inner   interface{} `xml:",innerxml"`
}
type xmlrpcMethodCall struct {
	XMLName    xml.Name      `xml:"methodCall"`
	MethodName string        `xml:"methodName"`
	Params     []xmlrpcParam `xml:"params>param"`
}

// buildXMLRPCRequest builds an XML-RPC request body for the given method and arguments.
func buildXMLRPCRequest(method string, args ...interface{}) string {
	params := make([]xmlrpcParam, len(args))
	for i, arg := range args {
		switch v := arg.(type) {
		case string:
			params[i] = xmlrpcParam{Value: xmlrpcValue{Inner: fmt.Sprintf("<string>%s</string>", xmlEscape(v))}}
		case int:
			params[i] = xmlrpcParam{Value: xmlrpcValue{Inner: fmt.Sprintf("<int>%d</int>", v)}}
		case map[string]interface{}:
			params[i] = xmlrpcParam{Value: xmlrpcValue{Inner: buildXMLRPCStruct(v)}}
		default:
			params[i] = xmlrpcParam{Value: xmlrpcValue{Inner: fmt.Sprintf("<string>%v</string>", v)}}
		}
	}
	call := xmlrpcMethodCall{
		MethodName: method,
		Params:     params,
	}
	out, _ := xml.MarshalIndent(call, "", "  ")
	return string(out)
}

// buildXMLRPCStruct builds an XML-RPC struct from a Go map.
func buildXMLRPCStruct(obj map[string]interface{}) string {
	var buf bytes.Buffer
	buf.WriteString("<struct>")
	for k, v := range obj {
		buf.WriteString("<member>")
		buf.WriteString("<name>")
		buf.WriteString(xmlEscape(k))
		buf.WriteString("</name>")
		buf.WriteString("<value>")
		switch val := v.(type) {
		case string:
			buf.WriteString(fmt.Sprintf("<string>%s</string>", xmlEscape(val)))
		case int:
			buf.WriteString(fmt.Sprintf("<int>%d</int>", val))
		default:
			buf.WriteString(fmt.Sprintf("<string>%v</string>", val))
		}
		buf.WriteString("</value>")
		buf.WriteString("</member>")
	}
	buf.WriteString("</struct>")
	return buf.String()
}

// xmlEscape escapes a string for XML.
func xmlEscape(s string) string {
	return html.EscapeString(s)
}

// AddZoneRecord adds a DNS record to a zone using the Loopia API.
func (c *LoopiaClient) AddZoneRecord(domain, subdomain string, recordObj map[string]interface{}) error {
	reqBody := buildXMLRPCRequest("addZoneRecord", c.Username, c.Password, domain, subdomain, recordObj)
	resp, err := c.HTTPClient.Post(c.Endpoint, "text/xml", bytes.NewBufferString(reqBody))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	respData, _ := io.ReadAll(resp.Body)
	// TODO: Parse XML-RPC response and handle errors
	fmt.Println(string(respData))
	return nil
}

// RemoveZoneRecord removes a DNS record from a zone using the Loopia API.
func (c *LoopiaClient) RemoveZoneRecord(domain, subdomain string, recordID int) error {
	reqBody := buildXMLRPCRequest("removeZoneRecord", c.Username, c.Password, domain, subdomain, recordID)
	resp, err := c.HTTPClient.Post(c.Endpoint, "text/xml", bytes.NewBufferString(reqBody))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	respData, _ := io.ReadAll(resp.Body)
	// TODO: Parse XML-RPC response and handle errors
	fmt.Println(string(respData))
	return nil
}

// GetZoneRecords fetches all DNS records for a given zone and subdomain using the Loopia API.
func (c *LoopiaClient) GetZoneRecords(domain, subdomain string) ([]map[string]interface{}, error) {
	reqBody := buildXMLRPCRequest("getZoneRecords", c.Username, c.Password, domain, subdomain)
	resp, err := c.HTTPClient.Post(c.Endpoint, "text/xml", bytes.NewBufferString(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	respData, _ := io.ReadAll(resp.Body)
	// TODO: Parse XML-RPC response and return records
	fmt.Println(string(respData))
	return nil, nil
}

// RealClientFactory creates a real LoopiaClient from provider config.
func RealClientFactory(ctx context.Context, config Config) (Client, error) {
	return NewLoopiaClient(config.Username, config.Password, config.Endpoint), nil
}
