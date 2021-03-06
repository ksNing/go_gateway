package public

const (
	ValidatorKey        = "ValidatorKey"
	TranslatorKey       = "TranslatorKey"
	AdminSessionInfoKey = "AdminSessionInfoKey"

	LoadTypeHTTP = 0
	LoadTypeTCP  = 1
	LoadTypeGRPC = 2

	HTTPRuleTypePrefixURL = 0
	HTTPRuleTypeDomain    = 1

	RedisFlowDayKey  = "redis_day_key"
	RedisFlowHourKey = "redis_hour_key"

	Flowtotal              = "flow_total"
	FlowServiceCountPrefix = "flow_service_"
	FlowAppCount           = "flow_app_"
)
