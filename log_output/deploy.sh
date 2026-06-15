docker build -t log-output:latest .
k3d image import log-output:latest
kubectl apply -f deployment.yaml

