function add_docker_metadata(tag, timestamp, record)
    -- Extract container name from Docker metadata
    if record["container_name"] then
        local container_name = record["container_name"]
        -- Extract service name from container name
        if string.match(container_name, "follower%-service") then
            record["service_name"] = "follower-service"
        elseif string.match(container_name, "auth%-service") then
            record["service_name"] = "auth-service"
        elseif string.match(container_name, "tour%-service") then
            record["service_name"] = "tour-service"
        elseif string.match(container_name, "blog%-service") then
            record["service_name"] = "blog-service"
        elseif string.match(container_name, "purchase%-service") then
            record["service_name"] = "purchase-service"
        else
            record["service_name"] = "unknown"
        end
    end
    return 1, timestamp, record
end