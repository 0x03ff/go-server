// ecdh-aes.js
export async function generateECDHKeyPairP384(hash) {
    try {
        // Generate ECDH key pair for P-384
        const keyPair = await crypto.subtle.generateKey(
            {
                name: "ECDH",
                namedCurve: "P-384"
            },
            true,
            ["deriveKey", "deriveBits"]
        );

        // Export private key in PKCS#8 DER format
        const privKeyDer = await crypto.subtle.exportKey("pkcs8", keyPair.privateKey);
        const privKeyPem = convertDerToPem(privKeyDer, "PRIVATE KEY");

        // Export public key in SPKI (PKIX) DER format
        const pubKeyDer = await crypto.subtle.exportKey("spki", keyPair.publicKey);
        const pubKeyPem = convertDerToPem(pubKeyDer, "PUBLIC KEY");

        // Convert PEM string to ArrayBuffer before encryption
        const privKeyBuffer = new TextEncoder().encode(privKeyPem).buffer;

        // Encrypt the private key (correct parameter order: key, data)
        const encryptedPrivateKey = await encryptWithAESCTR(
            await convertHashToCryptoKey(hash),
            privKeyBuffer
        );

        // Convert the encrypted data to base64 for JSON serialization
        const encryptedPrivateKeyBase64 = arrayBufferToBase64(encryptedPrivateKey);

        return {
            publicKey: pubKeyPem,
            privateKey: encryptedPrivateKeyBase64
        };
    } catch (err) {
        throw new Error(`Key generation failed: ${err.message}`);
    }
}



export async function ECDH_hash256(message) {
    const msgUint8 = new TextEncoder().encode(message);
    const hashBuffer = await crypto.subtle.digest('SHA-256', msgUint8);
    return new Uint8Array(hashBuffer);
}

async function convertHashToCryptoKey(hash) {
    return crypto.subtle.importKey(
        "raw",
        hash,
        { name: "AES-CTR" },
        false,
        ["encrypt","decrypt"]
    );
}

export function arrayBufferToBase64(buffer) {
    const bytes = new Uint8Array(buffer);
    let binary = '';
    for (let i = 0; i < bytes.length; i++) {
        binary += String.fromCharCode(bytes[i]);
    }
    return btoa(binary);
}

function base64ToArrayBuffer(base64) {
    const binary = atob(base64);
    const bytes = new Uint8Array(binary.length);
    for (let i = 0; i < binary.length; i++) {
        bytes[i] = binary.charCodeAt(i);
    }
    return bytes.buffer;
}



// Helper: Convert DER to PEM (with 64-char line breaks)
function convertDerToPem(derBuffer, keyType) {
    const base64 = btoa(
        String.fromCharCode(...new Uint8Array(derBuffer))
    );

    // Split base64 into 64-char lines
    const lines = [];
    for (let i = 0; i < base64.length; i += 64) {
        lines.push(base64.substring(i, i + 64));
    }

    return `-----BEGIN ${keyType}-----\n${lines.join('\n')}\n-----END ${keyType}-----`;
}





async function pemToDer(pemString) {
    // Remove PEM headers and whitespace
    const pemContent = pemString
        .replace(/-----BEGIN [^-]+-----/, '')
        .replace(/-----END [^-]+-----/, '')
        .replace(/\s+/g, '');

    // Base64 decode the content
    const binaryString = atob(pemContent);
    const derBytes = new Uint8Array(binaryString.length);
    for (let i = 0; i < binaryString.length; i++) {
        derBytes[i] = binaryString.charCodeAt(i);
    }
    return derBytes.buffer;
}

