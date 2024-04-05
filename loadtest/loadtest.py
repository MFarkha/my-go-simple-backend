import requests

# URL of the server to send requests to
URL = "http://localhost:3000"

# Number of requests to send
NUM_REQUESTS = 100

# Default timeout
TIMEOUT = 5

# Endpoint to hit
endpoints = ["health", "ready", "payload", "metrics"]

# Send HTTP GET requests for each endpoint
for i in range(NUM_REQUESTS):
    for endpoint in endpoints:
        url_endpoint = URL + "/" + endpoint
        try:
            response = requests.get(url_endpoint, timeout=TIMEOUT)
            # Print response status code for each request
            # print(f"Request {i+1}, EndPoint: {endpoint}, Status code - {response.status_code}")

        except requests.exceptions.RequestException as e:
            print(f"Error sending request {i+1} for the endpoint {endpoint}: {e}")

# Print metrics summary
try:
    response = requests.get(URL + "/metrics", timeout=TIMEOUT)
    # print(response.json())
    FORMAT_SPEC = "{:<15}"
    header_line = FORMAT_SPEC.format("Endpoint")
    header_line += FORMAT_SPEC.format("RequestCount")
    header_line += FORMAT_SPEC.format("TotalDuration")
    header_line += FORMAT_SPEC.format("AverageLatency")
    print(header_line)
    metrics = response.json()
    for endpoint in metrics:
        data_line = FORMAT_SPEC.format(endpoint)
        data_line += FORMAT_SPEC.format(metrics[endpoint]["RequestCount"])
        data_line += FORMAT_SPEC.format(metrics[endpoint]["TotalDuration"])
        data_line += FORMAT_SPEC.format(metrics[endpoint]["AverageLatency"])
        print(data_line)

except requests.exceptions.RequestException as e:
    print(f"Error sending request for metrics endpoint: {e}")
