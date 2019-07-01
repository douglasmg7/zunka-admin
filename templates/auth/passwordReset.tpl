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
    <h2 class="title">Redefinir senha</h2>

    <!-- email -->
    <label for="email">E-mail</label>
    <input type="text" id="email" name="email"  value={{.Email.Value}}>
    <p class="error"> {{.Email.Msg}} </p>

    <!-- submit -->
    <input type="submit" value="Redefinir">

    <!-- reset password -->
    <a class="reset-pass" href="/auth/reset_password">Esqueceu a senha?</a>

    <!-- signup -->
    <p>Não tem cadastro? </p>
    <a class="signup" href="/auth/signup">Criar cadastro</a>

    <!-- Foot messages. -->
    {{if .SuccessMsgFooter}} <div class="success-msg"> {{.SuccessMsgFooter}} </div> {{end}}
    {{if .WarnMsgFooter}} <div class="warn-msg"> {{.WarnMsgFooter}} </div> {{end}}
</form>
{{end}}
