package contexthelper

type requestIdKeyType struct{}
type userIdKeyType struct{}
type configKeyType struct{}
type clientIpKeyType struct{}
type dbKeyType struct{}
type rabbitKeyType struct{}
type accessTokenDataCtxKeyType struct{}

var (
	requestIdCtxKey       = requestIdKeyType{}
	userIdCtxKey          = userIdKeyType{}
	configCtxKey          = configKeyType{}
	clientIpCtxKey        = clientIpKeyType{}
	dbCtxKey              = dbKeyType{}
	rabbitCtxKey          = rabbitKeyType{}
	accessTokenDataCtxKey = accessTokenDataCtxKeyType{}
)
