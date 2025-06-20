package monitor

type Alert struct {
	Alert       string
	Expr        string
	For         string
	Labels      map[string]string
	Annotations map[string]string
}

type Rule struct {
	Name  string
	Rules []Alert
}

type Group struct {
	Groups []Rule
}

type Global struct {
	ResolveTimeout   string `yaml:"resolve_timeout,omitempty" json:"resolve_timeout"`
	SmtpFrom         string `yaml:"smtp_from,omitempty" json:"smtp_from"`
	SmtpSmarthost    string `yaml:"smtp_smarthost,omitempty" json:"smtp_smarthost"`
	SmtpAuthUsername string `yaml:"smtp_auth_username,omitempty" json:"smtp_auth_username"`
	SmtpAuthPassword string `yaml:"smtp_auth_password,omitempty" json:"smtp_auth_password"`
	WechatApiSecret  string `yaml:"wechat_api_secret,omitempty" json:"wechat_api_secret"`
	WechatApiCorpId  string `yaml:"wechat_api_corp_id,omitempty" json:"wechat_api_corp_id"`
}

type Route struct {
	GroupBy             []string            `yaml:"group_by,omitempty,flow" json:"group_by,omitempty,flow"`
	GroupWait           string              `yaml:"group_wait,omitempty" json:"group_wait,omitempty"`
	GroupInterval       string              `yaml:"group_interval,omitempty" json:"group_interval,omitempty"`
	RepeatInterval      string              `yaml:"repeat_interval,omitempty" json:"repeat_interval,omitempty"`
	Receiver            string              `yaml:"receiver,omitempty" json:"receiver,omitempty"`
	Continue            bool                `yaml:"continue,omitempty" json:"continue,omitempty"`
	Routes              []Route             `yaml:"routes,omitempty" json:"routes,omitempty"`
	Match               map[string]string   `yaml:"match,omitempty" json:"match,omitempty"`
	MatchRe             map[string]string   `yaml:"match_re,omitempty" json:"match_re,omitempty"`
	Matchers            []map[string]string `yaml:"matchers,omitempty" json:"matchers,omitempty"`
	MuteTimeIntervals   []string            `yaml:"mute_time_intervals,omitempty" json:"mute_time_intervals,omitempty"`
	ActiveTimeIntervals []string            `yaml:"active_time_intervals,omitempty" json:"active_time_intervals,omitempty"`
}

type NotifierConfig struct {
	SendResolved bool `yaml:"send_resolved" json:"send_resolved"`
}

type EmailConfig struct {
	NotifierConfig `yaml:",inline" json:",inline"`
	To             string            `yaml:"to,omitempty" json:"to,omitempty"`
	From           string            `yaml:"from,omitempty" json:"from,omitempty"`
	Headers        map[string]string `yaml:"headers,omitempty,flow" json:"headers,omitempty,flow"`
	HTML           string            `yaml:"html,omitempty" json:"html,omitempty"`
}

type WebhookConfig struct {
	NotifierConfig `yaml:",inline" json:",inline"`
	URL            string `yaml:"url,omitempty" json:"url,omitempty"`
}

type WechatConfig struct {
	NotifierConfig `yaml:",inline" json:",inline"`
	ToUser         string `yaml:"to_user,omitempty" json:"to_user,omitempty"`
	ToParty        string `yaml:"to_party,omitempty" json:"to_party,omitempty"`
	ToTag          string `yaml:"to_tag,omitempty" json:"to_tag,omitempty"`
	AgentID        string `yaml:"agent_id,omitempty" json:"agent_id,omitempty"`
	ApiSecret      string `yaml:"api_secret,omitempty" json:"api_secret,omitempty"`
	CorpId         string `yaml:"corp_id,omitempty" json:"corp_id,omitempty"`
	Message        string `yaml:"message,omitempty" json:"message,omitempty"`
}

type InhibitRule struct {
	SourceMatch map[string]string `yaml:"source_match,omitempty" json:"source_match,omitempty"`
	TargetMatch map[string]string `yaml:"target_match,omitempty" json:"target_match,omitempty"`
	Equal       []string          `yaml:"equal,omitempty,flow" json:"equal,omitempty,flow"`
}

type Receiver struct {
	Name           string          `yaml:"name" json:"name"`
	WebhookConfigs []WebhookConfig `yaml:"webhook_configs,omitempty" json:"webhook_configs,omitempty"`
	EmailConfigs   []EmailConfig   `yaml:"email_configs,omitempty" json:"email_configs,omitempty"`
	WechatConfigs  []WechatConfig  `yaml:"wechat_configs,omitempty" json:"wechat_configs,omitempty"`
}

type AlertManagerConfig struct {
	Global            map[string]interface{} `yaml:"global" json:"global"`
	Route             Route                  `yaml:"route" json:"route"`
	Templates         []string               `yaml:"templates,omitempty" json:"templates,omitempty"`
	Receivers         []Receiver             `yaml:"receivers,omitempty" json:"receivers,omitempty"`
	InhibitRules      []InhibitRule          `yaml:"inhibit_rules,omitempty" json:"inhibit_rules,omitempty"`
	MuteTimeIntervals []interface{}          `yaml:"mute_time_intervals,omitempty" json:"mute_time_intervals,omitempty"`
	TimeIntervals     []interface{}          `yaml:"time_intervals,omitempty" json:"time_intervals,omitempty"`
}
