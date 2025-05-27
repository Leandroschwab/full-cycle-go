#!/bin/bash

# Script to test token-based rate limiting
echo "Testing token-based rate limiting..."
echo "Sending requests with API token to http://localhost:8080/"

# API token to use
API_TOKEN="test-token-123"

# Number of requests to send (greater than the token limit)
NUM_REQUESTS=500

# Counter for successful requests
SUCCESS_COUNT=0

# Use a POSIX-compliant loop syntax
i=1
while [ $i -le $NUM_REQUESTS ]; do
    # Send request with API token and capture response code
    RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" -H "API_KEY: $API_TOKEN" http://localhost:8080/)
    
    if [ "$RESPONSE" -eq 200 ]; then
        SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
        echo "Request $i: Success (HTTP $RESPONSE)"
    else
        echo "Request $i: Failed (HTTP $RESPONSE)"
    fi
    
    # Small delay to make output readable
    sleep 0.0005    

    # Increment counter
    i=$((i + 1))
done

echo "Test completed. $SUCCESS_COUNT/$NUM_REQUESTS requests were successful."
echo "Token-based limit should be reached after configured number of requests."
