# initialize
./app &
export BUILD_PID=$!
go build test/callback_server.go
sleep 7

# send test request
curl --form file=@example/app.tar --form callback=http://localhost:8888/callback http://localhost:$PORT/builds

# wait callback
PORT=8888 ./callback_server
export RET=$?

# finish test
kill $BUILD_PID
exit $RET
