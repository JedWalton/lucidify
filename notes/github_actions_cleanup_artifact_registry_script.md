# List images sorted by creation date, get all but the latest three, and delete them
gcloud container images list-tags europe-docker.pkg.dev/sheffessions/sheffessions-docker-repository/sheffessions_api \
  --limit=999999 --sort-by=TIMESTAMP \
  --format='get(digest)' | tail -n +4 | xargs -I{} gcloud container images delete -q --force-delete-tags europe-docker.pkg.dev/sheffessions/sheffessions-docker-repository/sheffessions_api@{}

