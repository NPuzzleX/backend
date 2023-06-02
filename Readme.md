To serveapply
- go get .
- go run .

To docker
- docker build -t npuzzlex-backend .
- docker run -dp 8080:8080 npuzzlex-backend

To kubernetes
Make sure kubectl points to the correct context (desktop dev or DO k8s deploy)
- kubectl apply -f deployment.yml
- kubectl apply -f services.yml
- Check kubectl get deploy(services)