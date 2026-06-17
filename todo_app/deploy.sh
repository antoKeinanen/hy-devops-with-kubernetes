docker build -t todo-app-frontend:latest todo-frontend
docker build -t todo-app-backend:latest todo-backend
k3d image import todo-app-frontend:latest todo-app-backend:latest
kubectl apply -f manifests

