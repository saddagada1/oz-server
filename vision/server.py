import base64
import os, logging, requests, json
from flask import Flask, request
from flask_cors import CORS

server = Flask(__name__)

server.logger.setLevel(logging.DEBUG)

credentials = {"accessToken": ""}

def getCredentials():
    global credentials

    endpoint = "https://api.sandbox.ebay.com/identity/v1/oauth2/token" if os.environ.get("ENV") == "DEV" else "https://api.ebay.com/identity/v1/oauth2/token"

    client_id = os.environ.get("EBAY_CLIENT_ID")
    client_secret = os.environ.get("EBAY_CLIENT_SECRET")
    
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
        "scope": "https://api.ebay.com/oauth/api_scope"
    }
    
    response = requests.post(endpoint, headers=headers, data=data)
    
    if response.status_code == 200:
        credentials["accessToken"] = response.json()["access_token"]
    else:
        server.logger.error(response)
        response.raise_for_status()


@server.route("/search", methods=["POST"])
def search():
    server.logger.debug("search request")

    global credentials

    if credentials["accessToken"] == "":
        getCredentials();
        
    body = json.loads(request.data)
    image = body.get("image", None)

    if not image:
        return 'no image', 400
    
    endpoint = "https://api.ebay.com/buy/browse/v1/item_summary/search_by_image"

    headers = {
        "Authorization": f"Bearer {credentials.access_token}",
        "Content-Type": "application/json"
    }
    
    response = requests.post(endpoint, headers=headers, data={image})
    
    if response.status_code == 200:
        return response.json()
    else:
        server.logger.error(response)
        response.raise_for_status()

if __name__ == "__main__":
    server.run(host='0.0.0.0', port=5001)