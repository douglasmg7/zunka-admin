{{template "base" .}}

{{ define "embedded-css"}}
<style type="text/css">
    .content {
        max-width: var(--width-small);
    }
    a.reset-pass {
        display: block;
        margin: .2em 0 1em 0;
    }
    p {
        margin-bottom: 0;
    }
</style>
{{end}}

{{define "title"}}Cadastro{{end}}

{{define "header"}}{{end}}

{{define "content"}}
    <form class="content" action="/ns/auth/signup" method="post">
        <h2 class="title">Criar cadastro</h2>
        <!-- name -->
        <label for="name">Nome</label>
        <input type="text" id="name" name="name" value={{.Name.Value}}>
        <p class="error">{{.Name.Msg}}</p>

        <!-- email -->
        <label for="email">E-mail</label>
        <input type="text" id="email" name="email" value={{.Email.Value}}>
        <p class="error">{{.Email.Msg}}</p>

        <!-- password -->
        <label for="password">Senha</label>
        <input type="password" id="password" name="password" value={{.Password.Value}}>
        <p class="error">{{.Password.Msg}}</p>

        <!-- confirm password -->
        <label for="passwordConfirm">Confirme a senha</label>
        <input type="password" id="passwordConfirm" name="passwordConfirm" value={{.PasswordConfirm.Value}}>
        <p class="error">{{.PasswordConfirm.Msg}}</p>

        <!-- submit -->
        <button type="submit">Cadastrar</button>

        <!-- Foot message. -->
        {{if .SuccessMsg}} <div class="success-msg"> {{.SuccessMsg}} </div> {{end}}
        {{if .WarnMsg}} <div class="warn-msg"> {{.WarnMsg}} </div> {{end}}
    </form>
{{end}}
