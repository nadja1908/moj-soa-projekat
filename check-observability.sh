#!/bin/bash

echo "🔍 Observability Stack - Health Check Script"
echo "=============================================="

# Color codes
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

check_service() {
    local service_name=$1
    local url=$2
    local expected_code=${3:-200}
    
    echo -n "Checking $service_name... "
    
    if curl -s -o /dev/null -w "%{http_code}" "$url" | grep -q "$expected_code"; then
        echo -e "${GREEN}✓ OK${NC}"
        return 0
    else
        echo -e "${RED}✗ FAILED${NC}"
        return 1
    fi
}

echo ""
echo "🎯 Observability Services:"
echo "-------------------------"

check_service "Jaeger UI" "http://localhost:16686"
check_service "Grafana" "http://localhost:3001" "302"  # Grafana redirects to login
check_service "Prometheus" "http://localhost:9090"
check_service "Kibana" "http://localhost:5601" "302"   # Kibana redirects
check_service "Elasticsearch" "http://localhost:9200"
check_service "Node Exporter" "http://localhost:9100"
check_service "cAdvisor" "http://localhost:8081"

echo ""
echo "📊 Application Services:"
echo "------------------------"

check_service "Auth Service" "http://localhost:8003/health"
check_service "Auth Metrics" "http://localhost:8003/metrics"
check_service "API Gateway" "http://localhost:8000/health"
check_service "Stakeholders Service" "http://localhost:8001/health"

echo ""
echo "🧪 Testing Tracing:"
echo "-------------------"

echo "Sending test login request to generate traces..."
curl -s -X POST http://localhost:8000/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass"}' > /dev/null

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Test request sent${NC}"
    echo -e "${YELLOW}ℹ Check traces in Jaeger: http://localhost:16686${NC}"
else
    echo -e "${RED}✗ Failed to send test request${NC}"
fi

echo ""
echo "📈 Quick Access URLs:"
echo "---------------------"
echo "🎯 Jaeger (Tracing):     http://localhost:16686"
echo "📊 Grafana (Dashboards): http://localhost:3001 (admin/admin)"
echo "⚡ Prometheus (Metrics):  http://localhost:9090"
echo "🔍 Kibana (Logs):        http://localhost:5601"
echo ""
echo -e "${YELLOW}Note: Grafana is on port 3001 to avoid conflicts!${NC}"