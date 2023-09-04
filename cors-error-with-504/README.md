
# Start the frontend. 
```sh
cd express
npm i
node server.js
```

# Start the Gateway
```sh
# install nginx if not. If on Mac, run `brew install nginx`
cd nginx
ginx -c $(pwd)/proxy-8083-hello.conf
```

# Start the origin server.
```sh
cd goweb
go run main.go
```

