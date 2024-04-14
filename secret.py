import secrets

def generate_jwt_secret():
    # Generate a random JWT secret of 32 bytes (256 bits)
    jwt_secret = secrets.token_hex(32)
    
    return jwt_secret

def main():
    jwt_secret = generate_jwt_secret()
    print("Generated JWT Secret:", jwt_secret)

if __name__ == "__main__":
    main()
