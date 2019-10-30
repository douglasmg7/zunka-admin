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
        <label {{if .Name.Msg}} class="error" {{end}} for="name">
            {{if not .Name.Msg}} Nome {{end}}
            {{if .Name.Msg}} {{.Name.Msg}} {{end}}
         </label> 
        <input type="text" id="name" name="name" value={{.Name.Value}}>

        <!-- email -->
        <label {{if .Email.Msg}} class="error" {{end}} for="email">
            {{if not .Email.Msg}} Email {{end}}
            {{if .Email.Msg}} {{.Email.Msg}} {{end}}
        </label> 
        <input type="text" id="email" name="email" value={{.Email.Value}}>

        <!-- password -->
        <label {{if .Password.Msg}} class="error" {{end}} for="password">
            {{if not .Password.Msg}} Senha {{end}}
            {{if .Password.Msg}} {{.Password.Msg}} {{end}}
        </label> 
        <input type="password" id="password" name="password" value={{.Password.Value}}>

        <!-- confirm password -->
        <label {{if .PasswordConfirm.Msg}} class="error" {{end}} for="passwordConfirm">
            {{if not .PasswordConfirm.Msg}} Confirme a senha {{end}}
            {{if .PasswordConfirm.Msg}} {{.PasswordConfirm.Msg}} {{end}}
        </label> 
        <input type="password" id="passwordConfirm" name="passwordConfirm" value={{.PasswordConfirm.Value}}>

        <!-- submit -->
        <input class="btn btn-info" type="submit" value="Cadastrar">

        <!-- Foot message. -->
        {{if .SuccessMsg}} <div class="success-msg"> {{.SuccessMsg}} </div> {{end}}
        {{if .WarnMsg}} <div class="warn-msg"> {{.WarnMsg}} </div> {{end}}
    </form>
{{end}}
