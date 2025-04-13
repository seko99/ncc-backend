// MPPE to support encryption with MSCHAPv2
package radius

import (
	"crypto/md5"
	"crypto/sha1"
	"math/rand"
	"time"
)

// Pads used in key derivation
var SHSpad1 = []byte{
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
}
var SHSpad2 = []byte{
	0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2,
	0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2,
	0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2,
	0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2, 0xf2,
}

// "Magic" constants used in key derivations
var magic1 = []byte{
	0x54, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, 0x20, 0x74,
	0x68, 0x65, 0x20, 0x4d, 0x50, 0x50, 0x45, 0x20, 0x4d,
	0x61, 0x73, 0x74, 0x65, 0x72, 0x20, 0x4b, 0x65, 0x79,
}
var magic2 = []byte{
	0x4f, 0x6e, 0x20, 0x74, 0x68, 0x65, 0x20, 0x63, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x20, 0x73, 0x69, 0x64, 0x65, 0x2c, 0x20,
	0x74, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, 0x20, 0x74, 0x68,
	0x65, 0x20, 0x73, 0x65, 0x6e, 0x64, 0x20, 0x6b, 0x65, 0x79,
	0x3b, 0x20, 0x6f, 0x6e, 0x20, 0x74, 0x68, 0x65, 0x20, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x20, 0x73, 0x69, 0x64, 0x65,
	0x2c, 0x20, 0x69, 0x74, 0x20, 0x69, 0x73, 0x20, 0x74, 0x68,
	0x65, 0x20, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x20,
	0x6b, 0x65, 0x79, 0x2e,
}
var magic3 = []byte{
	0x4f, 0x6e, 0x20, 0x74, 0x68, 0x65, 0x20, 0x63, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x20, 0x73, 0x69, 0x64, 0x65, 0x2c, 0x20,
	0x74, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, 0x20, 0x74, 0x68,
	0x65, 0x20, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x20,
	0x6b, 0x65, 0x79, 0x3b, 0x20, 0x6f, 0x6e, 0x20, 0x74, 0x68,
	0x65, 0x20, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x20, 0x73,
	0x69, 0x64, 0x65, 0x2c, 0x20, 0x69, 0x74, 0x20, 0x69, 0x73,
	0x20, 0x74, 0x68, 0x65, 0x20, 0x73, 0x65, 0x6e, 0x64, 0x20,
	0x6b, 0x65, 0x79, 0x2e,
}

func getMasterKey(hashHash []byte, ntRes []byte) []byte {
	magic := []byte{
		0x54, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, 0x20, 0x74,
		0x68, 0x65, 0x20, 0x4d, 0x50, 0x50, 0x45, 0x20, 0x4d,
		0x61, 0x73, 0x74, 0x65, 0x72, 0x20, 0x4b, 0x65, 0x79,
	}

	hash := sha1.New()
	hash.Write(hashHash)
	hash.Write(ntRes)
	hash.Write(magic)
	return hash.Sum(nil)[:16]
}

func getAsymmetricStartKey(masterKey []byte, sessKeyLen int, isSend bool) []byte {
	var magic []byte

	if isSend {
		magic = magic3
	} else {
		magic = magic2
	}
	if len(SHSpad1) != 40 {
		panic("Should not get here")
	}
	if len(SHSpad2) != 40 {
		panic("Should not get here")
	}
	if len(magic) != 84 {
		panic("Should not get here")
	}

	hash := sha1.New()
	hash.Write(masterKey)
	hash.Write(SHSpad1)
	hash.Write(magic)
	hash.Write(SHSpad2)
	return hash.Sum(nil)[:sessKeyLen]
}

func masterKeys(pass string, ntResponse []byte) ([]byte, []byte) {
	// PasswordHashHash( NtPasswordHash(Password, PasswordHash) )
	hashHash := hashNtPasswordHash(ntPasswordHash(ntPassword(pass)))
	// GetMasterKey(PasswordHashHash, NtResponse, MasterKey)
	masterKey := getMasterKey(hashHash, ntResponse)

	sendKey := getAsymmetricStartKey(masterKey, 16, true)
	recvKey := getAsymmetricStartKey(masterKey, 16, false)
	return sendKey, recvKey
}

