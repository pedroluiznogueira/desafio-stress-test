# Stress Test

CLI em Go para teste de carga: dispara `--requests` requisições HTTP contra `--url`, distribuídas
por `--concurrency` chamadas simultâneas (worker pool sobre um channel pré-carregado, garantindo
o total exato mesmo sob concorrência), e imprime um relatório ao final.

## Build

```
docker build -t stress-test .
```

## Rodando

```
docker run stress-test --url=http://google.com --requests=1000 --concurrency=10
```

## Saída

```
Relatório de Stress Test
Tempo total gasto: 1.2s
Total de requests realizados: 1000
Requests com status 200: 980
Distribuição de outros status codes:
  404: 15
  500: 5
```
