{{template "base" .}}

{{ define "embedded-css"}}
<style type="text/css">
    .content {
        max-width: var(--width-small);
    }
    div.btns-row {
        display: flex;
        flex-wrap: wrap;
        justify-content: space-between;
    }
    div.btns-row *:not(:last-child) {
        margin-right: 1em;
    }
</style>
{{end}}

{{define "title"}}Alteração do número do celular{{end}}

{{define "content"}}
<form class="content" action="/user/change/mobile" method="post">
    <h2 class="title">Alteração do número de celular</h2>

    <!-- New mobile number -->
    <label for="new-mobile">Novo número de celular</label>
    <input type="text" id="new-mobile" name="new-mobile"  value={{.NewMobile.Value}}>
    <p class="error"> {{.NewMobile.Msg}} </p>

    <!-- Password -->
    <label for="password">Senha</label>
    <input type="password" id="password" name="password">
    <p class="error"> {{.Password.Msg}} </p>

    <!-- submit -->
    <div class="btns-row">
        <a class="button btn-danger" href="/user/account">Cancelar</a>
        <input type="submit" value="Test">
        <input class="btn-danger" type="submit" value="Salvar">
    </div>
</form>
{{end}}
