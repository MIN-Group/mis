/**
 * @Author: wzx
 * @Description: 兼容go crypto实现的基础定义
 * @Version: 1.0.0
 * @Date: 2021/1/15 下午11:04
 *@Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
 */

package minsecurity

import (
	"hash"
	"io"
	"strconv"
)

// Hash identities a cryptographic hash function that is implemented in another
// package.
type Hash uint

// HashFunc simply returns the value of h so that Hash implements SignerOpts.
func (h Hash) HashFunc() Hash {
	return h
}

const (
	SM3 Hash = 1 + iota
	MD5
	SHA256
	SHA512
	maxHash
)

var digestSizes = []uint8{
	SM3:    16,
	MD5:    16,
	SHA256: 32,
	SHA512: 64,
}

// Size returns the length, in bytes, of a digest resulting from the given hash
// function. It doesn't require that the hash function in question be linked
// into the program.
func (h Hash) Size() int {
	if h > 0 && h < maxHash {
		return int(digestSizes[h])
	}
	panic("crypto: Size of unknown hash function")
}

var hashes = make([]func() hash.Hash, maxHash)

// New returns a new hash.Hash calculating the given hash function. New panics
// if the hash function is not linked into the binary.
func (h Hash) New() hash.Hash {
	if h > 0 && h < maxHash {
		f := hashes[h]
		if f != nil {
			return f()
		}
	}
	panic("crypto: requested hash function #" + strconv.Itoa(int(h)) + " is unavailable")
}

// RegisterHash registers a function that returns a new instance of the given
// hash function. This is intended to be called from the init function in
// packages that implement hash functions.
func RegisterHash(h Hash, f func() hash.Hash) {
	if h >= maxHash {
		panic("crypto: RegisterHash of unknown hash function")
	}
	hashes[h] = f
}

// Available reports whether the given hash function is linked into the binary.
func (h Hash) Available() bool {
	return h < maxHash && hashes[h] != nil
}

// PublicKey represents a public key using an unspecified algorithm.
//函数的参数是兼容go crypto的实现，在国密中*Opts都是nil处理
type PublicKey interface {
	GetBytes() []byte
	SetBytes([]byte) error
	Encrypt(rand io.Reader, msg []byte, opts DecrypterOpts) (encryptedtext []byte, err error)
	Verify(msg []byte, digest []byte, opts SignerOpts) (bool, error)
}

// PrivateKey represents a private key using an unspecified algorithm.
type PrivateKey interface {
	GetBytes() []byte
	SetBytes([]byte) error
	Sign(rand io.Reader, digest []byte, opts SignerOpts) (signature []byte, err error)
	Decrypt(rand io.Reader, msg []byte, opts DecrypterOpts) (plaintext []byte, err error)
}

// SignerOpts contains options for signing with a Signer.
type SignerOpts interface {
	// HashFunc returns an identifier for the hash function used to produce
	// the message passed to Signer.Sign, or else zero to indicate that no
	// hashing was done.
	HashFunc() Hash
}

type DecrypterOpts interface{}
