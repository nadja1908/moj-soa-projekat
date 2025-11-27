#!/bin/bash
# Docker Stats to Prometheus Script

echo "# HELP docker_container_cpu_percent CPU percentage used by container"
echo "# TYPE docker_container_cpu_percent gauge"

echo "# HELP docker_container_memory_usage_mb Memory usage in MB by container"
echo "# TYPE docker_container_memory_usage_mb gauge"

docker stats --no-stream --format "{{.Container}} {{.CPUPerc}} {{.MemUsage}}" | while read line; do
    container=$(echo $line | awk '{print $1}')
    cpu=$(echo $line | awk '{print $2}' | sed 's/%//')
    mem=$(echo $line | awk '{print $3}' | sed 's/MiB//')
    
    echo "docker_container_cpu_percent{container=\"$container\"} $cpu"
    echo "docker_container_memory_usage_mb{container=\"$container\"} $mem"
done