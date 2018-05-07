#!/usr/bin/env bash
set -e
systemctl daemon-reload
systemctl enable factd.service
systemctl start factd.service