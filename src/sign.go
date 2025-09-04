package sdk

/*
#cgo LDFLAGS: -L${SRCDIR} -lorderffi
#cgo LDFLAGS: -Wl,-rpath,$ORIGIN
#include <stdlib.h>

char* get_order_hash_ffi(
    const char* position_id,
    const char* base_asset_id_hex,
    const char* base_amount,
    const char* quote_asset_id_hex,
    const char* quote_amount,
    const char* fee_asset_id_hex,
    const char* fee_amount,
    const char* expiration,
    const char* salt,
    const char* user_public_key_hex,
    const char* domain_name,
    const char* domain_version,
    const char* domain_chain_id,
    const char* domain_revision
);

char* sign_message_ffi(const char* message_hex, const char* private_key_hex);

void free_string(char* s);
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// GetOrderHash computes the order hash using the provided parameters.
func GetOrderHash(
	positionID, baseAssetIDHex, baseAmount,
	quoteAssetIDHex, quoteAmount,
	feeAssetIDHex, feeAmount,
	expiration, salt,
	userPublicKeyHex,
	domainName, domainVersion,
	domainChainID, domainRevision string,
) (string, error) {
	// allocate C strings
	toC := func(s string) *C.char {
		return C.CString(s)
	}
	args := []*C.char{
		toC(positionID), toC(baseAssetIDHex), toC(baseAmount),
		toC(quoteAssetIDHex), toC(quoteAmount),
		toC(feeAssetIDHex), toC(feeAmount),
		toC(expiration), toC(salt),
		toC(userPublicKeyHex),
		toC(domainName), toC(domainVersion),
		toC(domainChainID), toC(domainRevision),
	}
	// ensure we free them all
	for _, cstr := range args {
		defer C.free(unsafe.Pointer(cstr))
	}

	// call into Rust
	result := C.get_order_hash_ffi(
		args[0], args[1], args[2], args[3],
		args[4], args[5], args[6], args[7],
		args[8], args[9], args[10], args[11],
		args[12], args[13],
	)
	if result == nil {
		return "", fmt.Errorf("get_order_hash_ffi returned NULL")
	}
	defer C.free_string(result)

	return C.GoString(result), nil
}

// SignMessage signs a message using the provided private key.
// It returns the signature as a hex string, where v, r and s are concatenated as left-padded 64-character hex strings.
// The signature is in the format: {r}{s}{v}
func SignMessage(messageHex, privateKeyHex string) (string, error) {
	cmsg := C.CString(messageHex)
	defer C.free(unsafe.Pointer(cmsg))
	cpriv := C.CString(privateKeyHex)
	defer C.free(unsafe.Pointer(cpriv))

	sig := C.sign_message_ffi(cmsg, cpriv)
	if sig == nil {
		return "", fmt.Errorf("sign_message_ffi returned NULL")
	}
	defer C.free_string(sig)

	return C.GoString(sig), nil
}
