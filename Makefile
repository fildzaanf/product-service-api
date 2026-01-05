# --------------------------------------
# File name Grafana K6 load testing 
# --------------------------------------
K6_TEST_REST=test/load-testing/rest-test.js
K6_TEST_GRPC=test/load-testing/grpc-test.js

# --------------------------------------
# Run REST API load test
# --------------------------------------
k6-rest:
	@echo "Running Grafana K6 REST load test..."
	k6 run $(K6_TEST_REST)

# --------------------------------------
# Run gRPC load test
# --------------------------------------
k6-grpc:
	@echo "Running Grafana K6 gRPC load test..."
	k6 run $(K6_TEST_GRPC)

# --------------------------------------
# Run both tests sequentially
# --------------------------------------
k6-all: k6-rest k6-grpc
