# crawler.py
from flask import Blueprint

# ENV var
from dotenv import load_dotenv
import os
load_dotenv()
SECRET_KEY = os.getenv("X_AI_API_KEY")

crawler_blueprint = Blueprint('crawler', __name__)

@crawler_blueprint.route('/crawl', methods=['GET'])
def crawl_route():
    # Your crawling logic here
    return "Crawler route"

@crawler_blueprint.errorhandler(500)
def internal_error(error):
    app.logger.error('Server Error: %s', (error))
    return jsonify({"error": "Internal server error", "details": str(error)}), 500

@crawler_blueprint.errorhandler(404)
def not_found(error):
    app.logger.error('Not Found: %s', (error))
    return "404 error"
