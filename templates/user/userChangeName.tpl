{{template "base" .}}

{{ define "embedded-css"}}
<style type="text/css">
    .content {
        max-width: var(--width-small);
    }
</style>
{{end}}

{{define "title"}}Alteração do nome{{end}}

{{define "content"}}
<form class="content" action="/user/change/name" method="post">
    <h2 class="title">Alteração do nome</h2>

    <!-- New name -->
    <label for="new-name">Novo nome</label>
    <input type="text" id="new-name" name="new-name"  value={{.NewName.Value}}>
    <p class="error"> {{.NewName.Msg}} </p>

    <!-- Password -->
    <label for="password">Senha</label>
    <input type="password" id="password" name="password">
    <p class="error"> {{.Password.Msg}} </p>

    <!-- submit -->
    <input type="submit" value="Alterar">
</form>
{{end}}
