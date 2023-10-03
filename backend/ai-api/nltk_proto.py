import nltk
nltk.download('stopwords')
nltk.download('brown')

from nltk.corpus import brown
from nltk.tokenize.texttiling import TextTilingTokenizer

# Initialize the TextTilingTokenizer
tt = TextTilingTokenizer()

# Get a sample text from the Brown corpus
text = brown.raw()[:4000]

# Tokenize the text into topical sections
segments = tt.tokenize(text)

# Print the segments
for idx, segment in enumerate(segments, 1):
    print(f"Segment {idx}:\n{segment}\n{'-'*40}")
