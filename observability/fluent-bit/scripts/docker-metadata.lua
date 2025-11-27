function add_docker_metadata(tag, timestamp, record)
    -- Ekstraktuj container ID iz tag-a
    local container_id = string.match(tag, "docker%.(.+)")
    if container_id then
        record["container_id"] = container_id:sub(1, 12)
    end
    
    -- Parsiranje Docker JSON log format
    local actual_log = ""
    if record["log"] then
        actual_log = tostring(record["log"])
        record["message"] = actual_log  -- Postavi clean message
    end
    
    -- Parsiranje HTTP logova iz message
    if actual_log ~= "" then
        -- ASP.NET Core log format parsing
        local method, path, status = string.match(actual_log, "(%w+)%s+([%w/%-_%.%?&=]+)[%s%w]*%s+responded%s+(%d+)")
        if not method then
            -- Alternative format
            method, path = string.match(actual_log, "([A-Z]+)%s+([%w/%-_%.%?&=]+)")
            status = string.match(actual_log, "status[:%s]*(%d+)")
        end
        
        if method then
            record["http_method"] = method
        end
        if path then
            record["http_path"] = path
        end
        if status then
            record["http_status"] = status
            -- Status kategorije
            if status:match("^2") then
                record["http_status_category"] = "success"
            elseif status:match("^4") then
                record["http_status_category"] = "client_error"
            elseif status:match("^5") then
                record["http_status_category"] = "server_error"
            else
                record["http_status_category"] = "other"
            end
        end
        
        -- Log level extraction
        local level = string.match(actual_log, "%[(%w+)%]") or string.match(actual_log, "(%u+):")
        if level then
            record["log_level"] = string.lower(level)
        end
    end
    
    local log_content = string.lower(actual_log)
    
    -- Detektuje servis na osnovu log sadržaja
    if string.find(log_content, "follower") or string.find(log_content, "8080") or string.find(log_content, "kestrel") then
        record["service_name"] = "follower-service"
        record["container_name"] = "follower-service"
    elseif string.find(log_content, "auth") or string.find(log_content, "8003") then
        record["service_name"] = "auth-service"
        record["container_name"] = "auth-service"
    elseif string.find(log_content, "tour") or string.find(log_content, "8004") then
        record["service_name"] = "tour-service"
        record["container_name"] = "tour-service"
    elseif string.find(log_content, "blog") or string.find(log_content, "8002") then
        record["service_name"] = "blog-service"
        record["container_name"] = "blog-service"
    elseif string.find(log_content, "purchase") or string.find(log_content, "8005") then
        record["service_name"] = "purchase-service"
        record["container_name"] = "purchase-service"
    elseif string.find(log_content, "stakeholder") or string.find(log_content, "8001") then
        record["service_name"] = "stakeholders-service"
        record["container_name"] = "stakeholders-service"
    elseif string.find(log_content, "gateway") or string.find(log_content, "8000") then
        record["service_name"] = "api-gateway"
        record["container_name"] = "api-gateway"
    else
        record["service_name"] = "unknown"
        record["container_name"] = record["container_id"] or "unknown"
    end
    
    return 1, timestamp, record
end