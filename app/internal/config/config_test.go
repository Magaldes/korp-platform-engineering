package config

import "testing"

func TestCacheFromEnvDefaults(t *testing.T) {
	t.Setenv("KORP_REDIS_ADDR", "")
	t.Setenv("KORP_CACHE_TTL_SECONDS", "")
	cache, err := CacheFromEnv()
	if err != nil || cache.Address != "redis:6379" || cache.TTL.Seconds() != 60 {
		t.Fatalf("cache=%+v err=%v", cache, err)
	}
}

func TestCacheFromEnvRejectsInvalidTTL(t *testing.T) {
	for _, raw := range []string{"0", "-1", "not-a-number"} {
		t.Setenv("KORP_CACHE_TTL_SECONDS", raw)
		if _, err := CacheFromEnv(); err == nil {
			t.Fatalf("TTL %q was accepted", raw)
		}
	}
}
