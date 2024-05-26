# Run
go run main.go


# Put data 
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"key":"test1","value":"test11"}' \
  http://localhost:8080/data

curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"key":"test2","value":"test2"}' \
  http://localhost:8080/data

-- After snapshot  
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"key":"test3","value":"test3"}' \
  http://localhost:8080/data

curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"key":"test4","value":"test4"}' \
  http://localhost:8080/data


  curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"key":"test5","value":"test5"}' \
  http://localhost:8080/data

# Get all data
  curl http://localhost:8080/all


# Take snapshot
  curl  --request POST http://localhost:8080/snapshot


# Get all data snapshot
  curl http://localhost:8080/all-snapshot
