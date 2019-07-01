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

{{define "title"}}Autenticação{{end}}

{{define "content"}}
<form class="content" action="/auth/signin" method="post">
    <h2 class="title">Entrar</h2>

    <!-- Head messages. -->
    {{if .SuccessMsgHead}} <div class="success-msg"> {{.SuccessMsgHead}} </div> {{end}}
    {{if .WarnMsgHead}} <div class="warn-msg"> {{.WarnMsgHead}} </div> {{end}}

    <!-- email -->
    <label for="email">E-mail</label>
    <input type="text" id="email" name="email"  value={{.Email.Value}}>
    <p class="error"> {{.Email.Msg}} </p>

    <!-- password -->
    <label for="password">Senha</label>
    <input type="password" id="password" name="password" value={{.Password.Value}}>
    <p class="error">{{.Password.Msg}}</p>

    <!-- submit -->
    <input type="submit" value="Entrar">

    <!-- reset password -->
    <a class="reset-pass" href="/auth/password/recovery">Esqueceu a senha?</a>

    <!-- signup -->
    <p>Não tem cadastro? </p>
    <a class="signup" href="/auth/signup">Criar cadastro</a>

    <!-- Foot messages. -->
    {{if .SuccessMsgFooter}} <div class="success-msg"> {{.SuccessMsgFooter}} </div> {{end}}
    {{if .WarnMsgFooter}} <div class="warn-msg"> {{.WarnMsgFooter}} </div> {{end}}
</form>
{{end}}
