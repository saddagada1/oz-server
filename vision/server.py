import logging, requests, json
from flask import Flask, request
from flask_cors import CORS
from lib import auth

server = Flask(__name__)

server.config['MAX_CONTENT_LENGTH'] = 4 * 1024 * 1024

server.logger.setLevel(logging.DEBUG)

@server.route("/search", methods=["POST"])
def search():
    server.logger.debug("search request")

    body = request.json

    access_token = body.get("accessToken", None)
    refresh_token = body.get("refreshToken", None)

    token, update = auth.validate(access_token, refresh_token)
        
    image = body.get("image", None)

    if not image:
        return 'no image', 400
    
    endpoint = "https://api.ebay.com/buy/browse/v1/item_summary/search_by_image?limit=5"

    headers = {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/json"
    }
    
    response = requests.post(endpoint, headers=headers, json={"image": image})
    
    if response.status_code != 200:
        response.raise_for_status()
        
    return response.json()
        

if __name__ == "__main__":
    server.run(host='0.0.0.0', port=5001)