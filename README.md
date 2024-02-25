# Rinha de Backend 2023 Q3

Bem-vindo ao repositório da Rinha de Backend do terceiro trimestre de 2023, organizada por [Francisco Zanfranceschi](https://github.com/zanfranceschi).

## Repositório da Rinha

Caso queria mais informações sobre a rinha:

[Repositório da Rinha](https://github.com/zanfranceschi/rinha-de-backend-2023-q3)

## Métricas da Execução do Teste de Estresse

Consegui alcançar mais de 46500 inserções no banco de dados.

![Métricas da execução do teste de estresse](/misc/stress-test.png)

## Como Executar os Testes

Para executar os testes, utilize os seguintes comandos:

```bash
make up-dev && make test && make down-dev
```

## Executar o Teste de Estresse

Para realizar o teste de estresse, execute os seguintes comandos:

```bash
make build && make up && make stress && make down
```

## Repositório de Referências

[Repositório de Consulta](https://github.com/leorcvargas/rinha-go)

## Vídeos do Akita On Rails sobre a Rinha

Incorporei algumas abordagens apresentadas nos vídeos abaixo.

- [Vídeo 1](https://www.youtube.com/watch?v=EifK2a_5K_U&t=2352s&pp=ygUFYWtpdGE%3D)
- [Vídeo 2](https://www.youtube.com/watch?v=-yGHG3pnHLg&pp=ygUFYWtpdGE%3D)
