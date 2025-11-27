function add_docker_metadata(tag, timestamp, record)
    -- Ekstraktuj container ID iz tag-a
    local container_id = string.match(tag, "docker%.(.+)")
    if container_id then
        record["container_id"] = container_id:sub(1, 12)
    end
    
    -- Odredi naziv servisa na osnovu sadržaja loga ili container_id-ja
    -- Koristimo static mapiranje dok ne možemo dinamički da dobijamo imena
    
    -- Pokušaj da detektuje servis na osnovu log sadržaja
    local log_content = ""
    if record["log"] then
        log_content = string.lower(tostring(record["log"]))
    elseif record["message"] then
        log_content = string.lower(tostring(record["message"]))
    end
    
    if string.find(log_content, "follower") or string.find(log_content, "8080") then
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
        -- Default vrednosti
        record["service_name"] = "unknown"
        record["container_name"] = record["container_id"] or "unknown"
    end
    
    return 1, timestamp, record
end