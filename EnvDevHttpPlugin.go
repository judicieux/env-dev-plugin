package envdev_plugin

import (
        "github.com/LeakIX/l9format"
        "regexp"
        "strings"
)

type EnvDevHttpPlugin.go.go struct {
        l9format.ServicePluginBase
}

func (EnvDevHttpPlugin.go) GetVersion() (int, int, int) {
        return 0, 0, 2
}

func (EnvDevHttpPlugin.go) GetRequests() []l9format.WebPluginRequest {
        return []l9format.WebPluginRequest{{
                Method:  "GET",
                Path:    "/.env.example",
                Headers: map[string]string{},
                Body:    []byte(""),
        }}
}

func (EnvDevHttpPlugin.go) GetName() string {
        return "EnvDevHttpPlugin.go"
}

func (EnvDevHttpPlugin.go) GetStage() string {
        return "open"
}

func (plugin EnvDevHttpPlugin.go) Verify(request l9format.WebPluginRequest, response l9format.WebPluginResponse, event *l9format.L9Event, options map[string]string) (hasLeak bool) {
        if !request.EqualAny(plugin.GetRequests()) || response.Response.StatusCode != 200 {
                return false
        }
        lowerBody := strings.ToLower(string(response.Body))
        if len(lowerBody) < 10 {
                return false
        }

        regexPattern := `(?i)(app_env=|db_host=|\bAKIA[A-Z0-9]{16}\b|eu-smtp-outbound-1\.mimecast\.com|eu-smtp-outbound-2\.mimecast\.com|de-smtp-outbound-1\.mimecast\.com|de-smtp-outbound-2\.mimecast\.com|us-smtp-outbound-1\.mimecast\.com|us-smtp-outbound-2\.mimecast\.com|ca-smtp-outbound-1\.mimecast\.com|ca-smtp-outbound-2\.mimecast\.com|za-smtp-outbound-1\.mimecast\.co\.za|za-smtp-outbound-2\.mimecast\.co\.za|au-smtp-outbound-1\.mimecast\.com|au-smtp-outbound-2\.mimecast\.com|je-smtp-outbound-1\.mimecast-offshore\.com|je-smtp-outbound-2\.mimecast-offshore\.com|usb-smtp-outbound-1\.mimecast\.com|usb-smtp-outbound-2\.mimecast\.com|uspcom-smtp-outbound-1\.mimecast-pscom-us\.com|uspcom-smtp-outbound-2\.mimecast-pscom-us\.com|SG\.[a-zA-Z0-9_-]{22}\.[a-zA-Z0-9_-]{43}|smtp-relay\.sendinblue\.com|mail\.smtp2go\.com|smtp-relay\.brevo\.com|smtp\.mailgun\.org|pro\.turbo-smtp\.com|in-v3\.mailjet\.com|smtp\.postmarkapp\.com|mail\.smtpeter\.com|mail\.infomaniak\.com|smtp\.pepipost\.com|smtp-pulse\.com|smtp\.mandrillapp\.com|smtp\.sparkpostmail\.com|dapi[a-h0-9]{32})`

        match, _ := regexp.MatchString(regexPattern, lowerBody)
        if match {
                event.Service.Software.Name = "EnvironmentFile"
                event.Leak.Type = "config_leak"
                event.Leak.Severity = "high"
                event.AddTag("potential-leak")
                event.Summary = "Found sensitive information in /.env.example:\n" + string(response.Body)
                return true
        }

        return false
}