async function ecc_384_user_private() {
    const recover = sessionStorage.getItem('user_recover');
    if (!recover) throw new Error("Missing user_recover in session");
    
    const hash_recover = await ECDH_hash256(recover);
    const base64Pem = sessionStorage.getItem('encrypted_private_key');
    
    if (!base64Pem) throw new Error("Missing client private key in session");

    //Convert base64 string directly to ArrayBuffer 
    const encryptedData = base64ToArrayBuffer(base64Pem);
    // Decrypt the ArrayBuffer
    const decryptedBuffer = await decryptWithAESCTR(
        await convertHashToCryptoKey(hash_recover),
        encryptedData
    );

    // FIX 2: Convert decrypted ArrayBuffer to string (PEM content)
    const pemString = new TextDecoder().decode(decryptedBuffer);

    // Convert PEM string to DER
    const derBuffer = await pemToDer(pemString);

    return crypto.subtle.importKey(
        "pkcs8",
        derBuffer,
        { name: "ECDH", namedCurve: "P-384" },
        false,
        ["deriveBits"]
    );
}


async function ecc_384_server_public() {
    const base64Pem = sessionStorage.getItem('server_public_key');
    if (!base64Pem) throw new Error("Missing server public key in session");

    // Decode base64 to get raw PEM string
    const pemString = atob(base64Pem);

    // Convert PEM to DER
    const derBuffer = await pemToDer(pemString);

    // Import as Web Crypto key
    return crypto.subtle.importKey(
        "spki",
        derBuffer,
        { name: "ECDH", namedCurve: "P-384" },
        false,
        []
    );
}

async function getECDHSharedSecret() {
    try {
        const privateKey = await ecc_384_user_private();
        const publicKey = await ecc_384_server_public();

        // 1. DERIVE RAW SHARED SECRET (256 bits = 32 bytes)
        const sharedSecret = await crypto.subtle.deriveBits(
            { name: "ECDH", public: publicKey },
            privateKey,
            256
        );

        // 2. IMPORT AS AES-CTR KEY
        return crypto.subtle.importKey(
            "raw",
            sharedSecret,
            { name: "AES-CTR" },
            false,
            ["encrypt", "decrypt"]
        );

    } catch (error) {
        console.error("ECDH key derivation failed:", error);
        throw error;
    }
}

async function encryptWithAESCTR(key, arrayBuffer) {
    // CTR mode uses a 16-byte counter (96-bit nonce + 32-bit counter)
    const counter = crypto.getRandomValues(new Uint8Array(16));

    try {
        const ciphertext = await crypto.subtle.encrypt(
            { name: "AES-CTR", counter, length: 32 },  // Critical: length=32
            key,
            arrayBuffer
        );

        // Prepend counter to ciphertext
        const result = new Uint8Array(counter.length + ciphertext.byteLength);
        result.set(counter);
        result.set(new Uint8Array(ciphertext), counter.length);

        return result.buffer;
    } catch (error) {
        console.error("Encryption error:", error);
        throw error;
    }
}

async function decryptWithAESCTR(key, arrayBuffer) {
    const data = new Uint8Array(arrayBuffer);

    // CTR counter is 16 bytes
    const counterSize = 16;
    if (data.length < counterSize) {
        throw new Error("Ciphertext too short");
    }

    const counter = data.slice(0, counterSize);
    const ciphertextBytes = data.slice(counterSize);

    try {
        return await crypto.subtle.decrypt(
            { name: "AES-CTR", counter, length: 32 },  // Critical: length=32
            key,
            ciphertextBytes
        );
    } catch (error) {
        console.error("Decryption error:", error);
        throw error;
    }
}

export async function ECDH_encryption(fileArrayBuffer) {
    try {
        const sharedKey = await getECDHSharedSecret();
        return await encryptWithAESCTR(sharedKey, fileArrayBuffer);
    } catch (error) {
        console.error("ECDH encryption failed:", error);
        throw new Error(`ECDH encryption failed: ${error.message}`);
    }
}

export async function ECDH_decryption(fileArrayBuffer) {
    try {
        const sharedKey = await getECDHSharedSecret();
        return await decryptWithAESCTR(sharedKey, fileArrayBuffer);
    } catch (error) {
        console.error("ECDH decryption failed:", error);
        throw new Error(`ECDH decryption failed: ${error.message}`);
    }
}
