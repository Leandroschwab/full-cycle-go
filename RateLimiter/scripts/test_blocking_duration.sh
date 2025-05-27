#!/bin/bash

# Script to test the blocking duration feature
echo "Testing blocking duration..."
echo "This test will try to trigger the rate limit and then check if requests stay blocked"

# API token to use
API_TOKEN="test-token-123"

# Function to send requests until we get a 429 response
trigger_limit() {
    local req_count=0
    local blocked=false
    
    echo "Sending requests until rate limit is triggered..."
    
    while [ "$blocked" = false ] && [ $req_count -lt 200 ]; do
        req_count=$((req_count + 1))
        
        RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" -H "API_KEY: $API_TOKEN" http://localhost:8080/)
        
        if [ "$RESPONSE" -eq 429 ]; then
            blocked=true
            echo "Rate limit triggered after $req_count requests (HTTP 429)"
        else
            echo -n "."
        fi
        
        # Small delay
        sleep 0.00005
    done
    
    echo ""
    return 0
}

# Function to check if requests are still blocked
check_if_blocked() {
    RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" -H "API_KEY: $API_TOKEN" http://localhost:8080/)
    
    if [ "$RESPONSE" -eq 429 ]; then
        echo "Requests are still blocked (HTTP 429)"
        return 0
    else
        echo "Requests are allowed again (HTTP $RESPONSE)"
        return 1
    fi
}

# Trigger the rate limit
trigger_limit

# Record start time
START_TIME=$(date +%s)

# Check immediately after triggering
echo "Checking if requests are blocked immediately after hitting the limit..."
check_if_blocked

# Check every second until no longer blocked
echo "Will check every second to see when the block expires..."

# Use a POSIX-compliant loop
seconds=1
while [ $seconds -le 300 ]; do  # Check for up to 5 minutes (300 seconds)
    sleep 1
    
    echo "Checking after $seconds second(s)..."
    check_if_blocked
    
    if [ $? -eq 1 ]; then
        END_TIME=$(date +%s)
        DURATION=$((END_TIME - START_TIME))
        echo "Block has expired after $DURATION seconds"
        break
    fi
    
    seconds=$((seconds + 1))
done

if [ $seconds -gt 300 ]; then
    echo "Block didn't expire within the 5-minute test period"
fi

echo "Test completed."
