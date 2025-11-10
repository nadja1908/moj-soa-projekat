FROM openjdk:17-jdk-slim

# Install wget for health checks
RUN apt-get update && apt-get install -y wget && apt-get clean

# Download Eureka Server JAR (Spring Boot standalone)
WORKDIR /app

# Create a simple Spring Boot Eureka Server
COPY eureka-server.jar /app/eureka-server.jar

EXPOSE 8761

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=20s --retries=5 \
  CMD wget -q --timeout=10 -O /dev/null http://localhost:8761/actuator/health || exit 1

CMD ["java", "-jar", "eureka-server.jar"]