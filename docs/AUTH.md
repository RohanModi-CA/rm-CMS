# Authentication & Secrets

RM-CMS uses client-side RSA encryption. No server required.

## 1. The Helper API
Import the helper and fetch any secret from `site/auth/data_jsons/`.

```javascript
import { rmcmsAuth } from "/auth/rmcms-auth.js";

// Decrypt a secret from site/auth/data_jsons/ai_config.json
const apiKey = await rmcmsAuth.getProperty("ai_config.json", "gemini_api_key");
```

## 2. The Secret Manager
Use the built-in utility to encrypt new secrets:
`http://localhost:8080/auth/manager.html`

1. Select your user.
2. Enter a property name and the plaintext secret.
3. Click Encrypt and copy the resulting JSON.
4. Save it as a `.json` file in `site/auth/data_jsons/`.

## 3. Requirements
- **Login:** User must be logged in via `login.html`.
- **Session:** Keys are stored in **IndexedDB** and survive page refreshes.
- **Context:** Must be served over **HTTPS** or **localhost**.

## 4. Data Structure
Files in `site/auth/data_jsons/` follow this format:
```json
{
    "users": {
        "<username_hash>": {
            "property_name": "encrypted_base64_blob"
        }
    }
}
```
