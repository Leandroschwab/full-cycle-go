#!/bin/bash

# Script to test rate limiting with concurrent requests
echo "Testing rate limiting with concurrent requests..."
echo "Sending concurrent requests to http://localhost:8080/"

# Number of concurrent requests
CONCURRENCY=5

# Number of requests per client
REQUESTS_PER_CLIENT=10

# Function to send requests and count successes
send_requests() {
    local client_id=$1
    local success_count=0
    
    # Use POSIX-compliant while loop
    i=1
    while [ $i -le $REQUESTS_PER_CLIENT ]; do
        RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/)
        
        if [ "$RESPONSE" -eq 200 ]; then
            success_count=$((success_count + 1))
        fi
        
        # Small delay between requests
        sleep 0.1
        i=$((i + 1))
    done
    
    echo "Client $client_id: $success_count/$REQUESTS_PER_CLIENT requests successful"
}

# Start concurrent clients with POSIX-compliant while loop
client=1
while [ $client -le $CONCURRENCY ]; do
    send_requests $client &
    client=$((client + 1))
done

# Wait for all background processes to complete
wait

echo "Concurrent test completed."
echo "With high concurrency, rate limits should be reached quickly."
