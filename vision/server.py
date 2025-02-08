import logging, requests, json
from flask import Flask, request
from flask_cors import CORS
from lib import auth

server = Flask(__name__)

server.config['MAX_CONTENT_LENGTH'] = 4 * 1024 * 1024

server.logger.setLevel(logging.DEBUG)

session = requests.Session()
session.headers = {
    "Content-Type": "application/json"
}

@server.route("/connectAccount", methods=["POST"])
def connectAccount():
    server.logger.debug("connect account request")

    body = request.json

    code = body.get("code", None)

    if not code:
        return 'no code', 400
    
    return auth.connectAccount(code)

@server.route("/search", methods=["POST"])
def search():
    server.logger.debug("search request")

    global session
    body = request.json

    access_token = body.get("accessToken", None)
    refresh_token = body.get("refreshToken", None)

    token = auth.validate(access_token)
        
    image = body.get("image", None)

    if not image:
        return 'no image', 400
    
    endpoint = "https://api.ebay.com/buy/browse/v1/item_summary/search_by_image?limit=5"

    session.headers.update( {
        "Authorization": f"Bearer {token}",
    })
    
    payload = session.prepare_request(requests.Request("POST", url=endpoint, json={"image": image}))
    response = session.send(payload)
    
    if response.status_code != 200:
        return handleApiError(response, refresh_token, payload)
        
    return response.json(), None

def handleApiError(response, refresh_token, request_payload):
    server.logger.debug("api error")

    global session

    if response.status_code == 401:
        server.logger.debug("auth error")

        token, update = auth.authorize(refresh_token)

        server.logger.debug(token)
        session.headers.update({
            "Authorization": f"Bearer {token}",
        })

        response = session.send(request_payload)

        if response.status_code != 200:
            return response.raise_for_status(), None
        
        if update != True:
            return response.json(), None
        
        return response.json(), token
        
    return response.raise_for_status(), None

        

if __name__ == "__main__":
    server.run(host='0.0.0.0', port=5001)