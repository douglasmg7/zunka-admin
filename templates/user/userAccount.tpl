{{template "base" .}}

{{ define "embedded-css"}}
<style type="text/css">
    .content {
        max-width: var(--width-small);
    }
    .title {
        <--! color: var(--blue-light); -->
    }
    .title + h4 {
        margin-top: 0;
        opacity: .7;
    }
    a.button {
        display: block;
    }
</style>
{{end}}

{{define "title"}}Dados da conta{{end}}

{{define "header"}}{{end}}

{{define "content"}}
    <div class="content">
        <h2>Dados da conta</h2>

        <h4 class="title">Nome</h4>
        <h4>{{.Name}}</h4>
        <a href="/ns/user/change/name">Alterar</a>

        <h4 class="title">Email</h4>
        <h4>{{.Email}}</h4>
        <a href="/ns/user/change/email">Alterar</a>

        <h4 class="title">Número de celular</h4>
        <h4>{{if .Mobile}} {{.Mobile}} {{end}}</h4>
        <a href="/ns/user/change/mobile">Alterar</a>

        <h4 class="title">Senha</h4>
        <a href="/ns/user/change/password">Alterar</a>

        <h4 class="title">RG</h4>
        <h4>{{if .RG}} {{.RG}} {{end}}</h4>
        <a href="/ns/user/change/rg">Alterar</a>

        <h4 class="title">CPF</h4>
        <h4>{{if .CPF}} {{.CPF}} {{end}}</h4>
        <a href="/ns/user/change/cpf">Alterar</a>

        <h4 class="title">Apagar conta</h4>
        <a href="/ns/user/delete/account">Apagar</a>

        <!-- submit -->
        <a class="button" href="/ns/">Sair</a>
    </div>
{{end}}
