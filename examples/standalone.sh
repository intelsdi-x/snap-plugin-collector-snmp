#!/bin/bash

./build/linux/x86_64/snap-plugin-collector-snmp --config '{
    "setfile": "setfile.json",
    "snmp_agent_name": "host1",
    "snmp_agent_address": "127.0.0.1:161",
    "snmp_version": "v3",
    "network": "udp",
    "user_name": "bootstrap",
    "security_level": "AuthPriv",
    "auth_password": "password",
    "auth_protocol": "MD5",
    "priv_protocol": "DES",
    "priv_password": "password"
}'
