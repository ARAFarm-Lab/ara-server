SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
LOG_FILE="$SCRIPT_DIR/../log/ara-server.log"

LOG_ROTATE_CONFIG="/etc/logrotate.d/ara-server"
LOG_ROTATE_OPTIONS="{
    weekly
    rotate 2
    compress
    missingok
    notifempty
}"

# Create the logrotate configuration file
echo "$LOG_FILE $LOG_ROTATE_OPTIONS" | sudo tee "$LOG_ROTATE_CONFIG" > /dev/null

# Check if the crontab entry already exists before adding it
CRON_ENTRY="0 0 * * 0 /usr/sbin/logrotate $LOG_ROTATE_CONFIG"
if ! crontab -l | grep -q "$CRON_ENTRY"; then
    # Add the crontab entry to run logrotate weekly
    (crontab -l 2>/dev/null; echo "$CRON_ENTRY") | crontab -
fi
