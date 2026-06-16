docker build -t pingpong:latest .
k3d image import pingpong:latest
kubectl apply -f manifests

