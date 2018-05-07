#!/usr/bin/env bash
set -e
systemctl stop factd.service
systemctl disable factd.service