# TgBot
A simple Bot Project

## What I can DO

- To send a quote                      `/mingyan`
- Translate                                `/translate`
- Get weather conditions         `/weather`
- Get Music                                 `/music`
- Get coin value                          `/btc`  `/xmr`  `/eth`
- <u>***And So On!***</u>

## Quick Start

1. Create a config.go in config and fix blank below

   ```go
   package config
   
   const TOKEN1 = ""             //Main bot TOKEN
   const TOKEN2 = ""             //Test bot TOKEN,optional
   const DB_TOKEN = ""           //Database1
   const CONFIG_TOKEN = ""       //Database2
   const WEATHER_TOKEN = ""      //visit `dev.qweather.com` to apply a `TOKEN`
   const TP_IP1 = ""             
   const TP_IP2 = ""
   const SHORT_IP = ""          
   const BOT_CONFIG = ""         //webHook URL
   ```

2. Run

   ```shell
   go mod tidy
   go build
   ./bot &
   ```

## Join my telegram group

[touch it](https://t.me/AllenBot_Group)

