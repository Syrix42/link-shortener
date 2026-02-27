# Local Secrets (DO NOT COMMIT)

This folder is for **local development only**. It contains RSA keys used to sign/verify JWTs.

✅ **Keep the folder** in git so everyone knows where to put keys.  
❌ **Never commit private keys**. All `.pem` files in this folder are ignored by `.gitignore`.

## Files expected here

Create these files:

- `access_private.pem`  (RSA private key, used to sign Access tokens)
- `access_public.pem`   (RSA public key, used to verify Access tokens)
- `refresh_private.pem` (RSA private key, used to sign Refresh tokens)
- `refresh_public.pem`  (RSA public key, used to verify Refresh tokens)

## Generate keys (example)

### Access token keypair
```bash
# private key (PKCS#1)
openssl genrsa -out secrets/access_private.pem 2048

# public key
openssl rsa -in secrets/access_private.pem -pubout -out secrets/access_public.pem
