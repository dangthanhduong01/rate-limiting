# Rate limiting

# Setup redis:

<code>
    docker compose up
</code>

# Run:

<code>
    go run .
</code> 

# Build

`go run build`

# Test:

1. Test all

<code>
    curl -H "X-User-ID: user1" http://localhost:8080/fixed

    curl -H "X-User-ID: user1" http://localhost:8080/sliding

    curl -H "X-User-ID: user1" http://localhost:8080/token
</code>

2. Fixed Window

<code>
    for i in {1..6}; do
    curl -s -o /dev/null -w "%{http_code}\n" -H "X-User-ID: user123" http://localhost:8080/fixed
    done
</code>

3. Sliding Window
<code>
    for i in {1..11}; do
    curl -s -o /dev/null -w "%{http_code}\n" -H "X-User-ID: user456" http://localhost:8080/sliding
    done
</code>

4. Token Bucket
<code>
    for i in {1..12}; do
    curl -s -o /dev/null -w "%{http_code}\n" -H "X-User-ID: user789" http://localhost:8080/token
    done
</code>