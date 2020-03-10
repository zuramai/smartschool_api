# Installation

1. Clone the repository
  ```bash
  go get https://github.com/zuramai/smartschool_api.git
  git clone https://github.com/zuramai/smart_school_webapi.git
  ```

2. Run the mongodb
  ```bash
  sudo service mongod start
  ```

3. Run the API
  ```bash
  cd $GOPATH/src/github.com/zuramai/smartschool_api
  go run main.go
  ```

4. Run the Web
  ```bash
  cd /var/smart_school_webapi
  php artisan serve
  ```

5. Run the Socket Server
  ```bash
  cd /var/smart_school_webapi/websocket
  node index
  ```

6. Run the Attendance Detection
  ```bash
  cd $GOPATH/src/github.com/zuramai/smartschool_api/demo/CyberNet
  go run bundler.go
  ```

## Registration
1. Go to directory
```bash
cd $GOPATH/src/github.com/zuramai/smartschool_api/demo/Register
go run Binary.go
```
