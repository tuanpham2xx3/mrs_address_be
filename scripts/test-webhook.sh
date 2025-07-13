#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Function to print colored output
print_color() {
    echo -e "${1}${2}${NC}"
}

# Function to show usage
show_usage() {
    print_color $BLUE "🚀 Webhook Test Tool for Vietnam Admin API"
    echo "Usage: $0 <webhook_url> <webhook_secret> [environment]"
    echo "Example: $0 https://webhook1.iceteadev.site/ my_secret_key production"
    echo ""
    echo "Environments:"
    echo "  - staging (default)"
    echo "  - production"
    exit 1
}

# Function to generate HMAC signature
generate_signature() {
    local payload=$1
    local secret=$2
    echo -n "$payload" | openssl dgst -sha256 -hmac "$secret" | cut -d' ' -f2
}

# Function to create test payload
create_test_payload() {
    local environment=$1
    local timestamp=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    local commit_id="test-$(date +%s)"
    
    local ref="refs/heads/develop"
    local image_tag="develop-$commit_id"
    
    if [ "$environment" == "production" ]; then
        ref="refs/heads/main"
        image_tag="latest"
    fi
    
    cat <<EOF
{
  "ref": "$ref",
  "repository": {
    "name": "vietnam-admin-api",
    "full_name": "owner/vietnam-admin-api",
    "html_url": "https://github.com/owner/vietnam-admin-api",
    "clone_url": "https://github.com/owner/vietnam-admin-api.git",
    "default_branch": "main"
  },
  "pusher": {
    "name": "test-user",
    "email": "test@example.com"
  },
  "head_commit": {
    "id": "$commit_id",
    "message": "Test deployment to $environment",
    "timestamp": "$timestamp",
    "url": "https://github.com/owner/vietnam-admin-api/commit/$commit_id",
    "author": {
      "name": "Test User",
      "email": "test@example.com"
    }
  },
  "deployment": {
    "environment": "$environment",
    "image_tag": "$image_tag",
    "container_name": "vietnam-admin-api-$environment"
  }
}
EOF
}

# Check if required tools are installed
check_dependencies() {
    local missing=0
    
    if ! command -v curl &> /dev/null; then
        print_color $RED "❌ curl is required but not installed"
        missing=1
    fi
    
    if ! command -v openssl &> /dev/null; then
        print_color $RED "❌ openssl is required but not installed"
        missing=1
    fi
    
    if ! command -v jq &> /dev/null; then
        print_color $YELLOW "⚠️  jq is recommended for better JSON formatting"
    fi
    
    if [ $missing -eq 1 ]; then
        exit 1
    fi
}

# Main function
main() {
    # Check arguments
    if [ $# -lt 2 ]; then
        show_usage
    fi
    
    local webhook_url=$1
    local webhook_secret=$2
    local environment=${3:-staging}
    
    # Check dependencies
    check_dependencies
    
    print_color $CYAN "🔧 Testing webhook deployment for Vietnam Admin API"
    echo "Webhook URL: $webhook_url"
    echo "Environment: $environment"
    echo "Time: $(date)"
    echo ""
    
    # Create payload
    local payload=$(create_test_payload "$environment")
    
    print_color $PURPLE "📤 Sending webhook payload:"
    if command -v jq &> /dev/null; then
        echo "$payload" | jq .
    else
        echo "$payload"
    fi
    echo ""
    
    # Generate signature
    local signature=$(generate_signature "$payload" "$webhook_secret")
    local full_signature="sha256=$signature"
    
    # Generate delivery ID
    local delivery_id="test-$(date +%s)"
    
    print_color $YELLOW "📡 Sending request to webhook server..."
    
    # Send webhook request
    local temp_response=$(mktemp)
    local temp_headers=$(mktemp)
    
    local http_code=$(curl -s -w "%{http_code}" \
        -X POST \
        -H "Content-Type: application/json" \
        -H "X-Hub-Signature-256: $full_signature" \
        -H "X-GitHub-Event: push" \
        -H "X-GitHub-Delivery: $delivery_id" \
        -H "User-Agent: GitHub-Hookshot/test" \
        -D "$temp_headers" \
        -d "$payload" \
        "$webhook_url" > "$temp_response")
    
    # Print request details
    echo "Request Headers:"
    echo "  Content-Type: application/json"
    echo "  X-Hub-Signature-256: $full_signature"
    echo "  X-GitHub-Event: push"
    echo "  X-GitHub-Delivery: $delivery_id"
    echo "  User-Agent: GitHub-Hookshot/test"
    echo ""
    
    # Print response details
    print_color $PURPLE "📥 Response received:"
    echo "Status Code: $http_code"
    echo "Headers:"
    cat "$temp_headers" | sed 's/^/  /'
    echo "Body:"
    local response_body=$(cat "$temp_response")
    if command -v jq &> /dev/null && echo "$response_body" | jq . &> /dev/null; then
        echo "$response_body" | jq .
    else
        echo "$response_body"
    fi
    echo ""
    
    # Clean up temp files
    rm -f "$temp_response" "$temp_headers"
    
    # Check if request was successful
    if [ "$http_code" -eq 200 ]; then
        print_color $GREEN "✅ Webhook test completed successfully!"
        print_color $YELLOW "💡 Check your deployment server logs for deployment progress"
        
        # Suggest next steps
        echo ""
        echo "Next steps:"
        echo "1. Check deployment server logs"
        echo "2. Verify the application is running"
        echo "3. Test the API endpoints"
        
        if [ "$environment" == "production" ]; then
            echo "4. Monitor production deployment"
        else
            echo "4. Test staging environment"
        fi
        
        exit 0
    else
        print_color $RED "❌ Webhook test failed!"
        print_color $YELLOW "💡 Check webhook server configuration and try again"
        
        echo ""
        echo "Troubleshooting steps:"
        echo "1. Verify webhook URL is accessible"
        echo "2. Check webhook secret matches server configuration"
        echo "3. Ensure webhook server is running"
        echo "4. Check server logs for errors"
        echo "5. Verify environment variables are set correctly"
        
        exit 1
    fi
}

# Run main function with all arguments
main "$@" 