SIGECO - Sistema de Gerenciamento e Controle de Entradas
SIGECO é um sistema de desktop simples, desenvolvido em Go e utilizando a biblioteca Fyne, para registrar e monitorar o fluxo de entrada e saída de pessoas.

📸 Visão Geral
O sistema oferece uma interface limpa para cadastrar pessoas (na primeira entrada) e registrar seus horários de entrada e saída. Todos os eventos são exibidos em um log em tempo real.


✨ Funcionalidades
[x] Cadastro de Pessoas: Registra ID (RG/CPF), Nome e Telefone.
[x] Registro de Entrada: Marca o horário exato que a pessoa entrou.
[x] Registro de Saída: Atualiza o registro da pessoa com o horário exato que ela saiu.
[x] Log de Eventos: Exibe uma lista em tempo real de todas as entradas e saídas, formatadas com o horário.

💻 Tecnologias Utilizadas
Go (Golang): Linguagem principal do projeto.

Fyne: Biblioteca gráfica (GUI) 100% em Go para criar o front-end de desktop.

🚀 Como Executar o Projeto
Para rodar este projeto localmente, você precisará ter o Go (versão 1.19 ou superior) e um compilador C (gcc) instalados.

1. Pré-requisitos (Dependências do Fyne)
O Fyne precisa de algumas bibliotecas gráficas para funcionar. Instale-as de acordo com seu sistema operacional:

Linux (Debian/Ubuntu):

Bash

sudo apt install build-essential libgl1-mesa-dev xorg-dev
Windows (64-bits): A forma mais fácil é instalar o compilador TDM-GCC. Certifique-se de marcar a opção "Add to PATH" durante a instalação.

macOS: Instale o Xcode pela App Store ou rode o seguinte comando no terminal:

Bash

xcode-select --install
2. Clonar o Repositório
Bash

git clone https://github.com/LuanaMonteiro0/sigeco.git
cd sigeco
3. Baixar Dependências do Go
O Go cuidará disso automaticamente com o comando:

Bash

go mod tidy
4. Executar o Aplicativo
Bash

go run .
A primeira execução pode demorar um pouco para compilar tudo. As execuções seguintes serão quase instantâneas.