from flask import Flask, request, jsonify, abort
import os
from dotenv import load_dotenv

load_dotenv()

SECRET_KEY = os.getenv("AI_API_SECRET_KEY")
print(f"Loaded SECRET_KEY from .env: {SECRET_KEY}")  # Diagnostic print

app = Flask(__name__)

@app.route('/split_sentences', methods=['POST'])
def split_sentences():
    try:
        secret_key = request.headers.get('X-AI-API-KEY')
        print(f"Received AI_API_SECRET_KEY header: {secret_key}")  # Diagnostic print

        if not secret_key or secret_key != SECRET_KEY:
            print("Unauthorized due to mismatched or missing secret key.")  # Diagnostic print
            abort(401, description="Unauthorized")

        text = request.json.get('text')
        if not text:
            return jsonify({"error": "No text provided"}), 400


        return jsonify(text)
    except Exception as e:
        print(f"Exception occurred: {str(e)}")  # Diagnostic print
        return jsonify({"error": str(e)}), 500

@app.errorhandler(500)
def internal_error(error):
    app.logger.error('Server Error: %s', (error))
    return jsonify({"error": "Internal server error", "details": str(error)}), 500

@app.errorhandler(404)
def not_found(error):
    app.logger.error('Not Found: %s', (error))
    return "404 error"

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, debug=True)

