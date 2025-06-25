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

type Client interface {
	AddZoneRecord(domain, subdomain string, recordObj map[string]interface{}) error
	RemoveZoneRecord(domain, subdomain string, recordID int) error
	GetZoneRecords(domain, subdomain string) ([]map[string]interface{}, error)
}

type LoopiaClient struct {
	Username   string
	Password   string
	Endpoint   string
	HTTPClient *http.Client
}

func NewLoopiaClient(username, password, endpoint string) *LoopiaClient {
	return &LoopiaClient{
		Username:   username,
		Password:   password,
		Endpoint:   endpoint,
		HTTPClient: &http.Client{},
	}
}

// XML-RPC request/response structures and helpers
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

func xmlEscape(s string) string {
	return html.EscapeString(s)
}

// LoopiaClient implements Client interface
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

// Factory function for real client
func RealClientFactory(ctx context.Context, config Config) (Client, error) {
	return NewLoopiaClient(config.Username, config.Password, config.Endpoint), nil
}
