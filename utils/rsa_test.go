package utils

import (
	"testing"
)

var priKey = "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlDWFFJQkFBS0JnUUM5NzFWbTY5MnBkcjRyeUlIc0ZkY0pGQW9lTGtZNnE2VUlCNG1XQ0M4V0ExY1VwYjBhCkNwcHJSckQwdmFVamhNZW5uWk4za0I5U1lGU2doNHoxVWtYNUcrd2VhTGRCbVI3bmFEMURKdjBEbkRQTWtlMTUKUlE2WU1hanFCZjFkOElROEYwZW1uTU1tVWZuTFFHMlh4VmV3STJGMWVnak5qNWJLR1N6OVBSVEQ5UUlEQVFBQgpBb0dCQUljVStuWXlkZm1hVy9JanJsTkx6UjNGeE5SbU1NaDFYdS93L0dkWjlyTC9PU1dVSW9PczJ0cEU4b0Y5Cmh6OVZwZkdOM2wyQWdPWkRZS3l5K3d0V3NqQVd0R2RDRlhsQTZxT0JJaWJ1NHZiL0lBUzlBdWlpMDZsakFzaVoKNjZSTUxxSFNUbXZaaGVoejVISksrSTRNQ0wzQXA2WmV2NnlVc2NqMjd4dnlyb0JaQWtFQS9COG1ManJWSmd5RQpoMzBlWmdKUjdpVjQ3ZHl3anhmTlo2Vkcwb09LUVkrcEJMbDlZWXlqYUdBY1BacFF5UDdzdEcvVTh3QlhtREtQCkJJNjBGRTJlWXdKQkFNRGJTMEozeXF5NGNNN1FkM0VDdngwZzI5aXBuZldrSHA4RTFrUHcwdHhMbmlhbEpyengKd2dLVjhxTHZRc2xMMHhRNVZIcE5MSmpheDc5ZGVJYlpWOGNDUUIydUFRMmlLV053UjgyM1lmTzZSRERYd25PbAo3amI2STFrWE1NNHBaQVl4eGtEaklTcHhwdTdybVlkNitoV2ZSUGc4emdISlFZYU9OUjNoT3J2RkkyY0NRRXY1Cm1DcDFPcmpVYUV4eFA3eWJrbUtOUVU2WGM0MER2TFIwbVZ4bWtRc01GeCt1VEJaL1B5ajVuWDZtdHk3SjJqdkwKWWdaVVJNOXEwT29JanFUQkZwMENRUUQzUVVxR1RGZ28zMUl2UXp0NENXL25KTjN6dFZWa2tteWxPYzlnbnp6VwplMXRJNElUcXU5Nm82cE9tOWJtYzNvWG1hSFhUSGViVm9Ha1BpTVVoR2RaTwotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo="
var pubKey = "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlHZk1BMEdDU3FHU0liM0RRRUJBUVVBQTRHTkFEQ0JpUUtCZ1FDOTcxVm02OTJwZHI0cnlJSHNGZGNKRkFvZQpMa1k2cTZVSUI0bVdDQzhXQTFjVXBiMGFDcHByUnJEMHZhVWpoTWVublpOM2tCOVNZRlNnaDR6MVVrWDVHK3dlCmFMZEJtUjduYUQxREp2MERuRFBNa2UxNVJRNllNYWpxQmYxZDhJUThGMGVtbk1NbVVmbkxRRzJYeFZld0kyRjEKZWdqTmo1YktHU3o5UFJURDlRSURBUUFCCi0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQo="

func TestRsaNewKey(t *testing.T) {
	pri, pub, err := RsaNewKey()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(Base64ToString(pri))
	t.Log(Base64ToString(pub))
}

func TestRsaCrypt(t *testing.T) {
	signText := []byte("1111")
	pri, _ := Base64ToBytes(priKey)
	pub, _ := Base64ToBytes(pubKey)
	res, err := RsaEncrypt(signText, pub)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(Base64ToString(res))
	res, err = RsaDecrypt(res, pri)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(res))
}

func TestRsaSign(t *testing.T) {
	signText := []byte("1111")
	pri, _ := Base64ToBytes(priKey)
	pub, _ := Base64ToBytes(pubKey)
	sign, err := RsaSign(pri, signText)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(Base64ToString(sign))
	t.Log(RsaVerifySign(pub, signText, sign))
}
