docker build -t log-output-writer:latest writer
docker build -t log-output-server:latest server
k3d image import log-output-writer:latest log-output-server:latest
kubectl apply -f manifests

