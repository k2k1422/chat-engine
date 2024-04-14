import argparse
import json
import hashlib
import hmac
import base64
from datetime import datetime, timedelta

def base64url_encode(data):
    encoded = base64.urlsafe_b64encode(data)
    return encoded.decode('utf-8').rstrip('=')

def create_jwt(username, secret_key):
    header = {"alg": "HS256", "typ": "JWT"}
    header_encoded = base64url_encode(json.dumps(header).encode('utf-8'))

    exp_time = int((datetime.utcnow() + timedelta(days=1)).timestamp())
    payload = {"username": username, "exp": exp_time}
    payload_encoded = base64url_encode(json.dumps(payload).encode('utf-8'))

    signature = hmac.new(secret_key.encode('utf-8'), (header_encoded + "." + payload_encoded).encode('utf-8'), hashlib.sha256)
    signature_encoded = base64url_encode(signature.digest())

    jwt_token = header_encoded + "." + payload_encoded + "." + signature_encoded

    return jwt_token

def main():
    parser = argparse.ArgumentParser(description='Generate JWT token with custom values.')
    parser.add_argument('username', type=str, help='Username')
    parser.add_argument('--secret', type=str, default='f1152577d55a836ea26843b0433059ed9ba1add93e0e49767890c2c46852b8d8', help='Secret key for signing the token')

    args = parser.parse_args()

    token = create_jwt(args.username, args.secret)
    print("JWT Token:", token)

if __name__ == "__main__":
    main()
