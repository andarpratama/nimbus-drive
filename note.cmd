docker compose up --build


curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Andar", "email":"andar@example.com", "password":"secret123"}'
