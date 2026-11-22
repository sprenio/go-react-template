package contexthelper

type requestIdKeyType struct{}
type userIdKeyType struct{}
type configKeyType struct{}
type clientIpKeyType struct{}
type dbKeyType struct{}
type rabbitKeyType struct{}
type redisClientKeyType struct{}
type accessTokenDataCtxKeyType struct{}
type queryDateCtxKeyType struct{}
type queryYearMonthCtxKeyType struct{}
type currencyOrderKeyType struct{}
type rateLimiterKeyType struct{}

var (
	requestIdCtxKey       = requestIdKeyType{}
	userIdCtxKey          = userIdKeyType{}
	configCtxKey          = configKeyType{}
	clientIpCtxKey        = clientIpKeyType{}
	dbCtxKey              = dbKeyType{}
	rabbitCtxKey          = rabbitKeyType{}
	redisClientCtxKey     = redisClientKeyType{}
	accessTokenDataCtxKey = accessTokenDataCtxKeyType{}
	queryDateCtxKey       = queryDateCtxKeyType{}
	queryYearMonthCtxKey  = queryYearMonthCtxKeyType{}
	currencyOrderCtxKey   = currencyOrderKeyType{}
	rateLimiterCtxKey     = rateLimiterKeyType{}
)
