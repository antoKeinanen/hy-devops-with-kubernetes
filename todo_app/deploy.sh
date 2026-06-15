docker build -t todo-app:latest .
k3d image import todo-app:latest
kubectl apply -f manifests

