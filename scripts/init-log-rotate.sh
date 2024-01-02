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

echo "$LOG_FILE $LOG_ROTATE_OPTIONS" | sudo tee "$LOG_ROTATE_CONFIG" > /dev/null

CRON_ENTRY="0 0 * * 0 /usr/sbin/logrotate $LOG_ROTATE_CONFIG"
if ! crontab -l | grep -q "$CRON_ENTRY"; then
    (crontab -l 2>/dev/null; echo "$CRON_ENTRY") | crontab -e -
fi
