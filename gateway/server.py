import logging
import os, json, requests
from flask import Flask, request, redirect as next
from flask_cors import CORS

server = Flask(__name__)

server.config['MAX_CONTENT_LENGTH'] = 4 * 1024 * 1024

server.logger.setLevel(logging.DEBUG)

# CORS(server, resources={r"/*": {"origins": os.environ.get("CLIENT_URL"), "supports_credentials": True}})
@server.before_request
def validate():
    if request.path.startswith(('/auth', '/user', '/hello', '/redirect')):
        return None
    server.logger.debug("validate request")
    authorization = request.headers["Authorization"]
    response = requests.get(f"{os.environ.get("AUTH_ENDPOINT")}/validate", headers={'Authorization': authorization})
    if response.status_code != 200:
        return (response.json(), 401)
    return None

# HELPER ENDPOINTS
@server.route("/hello", methods=["GET"])
def hello():
    server.logger.debug("hello request")
    return "hello world", 400

@server.route("/redirect", methods=["GET"])
def redirect():
    server.logger.debug("redirect request")
    url = request.url.replace("http", os.environ.get("CLIENT_PROTOCOL"), 1)
    return next(url)

# AUTH/USER SERVICE ENDPOINTS
@server.route("/auth/login", methods=["POST"])
def login():
    server.logger.debug("login request")
    credentials = json.loads(request.data)
    response = requests.post(f"{os.environ.get("AUTH_ENDPOINT")}/login", json=credentials)
    return response.json(), response.status_code, {'Content-Type': 'application/json'}
    
@server.route("/auth/signup", methods=["POST"])
def signup():
    server.logger.debug("signup request")
    credentials = json.loads(request.data)
    response = requests.post(f"{os.environ.get("AUTH_ENDPOINT")}/signup", json=credentials)
    return response.json(), response.status_code, {'Content-Type': 'application/json'}
    
@server.route("/auth/refresh", methods=["POST"])
def refresh():
    server.logger.debug("refresh request")
    response = requests.post(f"{os.environ.get("AUTH_ENDPOINT")}/refresh", headers=request.headers)
    return response.json(), response.status_code, {'Content-Type': 'application/json'}

# VISION SERVICE ENDPOINTS
@server.route("/connectAccount", methods=["POST"])
def connectAccount():
    server.logger.debug("connectAccount request")
    params = json.loads(request.data)
    response = requests.post(f"{os.environ.get("VISION_ENDPOINT")}/connectAccount", json=params)
    return response.json(), response.status_code, {'Content-Type': 'application/json'}

@server.route("/searchByImage", methods=["POST"])
def searchByImage():
    server.logger.debug("searchByImage request")
    params = json.loads(request.data)
    response = requests.post(f"{os.environ.get("VISION_ENDPOINT")}/search", json=params)
    return response.json(), response.status_code, {'Content-Type': 'application/json'}

if __name__ == "__main__":
    server.run(host='0.0.0.0', port=8080)