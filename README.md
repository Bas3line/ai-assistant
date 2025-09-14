# API Testing Commands

## Public Endpoints

### Root & Health
```bash
curl http://localhost:8000/

curl http://localhost:8000/health
```

### Authentication
```bash
# Start Google OAuth
curl http://localhost:8000/api/auth/google

# OAuth Callback (get JWT token)
curl "http://localhost:8000/api/auth/google/callback?code=test_code&state=test_state"
```

## Protected Endpoints

### Get Token First
```bash
export TOKEN=$(curl -s "http://localhost:8000/api/auth/google/callback?code=test_code&state=test_state" | jq -r '.token')
```

### User Authentication
```bash
# Get current user info
curl -H "Authorization: Bearer $TOKEN" \
     http://localhost:8000/api/auth/me

# Logout
curl -X POST \
     -H "Authorization: Bearer $TOKEN" \
     http://localhost:8000/api/auth/logout
```

### AI Endpoints
```bash
# Ask AI question
curl -X POST \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $TOKEN" \
     -d '{
       "prompt": "What is the capital of France?",
       "provider": "gemini"
     }' \
     http://localhost:8000/api/ai/ask
```

### Email Endpoints
```bash
# Get emails
curl -H "Authorization: Bearer $TOKEN" \
     http://localhost:8000/api/emails/

# Send email
curl -X POST \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $TOKEN" \
     -d '{
       "to": "test@example.com",
       "subject": "Test Email",
       "body": "Hello World"
     }' \
     http://localhost:8000/api/emails/send
```

## Manual Token Testing
```bash
export TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZGVtb191c2VyXzEyMyIsImVtYWlsIjoiZGVtb0BleGFtcGxlLmNvbSIsImV4cCI6MTc1NzkxMjcwMywibmJmIjoxNzU3ODI2MzAzLCJpYXQiOjE3NTc4MjYzMDN9.Rjv0om9EvtEdeO9iO_VYDv2_bqqT1XM9puD2jvIrAjQ"

curl -H "Authorization: Bearer $TOKEN" http://localhost:8000/api/auth/me
```