echo "====================="
echo "Running Service ....."

# Stop the "app" service
echo "// Stopping 'app' service"
docker-compose down --remove-orphans

# Start the "app" service and rebuild it
echo "// Starting 'app' service and rebuilding"
docker-compose up -d --no-deps --build app

# Check if the "app" service is up
if docker-compose ps app | grep -q "Up"; then
    echo "Service is up!"
else
    echo "Failed to start the service."
    exit 1
fi

echo "====================="
