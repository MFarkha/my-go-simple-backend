import requests

# URL of the server to send requests to
url = 'http://localhost:3000/metrics'

# Number of requests to send
num_requests = 100

# Send HTTP GET requests
for i in range(num_requests):
    try:
        response = requests.get(url)
        # Print response status code for each request
        print(f'Request {i+1}: Status code - {response.status_code}')
        response.json()

    except requests.exceptions.RequestException as e:
        print(f'Error sending request {i+1}: {e}')
