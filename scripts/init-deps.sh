echo "====================="
echo "Initialize Dependencies"

echo "// Initialize MQTT Broker Service"
docker compose up -d --no-deps --build mqtt_broker

echo "// Initialize Database Service"
docker compose up -d --no-deps --build db