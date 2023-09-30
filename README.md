docker build . -t localhost:5000/jean/rinha-backend:latest && docker push localhost:5000/jean/rinha-backend:latest

docker-compose -f .\docker-compose.yml up -d

