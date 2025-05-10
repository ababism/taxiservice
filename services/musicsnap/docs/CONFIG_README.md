# MusicSnap Описание конфигураций

## Настройка
Для сборки сервиса используется .docker.env для `docker-compose`. И .env для тестирования.

Наполнение файлов .env.docker и .env можно найти в .docker.env.dist и .env.dist.

[Локальные](../config/config.local.yml) и [докер](../config/config.docker.yml) конфигурации в папке [../config](../config)

Важно также не забыть поменять конфиги в deployments/prometheus/prometheus.yml, если надо
```yaml
  - job_name: musicsnapService
    honor_timestamps: true
    scrape_interval: 15s
    scrape_timeout: 10s
    metrics_path: /system/metrics/prometheus
    scheme: http
    static_configs:
      - targets:
          - docker.for.mac.localhost:8080
```

Конфигурацию сервера можно также посмотреть в Go коде тут:
[config.go](..%2Finternal%2Fconfig%2Fconfig.go)