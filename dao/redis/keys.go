package redis

const (
	KeyPrefix = "yun:"
	KeyEmailSet = "email:"
	KeyInviteCode = "invite:"
)

func getRedisKey(key string) string {
	return KeyPrefix+key
}
