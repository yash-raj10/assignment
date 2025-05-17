
## Features

- **Generate Report**: Generates a report for a given user with a unique report ID and timestamp.
- **Health Check**: Simple health endpoint to verify the service is running.

---

## Requirements

- Go 1.20+
- `protoc` (Protocol Buffers compiler)
- `protoc-gen-go` and `protoc-gen-go-grpc` plugins installed
- [grpcurl](https://github.com/fullstorydev/grpcurl) for testing (or Postman with gRPC support)

---

## How to Run

### 1. Clone the Repository

```bash
git clone https://github.com/your-username/grpc-report-service.git
cd grpc-report-service
```

### 2. Generate gRPC Code

```bash
protoc --go_out=. --go-grpc_out=. pb/assignment.proto
```

### 3. Run the server

```bash
go run server/main.go
```
---

## How to Test

### 1. Health Check
- Method: GetHealth
- Request Type: Google.protobuf.Empty
- Response Type:

```bash
{
  "status": "ok"
}
```
- Test with grpcurl:
  
```bash
{
  grpcurl -plaintext localhost:8080 assignment.AssignmentService/GetHealth
}
```

### 2. Generate Report
- Method: GenerateReport
- Request Type:
```bash
{
  "user_id": "user1"
}
```
- Response Type:
```bash
{
  "user_id": "user1",
  "report_id": "user1-12345678",
  "created_at": "2025-05-17T12:34:56Z"
}

```
- Test with grpcurl:
  
```bash
{
 grpcurl -plaintext -d '{"user_id":"user1"}' localhost:8080 assignment.AssignmentService/GenerateReport

}
```

---
## Note
- Reports are stored in memory and refreshed automatically by a cron job every 10 seconds for users user1, user2, and user3.
- The answer to the scaling design is provided in design.md.
