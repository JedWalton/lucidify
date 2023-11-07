from flask import Flask, request, jsonify, abort

import re

# Blueprints
from chunker import chunker_blueprint

app = Flask(__name__)
app.register_blueprint(chunker_blueprint, url_prefix='/chunker')

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, debug=True)

