# chunker.py
from flask import Blueprint
from flask import Flask, request, jsonify, abort

# ENV Var
from dotenv import load_dotenv
import os
load_dotenv()
SECRET_KEY = os.getenv("X_AI_API_KEY")

import nltk
nltk.download('stopwords')
from nltk.tokenize.texttiling import TextTilingTokenizer

chunker_blueprint = Blueprint('chunker', __name__)

@chunker_blueprint.route('/split_text_to_chunks', methods=['POST'])
def split_sentences():
    try:
        secret_key = request.headers.get('X-AI-API-KEY')
        print(f"Received X_AI_API_KEY header: {secret_key}")  # Diagnostic print

        if not secret_key or secret_key != SECRET_KEY:
            print("Unauthorized due to mismatched or missing secret key.")  # Diagnostic print
            abort(401, description="Unauthorized")

        text = request.json.get('text', "")
        if not text:
            return jsonify({"error": "No text provided"}), 400

        tt = TextTilingTokenizer(w=20, k=10)

        try:
            segments = tt.tokenize(text)
        except ValueError as e:
            if str(e) == "No paragraph breaks were found(text too short perhaps?)":
                segments = [text]  # Return the same text as a single segment
            else:
                raise e  # If it's a different error, raise it to be caught by the outer exception handler

        return jsonify(segments)
    except Exception as e:
        print(f"Exception occurred: {str(e)}")  # Diagnostic print
        return jsonify({"error": str(e)}), 500

@chunker_blueprint.errorhandler(500)
def internal_error(error):
    app.logger.error('Server Error: %s', (error))
    return jsonify({"error": "Internal server error", "details": str(error)}), 500

@chunker_blueprint.errorhandler(404)
def not_found(error):
    app.logger.error('Not Found: %s', (error))
    return "404 error"
