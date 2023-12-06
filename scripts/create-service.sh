#!/bin/bash

username="syndica"
user_home="/home/$username"

# Service information
service_name="goserver"
service_executable="$user_home/$service_name.sh"

# Create the user
adduser "$username"

# Create the user home directory
mkdir -p "$user_home"

# Create the service executable
echo '#!/bin/bash' > "$service_executable"
echo 'while true; do' >> "$service_executable"
echo '  # Add your service logic here' >> "$service_executable"
echo '  sleep 10' >> "$service_executable"
echo 'done' >> "$service_executable"
chmod +x "$service_executable"

# Create the systemd service file
echo '[Unit]' > "/etc/systemd/system/$service_name.service"
echo "Description=$service_description" >> "/etc/systemd/system/$service_name.service"
echo '[Service]' >> "/etc/systemd/system/$service_name.service"
echo "User=$username" >> "/etc/systemd/system/$service_name.service"
echo "Group=$username" >> "/etc/systemd/system/$service_name.service"
echo "WorkingDirectory=$user_home" >> "/etc/systemd/system/$service_name.service"
echo "ExecStart=$service_executable" >> "/etc/systemd/system/$service_name.service"
echo "Restart=always" >> "/etc/systemd/system/$service_name.service"
echo '[Install]' >> "/etc/systemd/system/$service_name.service"
echo "WantedBy=multi-user.target" >> "/etc/systemd/system/$service_name.service"

# Start the service
systemctl daemon-reload
systemctl start "$service_name"

# Enable the service to start automatically at boot
systemctl enable "$service_name"

echo "User account '$username' created."
echo "Service '$service_name' created for user '$username'."
