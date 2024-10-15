import base64
import os
import time
import jwt
import requests

endpoint = "https://api.ebay.com/identity/v1/oauth2/token"
scope = "https://api.ebay.com/oauth/api_scope"
client_id = os.environ.get("EBAY_CLIENT_ID")
client_secret = os.environ.get("EBAY_CLIENT_SECRET")

client_token = None

def isTokenValid(token):
    try:
        claims = jwt.decode(token, os.environ.get("EBAY_CLIENT_SECRET"), algorithms=["HS256"])
        exp = claims.get("exp")

        if exp is not None:
            return time.time() > exp - 300

        return False
    except:
        return False

def refreshAuth(token):
    global endpoint
    global scope
    global client_id
    global client_secret

    if not client_id or not client_secret:
        raise ValueError("Client ID and Client Secret must be set in the environment variables")

    auth_str = f"{client_id}:{client_secret}"
    
    auth_bytes = base64.b64encode(auth_str.encode("utf-8"))
    auth_header = auth_bytes.decode("utf-8")
    
    headers = {
        "Content-Type": "application/x-www-form-urlencoded",
        "Authorization": f"Basic {auth_header}"
    }
    
    data = {
        "grant_type": "refresh_token",
        "refresh_token": token,
        "scope": scope
    }
    
    response = requests.post(endpoint, headers=headers, data=data)
    
    if response.status_code != 200:
        response.raise_for_status()

    return response.json()["access_token"]

def getCredentialsToken():
    global endpoint
    global scope
    global client_id
    global client_secret
    
    if not client_id or not client_secret:
        raise ValueError("Client ID and Client Secret must be set in the environment variables")

    auth_str = f"{client_id}:{client_secret}"
    
    auth_bytes = base64.b64encode(auth_str.encode("utf-8"))
    auth_header = auth_bytes.decode("utf-8")
    
    headers = {
        "Content-Type": "application/x-www-form-urlencoded",
        "Authorization": f"Basic {auth_header}"
    }
    
    data = {
        "grant_type": "client_credentials",
        "scope": scope
    }
    
    response = requests.post(endpoint, headers=headers, data=data)
        
    if response.status_code != 200:
            response.raise_for_status()

    return response.json()["access_token"]
        

def validate(access_token=None, refresh_token=None):
    global client_token
    try:
        if access_token is not None:
            if isTokenValid(access_token):
                return access_token, False
        
            if refresh_token is not None:
                return refreshAuth(refresh_token), True
        
        if isTokenValid(client_token):
                return client_token, False
        
        client_token = getCredentialsToken()
        return client_token, False
    except:
        return None, False

    

