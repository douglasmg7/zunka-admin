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
    label > span {
        font-weight: 500;
    }
</style>
{{end}}

{{define "title"}}Autenticação{{end}}

{{define "content"}}
<form class="content" action="/ns/auth/signin" method="post">
    <h2 class="title">Entrar</h2>

    <!-- Head messages. -->
    {{if .SuccessMsgHead}} <div class="success-msg"> {{.SuccessMsgHead}} </div> {{end}}
    {{if .WarnMsgHead}} <div class="warn-msg"> {{.WarnMsgHead}} </div> {{end}}

    <!-- email -->
    <label {{if .Email.Msg}} class="error" {{end}} for="email">
        {{if not .Email.Msg}} Email {{end}}
        {{if .Email.Msg}} {{.Email.Msg}} {{end}}
    </label> 
    <input type="text" id="email" name="email"  value={{.Email.Value}}>

    <!-- password -->
    <label {{if .Password.Msg}} class="error" {{end}} for="password">
        {{if not .Password.Msg}} Senha {{end}}
        {{if .Password.Msg}} {{.Password.Msg}} {{end}}
    </label> 
    <input type="password" id="password" name="password" value={{.Password.Value}}>

    <!-- submit -->
    <input class="btn btn-info" type="submit" value="Entrar">

    <!-- reset password -->
    <a class="reset-pass" href="/ns/auth/password/recovery">Esqueceu a senha?</a>

    <!-- signup -->
    <p>Não tem cadastro? </p>
    <a class="signup" href="/ns/auth/signup">Criar cadastro</a>

    <!-- Foot messages. -->
    {{if .SuccessMsgFooter}} <div class="success-msg"> {{.SuccessMsgFooter}} </div> {{end}}
    {{if .WarnMsgFooter}} <div class="warn-msg"> {{.WarnMsgFooter}} </div> {{end}}
</form>
{{end}}
