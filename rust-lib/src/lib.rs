use std::ffi::{CStr, CString};
use std::os::raw::c_char;

// bring your original functions into scope
use rust_crypto_lib_base::{get_order_hash, sign_message};
use starknet_crypto::Felt;

#[no_mangle]
pub extern "C" fn get_order_hash_ffi(
    pos_id: *const c_char,
    base_asset_id_hex: *const c_char,
    base_amount: *const c_char,
    quote_asset_id_hex: *const c_char,
    quote_amount: *const c_char,
    fee_asset_id_hex: *const c_char,
    fee_amount: *const c_char,
    expiration: *const c_char,
    salt: *const c_char,
    user_pubkey_hex: *const c_char,
    domain_name: *const c_char,
    domain_version: *const c_char,
    domain_chain_id: *const c_char,
    domain_revision: *const c_char,
) -> *mut c_char {
    // convert C strings → Rust Strings
    macro_rules! to_str { ($p:expr) => {
        unsafe { CStr::from_ptr($p).to_str().unwrap().to_string() }
    }}
    let result = get_order_hash(
        to_str!(pos_id),
        to_str!(base_asset_id_hex),
        to_str!(base_amount),
        to_str!(quote_asset_id_hex),
        to_str!(quote_amount),
        to_str!(fee_asset_id_hex),
        to_str!(fee_amount),
        to_str!(expiration),
        to_str!(salt),
        to_str!(user_pubkey_hex),
        to_str!(domain_name),
        to_str!(domain_version),
        to_str!(domain_chain_id),
        to_str!(domain_revision),
    )
    .map(|felt| felt.to_hex_string())            // or however you serialize `Felt`
    .unwrap_or_else(|e| e);

    CString::new(result).unwrap().into_raw()
}

#[no_mangle]
pub extern "C" fn sign_message_ffi(
    msg_hex: *const c_char,
    priv_hex: *const c_char,
) -> *mut c_char {
    let message = unsafe {
        let s = CStr::from_ptr(msg_hex).to_str().unwrap();
        Felt::from_hex_unchecked(s)
    };
    let privkey = unsafe {
        let s = CStr::from_ptr(priv_hex).to_str().unwrap();
        Felt::from_hex_unchecked(s)
    };

    let sig = sign_message(&message, &privkey)
        .map(|stark_sig: rust_crypto_lib_base::StarkSignature| stark_sig.to_hex_string())
        .unwrap_or_else(|e| e);

    CString::new(sig).unwrap().into_raw()
}

#[no_mangle]
pub extern "C" fn free_string(s: *mut c_char) {
    if s.is_null() { return; }
    unsafe { let _ = CString::from_raw(s); } // drops it
}


trait SerializableSig {
    fn to_hex_string(&self) -> String;
    #[allow(dead_code)]
    fn from_hex_string(hex: &str) -> Self;
}

impl SerializableSig for rust_crypto_lib_base::StarkSignature {
    fn to_hex_string(&self) -> String {
        // Serialize r and s as hex strings and knock off the 0x prefix
        let r_str = self.r.to_fixed_hex_string();
        let s_str = self.s.to_fixed_hex_string();
        let v_str = self.v.to_fixed_hex_string();

        format!(
            "{}{}{}",
            &r_str[2..], // skip "0x"
            &s_str[2..],
            &v_str[2..],
        )
    }

    fn from_hex_string(hex: &str) -> Self {
        // Deserialize from hex string
        let r = Felt::from_hex_unchecked(&format!("0x{}", &hex[0..64]));
        let s = Felt::from_hex_unchecked(&format!("0x{}", &hex[64..128]));
        let v = Felt::from_hex_unchecked(&format!("0x{}", &hex[128..192]));

        rust_crypto_lib_base::StarkSignature { r, s, v }
    }
}
