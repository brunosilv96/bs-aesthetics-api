# Sistema de Agendamento
- Excalidraw Link: https://excalidraw.com/#json=86KmmTihxzH2SvppqOe95,qVzg2w-Kczk2vmTj5WWmpg

## Dores do negócio
* Sem controle dos clientes
* Sem controle de agenda e horários
* Sem controle dos produtos disponíveis para venda

## Como é feito hoje?
* Agendamento e vendas feito pelo WPP

## Cliente - Site ou APP
* Cadastro com Login (Perfil):
    * Nome
    * E-mail
    * Celular ou WPP
    * Data de nascimento (para campanhas de aniversário)
* Listagem dos serviços oferecidos:
    * Descrição
    * Preço
* Listagem dos produtos oferecidos:
    * Descrição
    * Preço
    * Estoque
* Agenda de serviço com dias e horários disponíveis: 
    * Agendar
    * Remarcar
    * Cancelar
* Histórico de agendamentos
* Link simples com WPP

## Administrador - BackOffice - Site
* Autenticação
* Dashboard
    * Agendamento do dia
* CRUD de clientes (cadastro manual)
* CRUD de serviços (descrição, preço e tempo)
* CRUD de produtos (descrição, preço e disponibilidade)
* CRUD de dados da empresa:
    * Meios de contato 
    * Informações 
    * Dias e horários de funcionamento
* Gerenciamento de agenda (dia e horários):
    * Agendar
    * Alterar
    * Desmarcar 
    * Confirmar
    * Status: Confirmado, Não Compareceu, Pendente, Realizado, Cancelado
* Histórico de agendamentos do dia ou por cliente
* Relatório do dia e mensal
