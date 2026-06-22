docker build -t log-output-writer:latest writer
docker build -t log-output-server:latest server
k3d image import log-output-writer:latest log-output-server:latest
kubectl apply -f manifests
kubectl -n exercises rollout restart deployment/log-output
kubectl -n exercises rollout status deployment/log-output
