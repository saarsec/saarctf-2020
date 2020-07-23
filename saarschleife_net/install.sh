#!/usr/bin/env bash

set -eux

# Install the service on a fresh vulnbox. Target should be /home/<servicename>
# You get:
# - $SERVICENAME
# - $INSTALL_DIR
# - An user account with your name

# 1. Install dependencies
# moved to Dockerfile


# 2. Copy/move files
mv service/patches /tmp/
mv service/* "$INSTALL_DIR/"
# init script
cat - <<EOF > $INSTALL_DIR/run-server.sh
#!/usr/bin/env bash
if [ ! -f $INSTALL_DIR/backend/build/libs/Saarschleife-Server.jar ]; then
    echo "Building Server..."
    pushd .
    cd $INSTALL_DIR/backend/
    gradle shadowJar --no-daemon
    gradle --stop
    popd
fi

exec java -jar $INSTALL_DIR/backend/build/libs/Saarschleife-Server.jar
EOF
chmod +x $INSTALL_DIR/run-server.sh
# chown things
chown -R "$SERVICENAME:$SERVICENAME" "$INSTALL_DIR"

# 3. Build on box
pushd .
cd $INSTALL_DIR/backend
# Patch library
sudo -u saarschleife gradle dependencies --no-daemon
sudo -u saarschleife find $INSTALL_DIR/.gradle/caches/modules-*/files-*/org.litote.kmongo -name 'kmongo-coroutine-core-3.10.?.jar' \
     -exec cp -v /tmp/patches/kmongo-coroutine-core-3.10.1-patched.jar {} \;
sudo -u saarschleife gradle shadowJar --no-daemon
rm -f build/libs/*.jar src/main/kotlin/saarschleife/key.kt
popd


# 4. Configure startup for your service
# Typically use systemd for that:
# Install backend as systemd service
service-add-advanced "$INSTALL_DIR/run-server.sh" "$INSTALL_DIR/backend/" "Saarschleife.net social network" <<EOF
Restart=on-failure
RestartSec=10
EOF

# Fix bug with Kotlin daemon
sleep 5
rm -rf $INSTALL_DIR/.kotlin/daemon || true
