import logging
import os, json, requests
from flask import Flask, request
from flask_cors import CORS

server = Flask(__name__)

server.logger.setLevel(logging.DEBUG)

# CORS(server, resources={r"/*": {"origins": os.environ.get("CLIENT_URL"), "supports_credentials": True}})

# AUTH SERVICE ENDPOINTS

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
    response = requests.post(f"{os.environ.get("AUTH_ENDPOINT")}/refresh", request.data, headers=request.headers)
    return response.json(), response.status_code, {'Content-Type': 'application/json'}
    
@server.route("/validate", methods=["GET"])
def validate():
    server.logger.debug("validate request")
    response = requests.get(f"{os.environ.get("AUTH_ENDPOINT")}/validate", request.data, headers=request.headers)
    return response.json(), response.status_code, {'Content-Type': 'application/json'}

@server.route("/hello", methods=["GET"])
def hello():
    server.logger.debug("hello request")
    return "hello world", 400


# VISION SERVICE ENDPOINTS
@server.route("/searchByImage", methods=["POST"])
def searchByImage():
    server.logger.debug("searchByImage request")
    data = json.loads(request.data)
    response = requests.post(f"{os.environ.get("VISION_ENDPOINT")}/search", json=data)
    server.logger.debug(response)
    return response, response.status_code, {'Content-Type': 'application/json'}

if __name__ == "__main__":
    server.run(host='0.0.0.0', port=8080)