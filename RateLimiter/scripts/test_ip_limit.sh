#!/bin/bash

# Script to test IP-based rate limiting
echo "Testing IP-based rate limiting..."
echo "Sending requests to  http://localhost:8080/"

# Number of requests to send (greater than the IP limit)
NUM_REQUESTS=500

# Counter for successful requests
SUCCESS_COUNT=0

# Use POSIX-compliant while loop instead of bash-specific for loop
i=1
while [ $i -le $NUM_REQUESTS ]; do
    # Send request and capture response code
    RESPONSE=$(curl -s -o /dev/null -w "%{http_code}"  http://localhost:8080/)
    
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
echo "IP-based limit should be reached after configured number of requests."
