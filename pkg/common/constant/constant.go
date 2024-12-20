package constant

type Hashcode struct{}
type StatusInvitation string
type MetricAPI string
type Context string

const (
	CTX_ID_USER    Context   = "x-id-user"
	CTX_HASHCODE   Context   = "x-hashcode"
	CTX_AUTH_TOKEN Context   = "x-auth-token"
	METRIC_ID_USER MetricAPI = "id-user"

	AUDIO_FOLDER = "audio"
)