func multipleOf(val []byte, size int) []byte {
	if len(val)%size != 0 {
		val = append(val, make([]byte, size-(len(val)%size))...)
	}
	return val
}

func xor(a []byte, b []byte) []byte {
	if len(a) != len(b) {
		panic("Should not get here")
	}
	n := len(a)
	out := make([]byte, n)
	for i := 0; i < n; i++ {
		out[i] = a[i] ^ b[i]
	}
	return out
}

/*
   Construct a plaintext version of the String field by concate-
   nating the Key-Length and Key sub-fields.  If necessary, pad
   the resulting string until its length (in octets) is an even
   multiple of 16.  It is recommended that zero octets (0x00) be
   used for padding.  Call this plaintext P.

   Call the shared secret S, the pseudo-random 128-bit Request
   Authenticator (from the corresponding Access-Request packet) R,
   and the contents of the Salt field A.  Break P into 16 octet
   chunks p(1), p(2)...p(i), where i = len(P)/16.  Call the
   ciphertext blocks c(1), c(2)...c(i) and the final ciphertext C.
   Intermediate values b(1), b(2)...c(i) are required.  Encryption
   is performed in the following manner ('+' indicates
   concatenation)

   http://security.stackexchange.com/questions/35683/mppe-send-and-receive-key-derivation-from-ms-chapv2
   https://github.com/FreeRADIUS/freeradius-server/blob/5ea87f156381174ea24340db9b450d4eca8189c9/src/lib/radius.c#L623
*/
func tunnelPass(secret string, key []byte, reqAuth []byte, salt []byte) []byte {
	// concatenating the Key-Length and Key sub-fields.
	P := append([]byte{byte(len(key))}, key...)
	// If necessary, pad the resulting string until its length (in octets) is an even
	// multiple of 16.
	P = multipleOf(P, 16)

	var b [][]byte
	var c [][]byte
	var C []byte

	// Break P into 16 octet chunks where i = len(P)/16
	for i := 0; i < len(P)/16; i++ {
		p := P[i*16 : (i+1)*16]
		if i == 0 {
			// b(1) = MD5(S + R + A)
			{
				hash := md5.New()
				hash.Write([]byte(secret))
				hash.Write(reqAuth)
				hash.Write(salt)
				b = append(b, hash.Sum(nil))
			}
			// c(1) = p(1) xor b(1)
			{
				c = append(c, xor(p, b[i]))
			}
			// C = c(1)
			C = c[i]
		} else {
			// b(i) = MD5(S + c(i-1))
			{
				hash := md5.New()
				hash.Write([]byte(secret))
				hash.Write(c[i-1])
				b = append(b, hash.Sum(nil))
			}
			// c(i) = p(i) xor b(i)
			{
				c = append(c, xor(p, b[i]))
			}
			// C = C + c(i)
			C = append(C, c[i]...)
		}
	}

	if len(C)%16 != 0 {
		panic("Should not get here")
	}

	/*
	   The plaintext String field consists of three logical sub-fields:
	   - the Key-Length
	   - and Key sub-fields (both of which are required),
	   - and the optional Padding sub-field.

	   * The Key-Length sub-field is one octet in length and contains the length of the unencrypted Key
	    sub-field.
	   * The Key sub-field contains the actual encryption key.

	   * If the combined length (in octets) of the unencrypted Key-Length
	    and Key sub-fields is not an even multiple of 16, then the Padding
	    sub-field MUST be present.  If it is present, the length of the
	    Padding sub-field is variable, between 1 and 15 octets.
	*/
	plain := make([]byte, 2+len(C))
	plain[0] = salt[0]
	plain[1] = salt[1]
	copy(plain[2:], C)
	return plain
}

func salt(offset uint8) []byte {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	salt := uint8(r.Uint32())

	pfx := make([]byte, 2)
	pfx[0] = byte(0x80 | ((offset & 0x0f) << 3) | (salt & 0x07))
	pfx[1] = byte(salt)
	return pfx
}

func Mmpev2(secret string, pass string, reqAuth []byte, ntResponse []byte) ([]byte, []byte) {
	send, recv := masterKeys(pass, ntResponse)
	sendEnc := tunnelPass(secret, send, reqAuth, salt(0))
	recvEnc := tunnelPass(secret, recv, reqAuth, salt(1))

	return sendEnc, recvEnc
}
