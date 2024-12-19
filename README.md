# Raspberry Pi Kiosk Script

## Overview
This project provides scalable software for Raspberry Pi devices to display specific URLs in Chromium's kiosk mode. It supports HTML pages, video players, and more. Each device registers with the server, sends its unique identifier (UUID), checks for updated URLs periodically, and transmits performance metrics, including CPU temperature, load, memory usage, and logs.

## Features
- **Automatic Kiosk Mode**: Launch Chromium in kiosk mode with the specified URL.
- **Device Registration**: Register devices with a unique UUID.
- **Health Monitoring**: Send performance data (CPU temperature, load, memory usage, and logs) to the server.
- **URL Updates**: Periodically fetch the latest URL for display.
